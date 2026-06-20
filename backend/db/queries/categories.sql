-- name: ListCategories :many
SELECT id, name, type, created_at, updated_at
FROM categories
ORDER BY type, name;
