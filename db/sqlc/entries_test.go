package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/stretchr/testify/require"
)

const (
	YYYYMMDD = "2006-01-02"
)

func GetMadeUpDate(madeUpDate string) (time.Time, error) {
	return time.Parse(YYYYMMDD, madeUpDate)
}

func createRandomEntry(t *testing.T, user User) Entry {
	date, err := GetMadeUpDate("2022-12-11")
	require.NoError(t, err)
	require.NotEmpty(t, date)

	arg := CreateEntryParams{
		Owner:   user.Username,
		Name:    util.RandomString(4),
		DueDate: date,
		Amount:  util.RandomMoney(),
		Category: sql.NullString{
			String: util.RandomString(4),
			Valid:  true,
		},
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.Owner, entry.Owner)
	require.Equal(t, arg.Name, entry.Name)
	// require.Equal(t, arg.DueDate, entry.DueDate)
	require.NotZero(t, entry.DueDate)
	require.Equal(t, arg.Amount, entry.Amount)

	require.Equal(t, arg.Category, entry.Category)

	return entry
}

func TestCreateEntry(t *testing.T) {
	user := createRandomUser(t)
	createRandomEntry(t, user)
}

func TestUpdateEntry(t *testing.T) {
	user := createRandomUser(t)
	entry := createRandomEntry(t, user)

	arg := UpdateEntryParams{
		Owner:  user.Username,
		ID:     entry.ID,
		Amount: 10,
	}

	fmt.Println(entry.Name)

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)

	require.Equal(t, entry.Owner, updatedEntry.Owner)
	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, entry.DueDate, updatedEntry.DueDate)
	require.Equal(t, entry.Name, updatedEntry.Name)
	require.Equal(t, arg.Amount, updatedEntry.Amount)
}

func TestDeleteEntry(t *testing.T) {
	user := createRandomUser(t)
	entry := createRandomEntry(t, user)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	arg := GetEntryParams{
		Owner: user.Username,
		ID:    entry.ID,
	}

	deletedEntry, err := testQueries.GetEntry(context.Background(), arg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deletedEntry)
}

func TestDeleteEntries(t *testing.T) {
	user := createRandomUser(t)
	n := 5
	entries := make([]Entry, n)

	for i := 0; i < n; i++ {
		entries = append(entries, createRandomEntry(t, user))
	}

	err := testQueries.DeleteEntries(context.Background(), user.Username)
	require.NoError(t, err)

	for i := 0; i < n; i++ {
		getEntryParams := GetEntryParams{
			Owner: user.Username,
			ID:    entries[i].ID,
		}
		deletedEntry, err := testQueries.GetEntry(context.Background(), getEntryParams)
		require.Error(t, err)
		require.EqualError(t, err, sql.ErrNoRows.Error())
		require.Empty(t, deletedEntry)
	}
}

func TestGetEntry(t *testing.T) {
	user := createRandomUser(t)
	entry := createRandomEntry(t, user)

	arg := GetEntryParams{
		Owner: user.Username,
		ID:    entry.ID,
	}

	retrievedEntry, err := testQueries.GetEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedEntry)

	require.Equal(t, entry.ID, retrievedEntry.ID)
	require.Equal(t, entry.Name, retrievedEntry.Name)
	require.Equal(t, entry.Owner, retrievedEntry.Owner)
	require.Equal(t, entry.Amount, retrievedEntry.Amount)
	require.Equal(t, entry.DueDate, retrievedEntry.DueDate)
}

func TestGetEntries(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, user)
	}

	entries, err := testQueries.GetEntries(context.Background(), user.Username)
	require.NoError(t, err)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestGetCategories(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, user)
	}

	categories, err := testQueries.GetCategories(context.Background())
	require.NoError(t, err)

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}
