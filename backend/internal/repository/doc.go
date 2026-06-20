// Package repository はデータアクセス層。
//
// 型安全なクエリコードは sqlc により internal/repository/db 配下へ生成する
// （リポジトリルートで `sqlc generate` を実行）。生成物はコミット対象とし、
// repository ではそれをラップして service に提供する想定。
package repository
