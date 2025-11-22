-- name: GetRules :many
SELECT * FROM rules
WHERE username = $1;

-- name: CreateRule :one
INSERT INTO rules (
  username, rule_name, rule_type, rule_description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateRule :exec
UPDATE rules
SET rule_type = $1,
    rule_description = $2
WHERE username = $3 AND rule_name = $4;

-- name: DeleteRule :exec
DELETE FROM rules
WHERE id = $1;