package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/koma236/expense-tracker/backend/internal/handler"
	db "github.com/koma236/expense-tracker/backend/internal/repository/db"
	"github.com/koma236/expense-tracker/backend/internal/service"
)

func main() {
	// .env を読み込む（存在しなくてもエラーにしない）
	_ = godotenv.Load()

	sqlDB, err := openDB()
	if err != nil {
		log.Fatalf("DB の初期化に失敗しました: %v", err)
	}
	defer sqlDB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{getEnv("CORS_ALLOW_ORIGIN", "http://localhost:3000")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}))

	queries := db.New(sqlDB)
	health := handler.NewHealthHandler(sqlDB)
	categoryH := handler.NewCategoryHandler(queries)
	txH := handler.NewTransactionHandler(service.NewTransactionService(queries))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", health.Health)
		r.Get("/categories", categoryH.List)
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", txH.List)
			r.Post("/", txH.Create)
			r.Get("/{id}", txH.Get)
			r.Put("/{id}", txH.Update)
			r.Delete("/{id}", txH.Delete)
		})
	})

	addr := ":" + getEnv("APP_PORT", "8080")
	log.Printf("サーバーを起動します: http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("サーバーが異常終了しました: %v", err)
	}
}

// openDB は接続プールを用意する。実際の接続確認は /api/health の Ping で行うため、
// ここでは Ping せず、MySQL 未起動でもサーバー自体は起動できるようにしている。
func openDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo&charset=utf8mb4",
		getEnv("DB_USER", "app"),
		getEnv("DB_PASSWORD", "app"),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "expense_tracker"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
