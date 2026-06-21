-- name: GetMonthlyTotals :one
SELECT
  CAST(COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS SIGNED) AS income,
  CAST(COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS SIGNED) AS expense
FROM transactions
WHERE occurred_on BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date);

-- name: ListExpenseByCategory :many
SELECT
  t.category_id,
  c.name AS category_name,
  CAST(SUM(t.amount) AS SIGNED) AS total
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.type = 'expense'
  AND t.occurred_on BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
GROUP BY t.category_id, c.name
ORDER BY total DESC;

-- name: ListMonthlyTrend :many
SELECT
  DATE_FORMAT(occurred_on, '%Y-%m') AS ym,
  CAST(COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS SIGNED) AS income,
  CAST(COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS SIGNED) AS expense
FROM transactions
WHERE occurred_on BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
GROUP BY DATE_FORMAT(occurred_on, '%Y-%m')
ORDER BY ym;
