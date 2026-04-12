-- name: CreateRelationship :one
INSERT INTO relationships (from_tool_id, to_tool_id, kind, metadata) VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdateRelationship :one
UPDATE relationships SET from_tool_id = ?, to_tool_id = ?, kind = ?, metadata = ? WHERE id = ? RETURNING *;

-- name: DeleteRelationship :exec
DELETE FROM relationships WHERE id = ?;

-- name: GetRelationshipsByToolID :many
SELECT * FROM relationships WHERE from_tool_id = ? OR to_tool_id = ?;
