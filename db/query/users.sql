-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email, total_expenses
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;