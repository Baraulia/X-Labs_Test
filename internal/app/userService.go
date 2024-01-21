package app

//nolint:depguard
import (
	"context"
	"errors"
	"fmt"

	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (a *App) CreateUser(ctx context.Context, userDTO *models.User) (*models.User, error) {
	if len(userDTO.UserName) == 0 {
		a.logger.Error("empty username", nil)
		return nil, errors.New("empty username")
	}

	if valid := a.validator.IsEmail(userDTO.Email); !valid {
		a.logger.Error("invalid email", map[string]interface{}{"email": userDTO.Email})
		return nil, fmt.Errorf("invalid email: %s", userDTO.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password+a.SecretKey), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error("Error generating hash", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("error while generate hash: %w", err)
	}

	userDTO.Password = string(hashedPassword)

	return a.storage.CreateUser(ctx, userDTO)
}

func (a *App) UpdateUser(ctx context.Context, userDTO models.UpdateUserDTO, userID string) error {
	_, err := uuid.Parse(userID)
	if err != nil {
		a.logger.Error("invalid id(not UUID)", map[string]interface{}{"id": userID})
		return fmt.Errorf("invalid id(not UUID: %s", userID)
	}

	if userDTO.Email != nil {
		valid := a.validator.IsEmail(*userDTO.Email)
		if !valid {
			a.logger.Error("invalid email", map[string]interface{}{"email": *userDTO.Email})
			return fmt.Errorf("invalid email: %s", *userDTO.Email)
		}
	}

	if userDTO.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*userDTO.Password+a.SecretKey), bcrypt.DefaultCost)
		if err != nil {
			a.logger.Error("Error generating hash", map[string]interface{}{"error": err})
			return fmt.Errorf("error while generate hash: %w", err)
		}

		*userDTO.Password = string(hashedPassword)
	}

	return a.storage.UpdateUser(ctx, userDTO, userID)
}

func (a *App) DeleteUser(ctx context.Context, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		a.logger.Error("invalid id(not UUID)", map[string]interface{}{"id": id})
		return fmt.Errorf("invalid id(not UUID: %s", id)
	}

	return a.storage.DeleteUser(ctx, id)
}

func (a *App) GetUsers(ctx context.Context, offset, limit int) ([]models.User, int, error) {
	return a.storage.GetUsers(ctx, offset, limit)
}

func (a *App) GetOneUserByID(ctx context.Context, userID string) (*models.User, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		a.logger.Error("invalid id(not UUID)", map[string]interface{}{"id": userID})
		return nil, fmt.Errorf("invalid id(not UUID: %s", userID)
	}

	return a.storage.GetOneUserByID(ctx, userID)
}

func (a *App) GetOneUserByUsername(ctx context.Context, userName string) (*models.User, error) {
	return a.storage.GetOneUserByUsername(ctx, userName)
}

func (a *App) CheckPassword(ctx context.Context, userName, password string) (bool, error) {
	user, err := a.GetOneUserByUsername(ctx, userName)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+a.SecretKey))
	if err != nil {
		switch errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		case true:
			return false, errors.New("password does not match the hash")
		default:
			return false, err
		}
	}

	return true, nil
}
