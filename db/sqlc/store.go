package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/golang/mock/mockgen/model"
)

// Store gives all functions to execute db queries and transactions
type Store interface {
	Querier
	AddEntryTx(ctx context.Context, arg AddEntryTxParams) (AddEntryTxResult, error)
	DeleteEntryTx(ctx context.Context, arg DeleteEntryTxParams) (DeleteEntryTxResult, error)
	UpdateEntryTx(ctx context.Context, arg UpdateEntryTxParams) (UpdateEntryTxResult, error)
	DeleteUserTx(ctx context.Context, username string) error
	UpdateAccountTx(ctx context.Context, arg UpdateAccountTxParams) (UpdateAccountTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	fmt.Print(err)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

// Contains the input parameter of the add entry transaction
type AddEntryTxParams struct {
	Username string    `json:"username"`
	Name     string    `json:"name"`
	DueDate  time.Time `json:"due_date"`
	Amount   int64     `json:"amount"`
	Category string    `json:"category"`
}

// Contains the result of the Add entry transaction
type AddEntryTxResult struct {
	Entry Entry `json:"entry"`
	User  User  `json:"user"`
}

// Adds an entry and updates the total expense in the user
func (store *SQLStore) AddEntryTx(ctx context.Context, arg AddEntryTxParams) (AddEntryTxResult, error) {
	var result AddEntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}

		createEntryParams := CreateEntryParams{
			Owner:    arg.Username,
			Name:     arg.Name,
			DueDate:  arg.DueDate,
			Amount:   arg.Amount,
			Category: sql.NullString{String: arg.Category, Valid: true},
		}
		result.Entry, err = q.CreateEntry(ctx, createEntryParams)
		if err != nil {
			return err
		}

		totalExpenses := user.TotalExpenses
		entryAmount := result.Entry.Amount
		totalExpenses = totalExpenses + entryAmount

		updatedUserParams := UpdateUserParams{
			Username:      arg.Username,
			TotalExpenses: totalExpenses,
		}
		result.User, err = q.UpdateUser(ctx, updatedUserParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// Contains the input parameter of the update entry transaction
type UpdateEntryTxParams struct {
	Username string    `json:"username"`
	ID       int32     `json:"id"`
	Name     string    `json:"name"`
	DueDate  time.Time `json:"due_date"`
	Amount   int64     `json:"amount"`
	Category string    `json:"category"`
}

// Contains the result of the update entry transaction
type UpdateEntryTxResult struct {
	Entry Entry `json:"entry"`
	User  User  `json:"user"`
}

// Updates the amount of an entry and updates the total expense in the user
func (store *SQLStore) UpdateEntryTx(ctx context.Context, arg UpdateEntryTxParams) (UpdateEntryTxResult, error) {
	var result UpdateEntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}

		getEntryParams := GetEntryForUpdateParams{
			Owner: arg.Username,
			ID:    arg.ID,
		}
		entry, err := q.GetEntryForUpdate(ctx, getEntryParams)
		if err != nil {
			return err
		}

		totalExpenses := user.TotalExpenses
		entryAmount := entry.Amount

		changeInAmount := arg.Amount - entryAmount
		totalExpenses = totalExpenses + changeInAmount

		updatedUserParams := UpdateUserParams{
			Username:      arg.Username,
			TotalExpenses: totalExpenses,
		}
		result.User, err = q.UpdateUser(ctx, updatedUserParams)
		if err != nil {
			return err
		}

		var categoryValue sql.NullString
		if arg.Category == "" {
			categoryValue = sql.NullString{
				String: "",
				Valid:  false,
			}
		} else {
			categoryValue = sql.NullString{
				String: arg.Category,
				Valid:  true,
			}
		}

		updateEntryParams := UpdateEntryParams{
			Owner:    arg.Username,
			ID:       arg.ID,
			Name:     arg.Name,
			DueDate:  arg.DueDate,
			Amount:   arg.Amount,
			Category: categoryValue,
		}

		result.Entry, err = q.UpdateEntry(ctx, updateEntryParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// Contains the input parameter of the delete entry transaction
type DeleteEntryTxParams struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

// Contains the result of the delete entry transaction
type DeleteEntryTxResult struct {
	User User `json:"user"`
}

// Updates the amount of an entry and updates the total expense in the user
func (store *SQLStore) DeleteEntryTx(ctx context.Context, arg DeleteEntryTxParams) (DeleteEntryTxResult, error) {
	var result DeleteEntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.Username)
		if err != nil {
			return err
		}

		getEntryParams := GetEntryParams{
			Owner: arg.Username,
			ID:    arg.ID,
		}
		entry, err := q.GetEntry(ctx, getEntryParams)
		if err != nil {
			return err
		}

		totalExpenses := user.TotalExpenses
		entryAmount := entry.Amount
		totalExpenses = totalExpenses - entryAmount

		err = q.DeleteEntry(ctx, arg.ID)
		if err != nil {
			return err
		}

		updatedUserParams := UpdateUserParams{
			Username:      arg.Username,
			TotalExpenses: totalExpenses,
		}
		result.User, err = q.UpdateUser(ctx, updatedUserParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// Updates the amount of an user and updates the total expense in the user
func (store *SQLStore) DeleteUserTx(ctx context.Context, username string) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.DeleteEntries(ctx, username)
		if err != nil {
			return err
		}

		err = q.DeleteUser(ctx, username)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Contains the input parameter of the update entry transaction
type UpdateAccountTxParams struct {
	OrigUsername string `json:"origusername"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
}

// Contains the result of the update entry transaction
type UpdateAccountTxResult struct {
	User User `json:"user"`
}

// Updates the amount of an entry and updates the total expense in the user
func (store *SQLStore) UpdateAccountTx(ctx context.Context, arg UpdateAccountTxParams) (UpdateAccountTxResult, error) {
	var result UpdateAccountTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		updateEntryParams := UpdateEntriesOwnerParams{
			Owner:   arg.OrigUsername,
			Owner_2: arg.Username,
		}
		err = q.UpdateEntriesOwner(ctx, updateEntryParams)
		if err != nil {
			return err
		}

		updatedUserParams := UpdateUserInfoParams{
			Username:   arg.OrigUsername,
			Username_2: arg.Username,
			FullName:   arg.FullName,
			Email:      arg.Email,
		}
		result.User, err = q.UpdateUserInfo(ctx, updatedUserParams)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
