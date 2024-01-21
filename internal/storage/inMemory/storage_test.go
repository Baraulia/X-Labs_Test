package memorystorage

import (
	"context"
	"testing"

	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/Baraulia/X-Labs_Test/pkg/logger"
	"github.com/stretchr/testify/require"
)

var testUsers = []models.User{
	{
		Email:    "testEmail@gmail.com",
		UserName: "testUserName",
		Password: "testPassword",
		Admin:    false,
	},
	{
		Email:    "testEmail2@gmail.com",
		UserName: "testUserName2",
		Password: "testPassword2",
		Admin:    false,
	},
	{
		Email:    "testEmail3@gmail.com",
		UserName: "testUserName3",
		Password: "testPassword3",
		Admin:    false,
	},
	{
		Email:    "testEmail4@gmail.com",
		UserName: "testUserName4",
		Password: "testPassword4",
		Admin:    false,
	},
	{
		Email:    "testEmail5@gmail.com",
		UserName: "testUserName5",
		Password: "testPassword5",
		Admin:    false,
	},
	{
		Email:    "testEmail6@gmail.com",
		UserName: "testUserName6",
		Password: "testPassword6",
		Admin:    false,
	},
	{
		Email:    "testEmail7@gmail.com",
		UserName: "testUserName7",
		Password: "testPassword7",
		Admin:    false,
	},
	{
		Email:    "testEmail8@gmail.com",
		UserName: "testUserName8",
		Password: "testPassword8",
		Admin:    false,
	},
	{
		Email:    "testEmail9@gmail.com",
		UserName: "testUserName9",
		Password: "testPassword9",
		Admin:    false,
	},
	{
		Email:    "testEmail10@gmail.com",
		UserName: "testUserName10",
		Password: "testPassword10",
		Admin:    false,
	},
}

func TestCreateUser(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := NewUserStorage(logg)
	ctx := context.Background()

	user, err := storage.CreateUser(ctx, &testUsers[0])
	require.NoError(t, err)

	userFromStorage, ok := storage.users[user.ID]
	if !ok {
		t.Error("User not added to the storage")
	}

	require.Equal(t, testUsers[0].UserName, userFromStorage.UserName)
}

func TestUpdateUser(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := NewUserStorage(logg)
	ctx := context.Background()

	user, err := storage.CreateUser(ctx, &testUsers[0])
	require.NoError(t, err)

	oldUserName := user.UserName
	newUsername := "newUsername"
	err = storage.UpdateUser(ctx, models.UpdateUserDTO{
		UserName: &newUsername,
	}, user.ID)
	require.NoError(t, err)

	userFromStorageNew, ok := storage.users[user.ID]
	if !ok {
		t.Error("Event not added to the storage")
	}

	require.NotEqual(t, oldUserName, userFromStorageNew.UserName)
	require.Equal(t, userFromStorageNew.UserName, newUsername)
}

func TestGetUsers(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := NewUserStorage(logg)
	ctx := context.Background()

	for _, user := range testUsers {
		_, err = storage.CreateUser(ctx, &user)
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		offset        int
		limit         int
		expectedCount int
	}{
		{"simple case", 1, 2, 2},
		{"offset is bigger than count users", 11, 2, 0},
		{"limit out of range", 9, 5, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			users, count, err := storage.GetUsers(ctx, test.offset, test.limit)
			require.NoError(t, err)
			require.Equal(t, test.expectedCount, len(users))
			require.Equal(t, count, len(testUsers))
		})
	}
}

func TestGetUserById(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := NewUserStorage(logg)
	ctx := context.Background()

	user, err := storage.CreateUser(ctx, &testUsers[0])
	require.NoError(t, err)

	userFromStorage, ok := storage.users[user.ID]
	if !ok {
		t.Error("User not added to the storage")
	}

	userByID, err := storage.GetOneUserByID(ctx, user.ID)
	require.NoError(t, err)

	require.Equal(t, userByID.UserName, userFromStorage.UserName)
}

func TestGetUserByUserName(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := NewUserStorage(logg)
	ctx := context.Background()

	user, err := storage.CreateUser(ctx, &testUsers[0])
	require.NoError(t, err)

	userFromStorage, ok := storage.users[user.ID]
	if !ok {
		t.Error("User not added to the storage")
	}

	userByUserName, err := storage.GetOneUserByUsername(ctx, user.UserName)
	require.NoError(t, err)

	require.Equal(t, userByUserName.UserName, userFromStorage.UserName)
}
