package app

import (
	"context"
	"testing"

	"github.com/Baraulia/X-Labs_Test/internal/app/mocks"
	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/Baraulia/X-Labs_Test/pkg/logger"
	"github.com/Baraulia/X-Labs_Test/pkg/validation"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(s *mocks.MockStorageInterface, dto *models.User)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()

	validator := validation.New()

	testTable := []struct {
		name          string
		inputData     models.User
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "successful",
			inputData: models.User{
				Email:    "test@gmail.com",
				UserName: "testUserName",
				Password: "test",
				Admin:    false,
			},
			mockBehavior: func(s *mocks.MockStorageInterface, dto *models.User) {
				s.EXPECT().CreateUser(ctx, dto).Return(&models.User{
					ID:       uuid.New().String(),
					Email:    "test@gmail.com",
					UserName: "testUserName",
					Password: "test",
					Admin:    false,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "invalid email",
			inputData: models.User{
				Email:    "test&gmail.com",
				UserName: "testUserName",
				Password: "test",
				Admin:    false,
			},
			mockBehavior:  func(s *mocks.MockStorageInterface, dto *models.User) {},
			expectedError: true,
		},
		{
			name: "empty username",
			inputData: models.User{
				Email:    "test@gmail.com",
				Password: "test",
				Admin:    false,
			},
			mockBehavior:  func(s *mocks.MockStorageInterface, dto *models.User) {},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			storage := mocks.NewMockStorageInterface(c)
			testCase.mockBehavior(storage, &testCase.inputData)
			app := NewApp(logg, storage, validator, "")

			_, err = app.CreateUser(ctx, &testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type mockBehavior func(s *mocks.MockStorageInterface, dto models.UpdateUserDTO, userId string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	newUUID := uuid.New().String()
	validator := validation.New()
	ctx := context.Background()

	validEmail := "test@gmail.com"
	invalidEmail := "test&gmail.com"

	testTable := []struct {
		name          string
		inputData     models.UpdateUserDTO
		inputID       string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "successful",
			inputData: models.UpdateUserDTO{
				Email: &validEmail,
			},
			inputID: newUUID,
			mockBehavior: func(s *mocks.MockStorageInterface, dto models.UpdateUserDTO, userId string) {
				s.EXPECT().UpdateUser(ctx, dto, userId).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "invalid email",
			inputData: models.UpdateUserDTO{
				Email: &invalidEmail,
			},
			inputID:       newUUID,
			mockBehavior:  func(s *mocks.MockStorageInterface, dto models.UpdateUserDTO, userId string) {},
			expectedError: true,
		},
		{
			name: "invalid id",
			inputData: models.UpdateUserDTO{
				Email: &validEmail,
			},
			inputID:       "invalid uuid",
			mockBehavior:  func(s *mocks.MockStorageInterface, dto models.UpdateUserDTO, userId string) {},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			storage := mocks.NewMockStorageInterface(c)
			testCase.mockBehavior(storage, testCase.inputData, testCase.inputID)
			app := NewApp(logg, storage, validator, "")

			err = app.UpdateUser(ctx, testCase.inputData, testCase.inputID)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetOneUserById(t *testing.T) {
	type mockBehavior func(s *mocks.MockStorageInterface, userId string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	newUUID := uuid.New().String()
	validator := validation.New()
	ctx := context.Background()

	testTable := []struct {
		name          string
		inputID       string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name:    "successful",
			inputID: newUUID,
			mockBehavior: func(s *mocks.MockStorageInterface, userId string) {
				s.EXPECT().GetOneUserByID(ctx, userId).Return(&models.User{
					ID:       uuid.New().String(),
					Email:    "test@gmail.com",
					UserName: "testUserName",
					Password: "test",
					Admin:    false,
				}, nil)
			},
			expectedError: false,
		},
		{
			name:          "invalid id",
			inputID:       "invalid uuid",
			mockBehavior:  func(s *mocks.MockStorageInterface, userId string) {},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			storage := mocks.NewMockStorageInterface(c)
			testCase.mockBehavior(storage, testCase.inputID)
			app := NewApp(logg, storage, validator, "")

			_, err = app.GetOneUserByID(ctx, testCase.inputID)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type mockBehavior func(s *mocks.MockStorageInterface, userId string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	newUUID := uuid.New().String()
	validator := validation.New()
	ctx := context.Background()

	testTable := []struct {
		name          string
		inputID       string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name:    "successful",
			inputID: newUUID,
			mockBehavior: func(s *mocks.MockStorageInterface, userId string) {
				s.EXPECT().DeleteUser(ctx, userId).Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "invalid id",
			inputID:       "invalid uuid",
			mockBehavior:  func(s *mocks.MockStorageInterface, userId string) {},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			storage := mocks.NewMockStorageInterface(c)
			testCase.mockBehavior(storage, testCase.inputID)
			app := NewApp(logg, storage, validator, "")

			err = app.DeleteUser(ctx, testCase.inputID)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
