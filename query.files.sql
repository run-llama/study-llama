-- name: CreateFile :one
INSERT INTO files (
  username, file_name, file_category
) VALUES (
  $1, $2, $3
)
RETURNING *;