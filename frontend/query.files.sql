-- name: GetFiles :many
SELECT * FROM files
WHERE username = $1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;