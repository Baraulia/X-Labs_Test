package app

//nolint:depguard
import (
	"context"

	"github.com/Baraulia/X-Labs_Test/internal/models"
)

type App struct {
	logger    Logger
	storage   StorageInterface
	validator Validator
	SecretKey string
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type Validator interface {
	IsEmail(email string) bool
}

//go:generate mockgen -destination mocks/storageMock.go -package mocks github.com/Baraulia/X-Labs_Test/internal/app StorageInterface
type StorageInterface interface {
	CreateUser(ctx context.Context, userDTO *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, userDTO models.UpdateUserDTO, userID string) error
	DeleteUser(ctx context.Context, userID string) error
	GetUsers(ctx context.Context, offset, limit int) ([]models.User, int, error)
	GetOneUserByID(ctx context.Context, userID string) (*models.User, error)
	GetOneUserByUsername(ctx context.Context, userName string) (*models.User, error)
}

func NewApp(logger Logger, storage StorageInterface, validator Validator, secretKey string) *App {
	return &App{logger: logger, storage: storage, validator: validator, SecretKey: secretKey}
}
