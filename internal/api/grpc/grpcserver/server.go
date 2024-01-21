package grpcserver

//nolint:depguard
import (
	"context"

	"github.com/Baraulia/X-Labs_Test/internal/api"
	"github.com/Baraulia/X-Labs_Test/internal/api/grpc/pb"
	"github.com/Baraulia/X-Labs_Test/internal/app"
	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
)

type Server struct {
	service api.ServiceInterface
	logger  app.Logger
	pb.UnimplementedUserServiceServer
}

func NewServer(service api.ServiceInterface, logger app.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	result, err := s.service.CreateUser(ctx, &models.User{
		Email:    req.Email,
		UserName: req.Username,
		Password: req.Password,
		Admin:    req.Admin,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{User: convert(*result)}, nil
}

func (s Server) UpdateUser(ctx context.Context, req *pb.ChangeUserRequest) (*empty.Empty, error) {
	var updateUserDTO models.UpdateUserDTO

	if req.Email != "" {
		updateUserDTO.Email = &req.Email
	}

	if req.Username != "" {
		updateUserDTO.UserName = &req.Username
	}

	if req.Password != "" {
		updateUserDTO.Password = &req.Password
	}

	err := s.service.UpdateUser(ctx, updateUserDTO, req.Id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, count, err := s.service.GetUsers(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
	}

	pbUsers := make([]*pb.User, 0, len(users))

	for _, user := range users {
		pbUsers = append(pbUsers, convert(user))
	}

	return &pb.GetUsersResponse{
		Users:      pbUsers,
		TotalUsers: int32(count),
	}, nil
}

func (s Server) GetOneUserByID(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetOneUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{User: convert(*user)}, nil
}

func (s Server) GetOneUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetOneUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{User: convert(*user)}, nil
}

func convert(user models.User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.UserName,
		Password: user.Password,
		Admin:    user.Admin,
	}
}
