-- name: ListTransactions :many
SELECT
  t.id, t.occurred_on, t.amount, t.type, t.category_id,
  c.name AS category_name,
  t.memo, t.created_at, t.updated_at
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.occurred_on BETWEEN sqlc.arg(start_date) AND sqlc.arg(end_date)
  AND (sqlc.narg(type) IS NULL OR t.type = sqlc.narg(type))
  AND (sqlc.narg(category_id) IS NULL OR t.category_id = sqlc.narg(category_id))
ORDER BY t.occurred_on DESC, t.id DESC;

-- name: GetTransaction :one
SELECT
  t.id, t.occurred_on, t.amount, t.type, t.category_id,
  c.name AS category_name,
  t.memo, t.created_at, t.updated_at
FROM transactions t
JOIN categories c ON c.id = t.category_id
WHERE t.id = ?;

-- name: CreateTransaction :execlastid
INSERT INTO transactions (occurred_on, amount, type, category_id, memo)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateTransaction :execrows
UPDATE transactions
SET occurred_on = ?, amount = ?, type = ?, category_id = ?, memo = ?
WHERE id = ?;

-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = ?;
