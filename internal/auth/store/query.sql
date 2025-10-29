-- name: CreateUser :one
INSERT INTO users (
    name, email, password_hash
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: LogIn :one
SELECT * FROM users
WHERE email = $1 AND password_hash = $2;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET 
    name = $2,
    email = $3,
    password_hash = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, refresh_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET
    refresh_token = $2,
    expires_at = $3,
    created_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: CheckRefreshToken :one
SELECT EXISTS (
    SELECT 1
    FROM refresh_tokens
    WHERE refresh_token = $1
      AND revoked = false
      AND expires_at > NOW()
);

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE user_id = $1;

-- name: LogoutRefreshToken :one
UPDATE refresh_tokens
SET revoked = true
WHERE refresh_token = $1
RETURNING *;