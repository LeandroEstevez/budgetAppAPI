-- name: CreateEntry :one
INSERT INTO entries (
  owner, name, due_date, amount
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetEntries :many
SELECT * FROM entries
WHERE owner = $1;

-- name: GetEntry :one
SELECT * FROM entries
WHERE owner = $1 AND id = $2;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $3
WHERE owner = $1 AND id = $2
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;