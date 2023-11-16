package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := pkg.HashPassword(pkg.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Email:          pkg.RandomEmail(),
		HashedPassword: hashedPassword,
		Name:           pkg.RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Name, user.Name)
	require.NotZero(t, user.CreatedAt)

	return user
}

func Test_create_user(t *testing.T) {
	createRandomUser(t)
}

func Test_get_user(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func Test_list_users(t *testing.T) {
	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func Test_get_user_for_update(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserForUpdate(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Name, user2.Name)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Phone, user2.Phone)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserInfo(t *testing.T) {
	oldUser := createRandomUser(t)

	newName := pkg.RandomString(6)
	userUpdated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID:   oldUser.ID,
		Name: newName,
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Name, userUpdated.Name)
	require.Equal(t, newName, userUpdated.Name)
	require.Equal(t, oldUser.Email, userUpdated.Email)
	require.Equal(t, oldUser.HashedPassword, userUpdated.HashedPassword)
}

func Test_delete_user(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(),user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
