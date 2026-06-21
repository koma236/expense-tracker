-- name: ListCategories :many
SELECT id, name, type, created_at, updated_at
FROM categories
ORDER BY type, name;

-- name: ListCategoriesByType :many
SELECT id, name, type, created_at, updated_at
FROM categories
WHERE (sqlc.narg(type) IS NULL OR type = sqlc.narg(type))
ORDER BY type, name;

-- name: GetCategory :one
SELECT id, name, type, created_at, updated_at
FROM categories
WHERE id = ?;

-- name: CreateCategory :execlastid
INSERT INTO categories (name, type)
VALUES (?, ?);

-- name: UpdateCategory :execrows
UPDATE categories
SET name = ?, type = ?
WHERE id = ?;

-- name: DeleteCategory :execrows
DELETE FROM categories
WHERE id = ?;

-- name: CountCategoryByNameType :one
SELECT COUNT(*) AS cnt
FROM categories
WHERE name = ? AND type = ?;

-- name: CountCategoryByNameTypeExcludingID :one
SELECT COUNT(*) AS cnt
FROM categories
WHERE name = ? AND type = ? AND id <> ?;

-- name: CountCategoryUsage :one
SELECT
  (SELECT COUNT(*) FROM transactions t WHERE t.category_id = ?) +
  (SELECT COUNT(*) FROM budgets b WHERE b.category_id = ?) AS usage_count;
