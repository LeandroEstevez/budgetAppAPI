// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteEntries(ctx context.Context, owner string) error
	DeleteEntry(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, username string) error
	GetCategories(ctx context.Context) ([]sql.NullString, error)
	GetEntries(ctx context.Context, owner string) ([]Entry, error)
	GetEntry(ctx context.Context, arg GetEntryParams) (Entry, error)
	GetEntryForUpdate(ctx context.Context, arg GetEntryForUpdateParams) (Entry, error)
	GetUser(ctx context.Context, username string) (User, error)
	GetUserForUpdate(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
