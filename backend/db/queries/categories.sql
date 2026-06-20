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
