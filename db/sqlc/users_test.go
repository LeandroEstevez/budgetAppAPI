package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams {
		Username: util.RandomString(6),
		HashedPassword: util.RandomString(3),
		FullName: util.RandomFullName(),
		Email: util.RandomEmail(),
		TotalExpenses:  util.RandomMoney(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	// TODO: compare hashed password
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.TotalExpenses, user.TotalExpenses)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := CreateRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), newUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, newUser.Username)
	// TODO: compare hashed password
	require.Equal(t, user.FullName, newUser.FullName)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.TotalExpenses, newUser.TotalExpenses)
}

func TestUpdateUser(t *testing.T) {
	userv1 := CreateRandomUser(t)

	arg := UpdateUserParams {
		Username: userv1.Username,
		TotalExpenses: util.RandomMoney(),
	}

	userv2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, userv2)

	require.Equal(t, userv1.Username, userv2.Username)
	// TODO: compare hashed password
	require.Equal(t, userv1.FullName, userv2.FullName)
	require.Equal(t, userv1.Email, userv2.Email)
	require.Equal(t, arg.TotalExpenses, userv2.TotalExpenses)
}

func TestDeleteUser(t *testing.T) {
	newUser := CreateRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), newUser.Username)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), newUser.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i ++ {
		CreateRandomUser(t)
	}

	arg := ListUsersParams {
		Limit: 5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

