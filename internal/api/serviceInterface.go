package api

//nolint:depguard
import (
	"context"

	"github.com/Baraulia/X-Labs_Test/internal/models"
)

//go:generate mockgen -destination serviceMocks/serviceMock.go -package serviceMocks github.com/Baraulia/X-Labs_Test/internal/api ServiceInterface
type ServiceInterface interface {
	CreateUser(ctx context.Context, userDTO *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, userDTO models.UpdateUserDTO, userID string) error
	DeleteUser(ctx context.Context, userID string) error
	GetUsers(ctx context.Context, offset, limit int) ([]models.User, int, error)
	GetOneUserByID(ctx context.Context, userID string) (*models.User, error)
	GetOneUserByUsername(ctx context.Context, userName string) (*models.User, error)
	CheckPassword(ctx context.Context, username, password string) (bool, error)
}
