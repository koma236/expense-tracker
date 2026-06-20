CREATE TABLE transactions (
  id          BIGINT       NOT NULL AUTO_INCREMENT,
  occurred_on DATE         NOT NULL,
  amount      INT          NOT NULL,
  type        ENUM('income','expense') NOT NULL,
  category_id BIGINT       NOT NULL,
  memo        VARCHAR(255) NULL,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_transactions_occurred_on (occurred_on),
  KEY idx_transactions_category_id (category_id),
  CONSTRAINT fk_transactions_category_id
    FOREIGN KEY (category_id) REFERENCES categories (id),
  CONSTRAINT chk_transactions_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
