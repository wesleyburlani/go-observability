SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = @id;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = @username;

-- name: UpdateUser :one
UPDATE users SET
  username = coalesce(sqlc.narg(username), username),
  email = coalesce(sqlc.narg(email), email),
  password = coalesce(sqlc.narg(password), password)
WHERE id = @id
RETURNING *;

-- name: DeleteUserById :one
DELETE FROM users
WHERE id = @id
RETURNING *;

-- name: DeleteUserByEmail :one
DELETE FROM users
WHERE email = @email
RETURNING *;

-- name: DeleteUserByUsername :one
DELETE FROM users
WHERE username = @username
RETURNING *;
