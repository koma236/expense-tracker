CREATE TABLE budgets (
  id          BIGINT   NOT NULL AUTO_INCREMENT,
  category_id  BIGINT   NULL,
  `year_month` CHAR(7)  NOT NULL,
  amount       INT      NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_budgets_year_month (`year_month`),
  UNIQUE KEY uq_budgets_category_year_month (category_id, `year_month`),
  CONSTRAINT fk_budgets_category_id
    FOREIGN KEY (category_id) REFERENCES categories (id),
  CONSTRAINT chk_budgets_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
