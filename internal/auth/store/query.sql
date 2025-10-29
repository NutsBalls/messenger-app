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

-- name: UpdateRefreshToken :exec
UPDATE users
SET
    crypted_refresh_token = $2
WHERE id = $1;