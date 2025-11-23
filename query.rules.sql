-- name: GetRules :many
SELECT * FROM rules
WHERE username = $1;