-- name: CreateRelationship :one
INSERT INTO relationships (from_tool_id, to_tool_id, kind, metadata) VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateRelationship :one
UPDATE relationships SET from_tool_id = ?, to_tool_id = ?, kind = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteRelationship :exec
DELETE FROM relationships WHERE id = ?;

-- name: GetRelationshipsByToolID :many
SELECT * FROM relationships WHERE from_tool_id = ? OR to_tool_id = ?;

-- name: GetRelationshipByID :one
SELECT * FROM relationships WHERE id = ? LIMIT 1;

-- name: GetToolAlternatives :many
SELECT * FROM relationships WHERE (from_tool_id = ? OR to_tool_id = ?) AND kind = 'alternative_to';