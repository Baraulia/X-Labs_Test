package grpcserver

import (
	"context"
	"errors"
	"testing"

	"github.com/Baraulia/X-Labs_Test/internal/api/grpc/pb"
	"github.com/Baraulia/X-Labs_Test/internal/api/serviceMocks"
	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/Baraulia/X-Labs_Test/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, dto *models.User)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.ChangeUserRequest
		convertData    *models.User
		expectedResult *pb.UserResponse
		mockBehavior   mockBehavior
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.ChangeUserRequest{
				Email:    "test@gmail.com",
				Username: "testUserName",
				Password: "test",
				Admin:    false,
			},
			convertData: &models.User{
				Email:    "test@gmail.com",
				UserName: "testUserName",
				Password: "test",
				Admin:    false,
			},
			expectedResult: &pb.UserResponse{
				User: &pb.User{
					Id:       newUUID,
					Email:    "test@gmail.com",
					Username: "testUserName",
					Password: "test",
					Admin:    false,
				},
			},
			mockBehavior: func(s *serviceMocks.MockServiceInterface, dto *models.User) {
				s.EXPECT().CreateUser(ctx, dto).Return(&models.User{
					ID:       newUUID,
					Email:    "test@gmail.com",
					UserName: "testUserName",
					Password: "test",
					Admin:    false,
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.ChangeUserRequest{
				Email:    "test&gmail.com",
				Username: "testUserName",
				Password: "test",
				Admin:    false,
			},
			convertData: &models.User{
				Email:    "test&gmail.com",
				UserName: "testUserName",
				Password: "test",
				Admin:    false,
			},
			mockBehavior: func(s *serviceMocks.MockServiceInterface, dto *models.User) {
				s.EXPECT().CreateUser(ctx, dto).Return(nil, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.convertData)
			server := NewServer(service, logg)

			response, err := server.CreateUser(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, response)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, dto models.UpdateUserDTO, id string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	email := "test@gmail.com"

	testTable := []struct {
		name          string
		inputData     *pb.ChangeUserRequest
		convertData   models.UpdateUserDTO
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "successful",
			inputData: &pb.ChangeUserRequest{
				Email: email,
			},
			convertData: models.UpdateUserDTO{
				Email: &email,
			},
			mockBehavior: func(s *serviceMocks.MockServiceInterface, dto models.UpdateUserDTO, id string) {
				s.EXPECT().UpdateUser(ctx, dto, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.ChangeUserRequest{
				Email: email,
			},
			convertData: models.UpdateUserDTO{
				Email: &email,
			},
			mockBehavior: func(s *serviceMocks.MockServiceInterface, dto models.UpdateUserDTO, id string) {
				s.EXPECT().UpdateUser(ctx, dto, id).Return(errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.convertData, testCase.inputData.Id)
			server := NewServer(service, logg)

			_, err := server.UpdateUser(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetOneUserByID(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, id string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.GetUserByIdRequest
		id             string
		mockBehavior   mockBehavior
		expectedResult *pb.UserResponse
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.GetUserByIdRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *serviceMocks.MockServiceInterface, id string) {
				s.EXPECT().GetOneUserByID(ctx, id).Return(&models.User{
					ID:       newUUID,
					Email:    "test@gmail.com",
					UserName: "test username",
					Password: "pass",
					Admin:    false,
				}, nil)
			},
			expectedResult: &pb.UserResponse{User: &pb.User{
				Id:       newUUID,
				Email:    "test@gmail.com",
				Username: "test username",
				Password: "pass",
				Admin:    false,
			}},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.GetUserByIdRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *serviceMocks.MockServiceInterface, id string) {
				s.EXPECT().GetOneUserByID(ctx, id).Return(nil, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.id)
			server := NewServer(service, logg)

			result, err := server.GetOneUserByID(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, result)
			}
		})
	}
}

func TestGetOneUserByUsername(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, username string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.GetUserByUsernameRequest
		username       string
		mockBehavior   mockBehavior
		expectedResult *pb.UserResponse
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.GetUserByUsernameRequest{
				Username: "username",
			},
			username: "username",
			mockBehavior: func(s *serviceMocks.MockServiceInterface, username string) {
				s.EXPECT().GetOneUserByUsername(ctx, username).Return(&models.User{
					ID:       newUUID,
					Email:    "test@gmail.com",
					UserName: "username",
					Password: "pass",
					Admin:    false,
				}, nil)
			},
			expectedResult: &pb.UserResponse{User: &pb.User{
				Id:       newUUID,
				Email:    "test@gmail.com",
				Username: "username",
				Password: "pass",
				Admin:    false,
			}},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.GetUserByUsernameRequest{
				Username: "username",
			},
			username: "username",
			mockBehavior: func(s *serviceMocks.MockServiceInterface, username string) {
				s.EXPECT().GetOneUserByUsername(ctx, username).Return(nil, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.username)
			server := NewServer(service, logg)

			result, err := server.GetOneUserByUsername(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, result)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, id string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.DeleteUserRequest
		id             string
		mockBehavior   mockBehavior
		expectedResult *pb.DeleteUserResponse
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.DeleteUserRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *serviceMocks.MockServiceInterface, id string) {
				s.EXPECT().DeleteUser(ctx, id).Return(nil)
			},
			expectedResult: &pb.DeleteUserResponse{Success: true},
			expectedError:  false,
		},
		{
			name: "error from service",
			inputData: &pb.DeleteUserRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *serviceMocks.MockServiceInterface, id string) {
				s.EXPECT().DeleteUser(ctx, id).Return(errors.New("service error"))
			},
			expectedResult: &pb.DeleteUserResponse{Success: false},
			expectedError:  true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.id)
			server := NewServer(service, logg)

			result, err := server.DeleteUser(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
				require.Equal(t, testCase.expectedResult, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, result)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	type mockBehavior func(s *serviceMocks.MockServiceInterface, offset, limit int)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID1 := uuid.New().String()
	newUUID2 := uuid.New().String()
	totalUsers := 10

	testTable := []struct {
		name           string
		inputData      *pb.GetUsersRequest
		offset         int
		limit          int
		expectedResult *pb.GetUsersResponse
		mockBehavior   mockBehavior
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.GetUsersRequest{
				Offset: 10,
				Limit:  2,
			},
			offset: 10,
			limit:  2,
			expectedResult: &pb.GetUsersResponse{
				Users: []*pb.User{
					{
						Id:       newUUID1,
						Email:    "test@gmail.com",
						Username: "testUserName",
						Password: "test",
						Admin:    false,
					},
					{
						Id:       newUUID2,
						Email:    "test2@gmail.com",
						Username: "testUserName2",
						Password: "test2",
						Admin:    false,
					},
				},
				TotalUsers: int32(totalUsers),
			},
			mockBehavior: func(s *serviceMocks.MockServiceInterface, offset, limit int) {
				s.EXPECT().GetUsers(ctx, offset, limit).Return([]models.User{
					{
						ID:       newUUID1,
						Email:    "test@gmail.com",
						UserName: "testUserName",
						Password: "test",
						Admin:    false,
					},
					{
						ID:       newUUID2,
						Email:    "test2@gmail.com",
						UserName: "testUserName2",
						Password: "test2",
						Admin:    false,
					},
				}, totalUsers, nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.GetUsersRequest{
				Offset: 10,
				Limit:  2,
			},
			offset: 10,
			limit:  2,
			mockBehavior: func(s *serviceMocks.MockServiceInterface, offset, limit int) {
				s.EXPECT().GetUsers(ctx, offset, limit).Return(nil, 0, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := serviceMocks.NewMockServiceInterface(c)
			testCase.mockBehavior(service, testCase.offset, testCase.limit)
			server := NewServer(service, logg)

			response, err := server.GetUsers(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, response)
			}
		})
	}
}
