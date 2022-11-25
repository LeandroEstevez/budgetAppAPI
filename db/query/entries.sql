-- name: CreateEntry :one
INSERT INTO entries (
  owner, name, due_date, amount, category
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetEntries :many
SELECT * FROM entries
WHERE owner = $1;

-- name: GetEntry :one
SELECT * FROM entries
WHERE owner = $1 AND id = $2;

-- name: GetEntryForUpdate :one
SELECT * FROM entries
WHERE owner = $1 AND id = $2
FOR NO KEY UPDATE;

-- name: GetCategories :many
SELECT category FROM entries
WHERE category != '' AND category IS NOT NULL;

-- name: UpdateEntry :one
UPDATE entries
SET name = $3, due_date = $4, amount = $5, category = $6
WHERE owner = $1 AND id = $2
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;

-- name: DeleteEntries :exec
DELETE FROM entries
WHERE owner = $1;