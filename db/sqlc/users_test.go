package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword , err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams {
		Username: util.RandomString(6),
		HashedPassword: hashedPassword,
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
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)

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
	userv1 := createRandomUser(t)

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
	newUser := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), newUser.Username)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), newUser.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i ++ {
		createRandomUser(t)
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

