package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/stretchr/testify/require"
)

func TestAddEntryTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)

	date, err := GetMadeUpDate("2022-12-11")
	require.NoError(t, err)
	require.NotEmpty(t, date)

	addEntryTxParams :=  AddEntryTxParams {
		Username: user.Username,
		Name: util.RandomString(6),
		DueDate: date,
		Amount: util.RandomMoney(),
	}

	result, err := store.AddEntryTx(context.Background(), addEntryTxParams)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	totalExpenses := user.TotalExpenses
	entryAmount := addEntryTxParams.Amount
	totalExpenses = totalExpenses + entryAmount

	require.NotEmpty(t, result.User)
	require.Equal(t, totalExpenses, result.User.TotalExpenses)

	require.NotEmpty(t, result.Entry)
	require.Equal(t, entryAmount, result.Entry.Amount)
}

func TestUpdateEntryTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)
	entry := createRandomEntry(t, user)

	amount := int64(10)

	result, err := store.UpdateEntryTx(context.Background(), UpdateEntryTxParams {
		Username: user.Username,
		ID: entry.ID,
		Amount: amount,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	totalExpenses := user.TotalExpenses
	entryAmount := entry.Amount

	changeInAmount := amount - entryAmount
	totalExpenses = totalExpenses + changeInAmount

	require.NotEmpty(t, result.User)
	require.Equal(t, totalExpenses, result.User.TotalExpenses)

	require.NotEmpty(t, result.Entry)
	require.Equal(t, amount, result.Entry.Amount)
}

func TestDeleteEntryTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)
	entry := createRandomEntry(t, user)

	result, err := store.DeleteEntryTx(context.Background(), DeleteEntryTxParams {
		Username: user.Username,
		ID: entry.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	totalExpenses := user.TotalExpenses
	entryAmount := entry.Amount
	totalExpenses = totalExpenses - entryAmount

	require.NotEmpty(t, result.User)
	require.Equal(t, totalExpenses, result.User.TotalExpenses)
}

func TestDeleteUserTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)
	n := 5
	entries := make([]Entry, n)

	for i := 0; i < n; i++ {
		entries = append(entries, createRandomEntry(t, user))
	}

	err := store.DeleteUserTx(context.Background(), user.Username)
	require.NoError(t, err)

	deletedUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedUser)

	for i := 0; i < n; i++ {
		getEntryParams := GetEntryParams {
			Owner: user.Username,
			ID: entries[i].ID,
		}
		deletedEntry, err := testQueries.GetEntry(context.Background(), getEntryParams)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, deletedEntry)
	}
}

func TestConcurrentAddEntryTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)

	fmt.Println("Before >>", user.TotalExpenses)

	date, err := GetMadeUpDate("2022-12-11")
	require.NoError(t, err)
	require.NotEmpty(t, date)

	// run n concurrent additions
	n := 10

	errs := make(chan error)
	results := make(chan AddEntryTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.AddEntryTx(context.Background(), AddEntryTxParams {
				Username: user.Username,
				Name: util.RandomString(6),
				DueDate: date,
				Amount: int64(10),
			})

			errs <- err
			results <- result
		}()
	}

	totalExpenses := user.TotalExpenses
	// check results
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		entryAmount := int64(10)
		totalExpenses = totalExpenses + entryAmount

		fmt.Println("After >>", result.User.TotalExpenses)

		require.NotEmpty(t, result.User)
		// require.Equal(t, totalExpenses, result.User.TotalExpenses)

		require.NotEmpty(t, result.Entry)
		require.Equal(t, entryAmount, result.Entry.Amount)
	}

	updatedUser, err := store.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, totalExpenses, updatedUser.TotalExpenses)
}

func TestConcurrentUpdateEntryTx(t *testing.T) {
	store := NewStore(testDB)

	user := CreateRandomUser(t)
	entry := createRandomEntry(t, user)
	fmt.Println("Before >>", user.TotalExpenses)
	fmt.Println("Before >>", entry.Amount)

	amount := int64(0)

	// run n concurrent updates
	n := 10

	errs := make(chan error)
	results := make(chan UpdateEntryTxResult)

	for i := 0; i < n; i++ {
		amount += 10
		go func() {
			result, err := store.UpdateEntryTx(context.Background(), UpdateEntryTxParams {
				Username: user.Username,
				ID: entry.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		totalExpenses := user.TotalExpenses
		entryAmount := entry.Amount

		changeInAmount := amount - entryAmount
		totalExpenses = totalExpenses + changeInAmount

		fmt.Println("After >>", result.User.TotalExpenses)
		fmt.Println("After >>",  result.Entry.Amount)

		require.NotEmpty(t, result.User)
		require.Equal(t, totalExpenses, result.User.TotalExpenses)

		require.NotEmpty(t, result.Entry)
		require.Equal(t, amount, result.Entry.Amount)
	}
}