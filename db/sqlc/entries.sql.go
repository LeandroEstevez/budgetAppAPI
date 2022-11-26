// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: entries.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
  owner, name, due_date, amount, category
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, owner, name, due_date, amount, category
`

type CreateEntryParams struct {
	Owner    string         `json:"owner"`
	Name     string         `json:"name"`
	DueDate  time.Time      `json:"due_date"`
	Amount   int64          `json:"amount"`
	Category sql.NullString `json:"category"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry,
		arg.Owner,
		arg.Name,
		arg.DueDate,
		arg.Amount,
		arg.Category,
	)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.DueDate,
		&i.Amount,
		&i.Category,
	)
	return i, err
}

const deleteEntries = `-- name: DeleteEntries :exec
DELETE FROM entries
WHERE owner = $1
`

func (q *Queries) DeleteEntries(ctx context.Context, owner string) error {
	_, err := q.db.ExecContext(ctx, deleteEntries, owner)
	return err
}

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1
`

func (q *Queries) DeleteEntry(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteEntry, id)
	return err
}

const getCategories = `-- name: GetCategories :many
SELECT category FROM entries
WHERE category != '' AND category IS NOT NULL
`

func (q *Queries) GetCategories(ctx context.Context) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []sql.NullString{}
	for rows.Next() {
		var category sql.NullString
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntries = `-- name: GetEntries :many
SELECT id, owner, name, due_date, amount, category FROM entries
WHERE owner = $1
`

func (q *Queries) GetEntries(ctx context.Context, owner string) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntries, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Name,
			&i.DueDate,
			&i.Amount,
			&i.Category,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEntry = `-- name: GetEntry :one
SELECT id, owner, name, due_date, amount, category FROM entries
WHERE owner = $1 AND id = $2
`

type GetEntryParams struct {
	Owner string `json:"owner"`
	ID    int32  `json:"id"`
}

func (q *Queries) GetEntry(ctx context.Context, arg GetEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, arg.Owner, arg.ID)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.DueDate,
		&i.Amount,
		&i.Category,
	)
	return i, err
}

const getEntryForUpdate = `-- name: GetEntryForUpdate :one
SELECT id, owner, name, due_date, amount, category FROM entries
WHERE owner = $1 AND id = $2
FOR NO KEY UPDATE
`

type GetEntryForUpdateParams struct {
	Owner string `json:"owner"`
	ID    int32  `json:"id"`
}

func (q *Queries) GetEntryForUpdate(ctx context.Context, arg GetEntryForUpdateParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntryForUpdate, arg.Owner, arg.ID)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.DueDate,
		&i.Amount,
		&i.Category,
	)
	return i, err
}

const updateEntry = `-- name: UpdateEntry :one
UPDATE entries
SET name = $3, due_date = $4, amount = $5, category = $6
WHERE owner = $1 AND id = $2
RETURNING id, owner, name, due_date, amount, category
`

type UpdateEntryParams struct {
	Owner    string         `json:"owner"`
	ID       int32          `json:"id"`
	Name     string         `json:"name"`
	DueDate  time.Time      `json:"due_date"`
	Amount   int64          `json:"amount"`
	Category sql.NullString `json:"category"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, updateEntry,
		arg.Owner,
		arg.ID,
		arg.Name,
		arg.DueDate,
		arg.Amount,
		arg.Category,
	)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Name,
		&i.DueDate,
		&i.Amount,
		&i.Category,
	)
	return i, err
}
