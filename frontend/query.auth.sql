-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUserTokensLogin :exec
UPDATE users
SET session_token = $1,
    csrf_token = $2
WHERE username = $3;

-- name: UpdateUserTokensLogout :exec
UPDATE users
SET session_token = NULL,
    csrf_token = NULL
WHERE session_token = $1 AND csrf_token = $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

-- name: GetUserBySessionToken :one
SELECT * FROM users
WHERE session_token = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password
) VALUES (
  $1, $2
)
RETURNING *;
