-- name: GetToolBySlug :one
SELECT * FROM tools WHERE slug = ? LIMIT 1;

-- name: GetToolByID :one
SELECT * FROM tools WHERE id = ? LIMIT 1;

-- name: CreateTool :one
INSERT INTO tools (name, slug, category, sub_type, prolang, release_year, devstatus, details, use_cases, tags, website, github) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateTool :one
UPDATE tools SET name = ?, slug = ?, category = ?, sub_type = ?, prolang = ?, release_year = ?, devstatus = ?, details = ?, use_cases = ?, tags = ?, website = ?, github = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteTool :exec
DELETE FROM tools WHERE slug = ?;

-- name: GetToolsByIDs :many
SELECT * FROM tools WHERE id IN (sqlc.slice('ids')) ORDER BY name ASC;