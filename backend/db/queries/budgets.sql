-- name: ListBudgets :many
SELECT
  b.id, b.category_id, b.`year_month`, b.amount, b.created_at, b.updated_at,
  c.name AS category_name
FROM budgets b
LEFT JOIN categories c ON c.id = b.category_id
WHERE b.`year_month` = ?
ORDER BY (b.category_id IS NOT NULL), c.name;

-- name: GetBudget :one
SELECT
  b.id, b.category_id, b.`year_month`, b.amount, b.created_at, b.updated_at,
  c.name AS category_name
FROM budgets b
LEFT JOIN categories c ON c.id = b.category_id
WHERE b.id = ?;

-- name: CreateBudget :execlastid
INSERT INTO budgets (category_id, `year_month`, amount)
VALUES (?, ?, ?);

-- name: UpdateBudget :execrows
UPDATE budgets
SET amount = ?
WHERE id = ?;

-- name: DeleteBudget :execrows
DELETE FROM budgets
WHERE id = ?;

-- name: CountMonthWideBudget :one
-- 月全体予算（category_id IS NULL）の重複チェック。
-- MySQL の UNIQUE 制約は NULL の重複を許すため、アプリ側でも検査する。
SELECT COUNT(*) AS cnt
FROM budgets
WHERE `year_month` = ? AND category_id IS NULL;

-- name: CountCategoryBudget :one
-- カテゴリ別予算の重複チェック。
SELECT COUNT(*) AS cnt
FROM budgets
WHERE `year_month` = ? AND category_id = ?;
