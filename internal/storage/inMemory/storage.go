package memorystorage

//nolint:depguard
import (
	"context"
	"fmt"
	"sync"

	"github.com/Baraulia/X-Labs_Test/internal/app"
	"github.com/Baraulia/X-Labs_Test/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	mu              sync.RWMutex
	users           map[string]*models.User
	indexByEmail    map[string]string
	indexByUsername map[string]string
	listIds         []string
	logger          app.Logger
}

func NewUserStorage(logger app.Logger) *UserStorage {
	return &UserStorage{
		users:           make(map[string]*models.User),
		indexByEmail:    make(map[string]string),
		indexByUsername: make(map[string]string),
		logger:          logger,
	}
}

func (us *UserStorage) InitAdmin(initAdminName, initAdminPassword, secretKey string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initAdminPassword+secretKey), bcrypt.DefaultCost)
	if err != nil {
		us.logger.Error("Error generating hash", map[string]interface{}{"error": err})
		return fmt.Errorf("error while generate hash: %w", err)
	}
	newUUID := uuid.New().String()

	us.users[newUUID] = &models.User{
		Email:    "admin@gmail.com",
		UserName: initAdminName,
		Password: string(hashedPassword),
		Admin:    true,
	}

	us.listIds = append(us.listIds, newUUID)
	us.indexByEmail["admin@gmail.com"] = newUUID
	us.indexByUsername[initAdminName] = newUUID

	return nil
}

func (us *UserStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if existingID, exists := us.indexByEmail[user.Email]; exists {
		us.logger.Error(
			"user with a such email already exists", map[string]interface{}{"email": user.Email, "id": existingID})
		return nil, fmt.Errorf("user with email %s already exists (ID: %s)", user.Email, existingID)
	}

	if existingID, exists := us.indexByUsername[user.UserName]; exists {
		us.logger.Error(
			"user with a such username already exists", map[string]interface{}{"username": user.UserName, "id": existingID})
		return nil, fmt.Errorf("user with username %s already exists (ID: %s)", user.UserName, existingID)
	}

	newUUID := uuid.New().String()
	us.users[newUUID] = user
	us.indexByEmail[user.Email] = newUUID
	us.indexByUsername[user.UserName] = newUUID
	us.listIds = append(us.listIds, newUUID)
	user.ID = newUUID

	return user, nil
}

func (us *UserStorage) DeleteUser(ctx context.Context, userID string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	user, exists := us.users[userID]
	if !exists {
		us.logger.Error("user with a such ID does not exist", map[string]interface{}{"id": userID})
		return fmt.Errorf("user with ID %s not found", userID)
	}

	delete(us.users, userID)
	delete(us.indexByEmail, user.Email)
	delete(us.indexByUsername, user.UserName)

	if err := us.removeID(userID); err != nil {
		return err
	}
	us.logger.Info("user was deleted", nil)

	return nil
}

func (us *UserStorage) UpdateUser(ctx context.Context, userDTO models.UpdateUserDTO, userID string) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	user, exists := us.users[userID]
	if !exists {
		us.logger.Error("user with a such ID does not exist", map[string]interface{}{"id": userID})
		return fmt.Errorf("user with ID %s not found", userID)
	}

	if userDTO.UserName != nil {
		if existingID, ex := us.indexByUsername[*userDTO.UserName]; ex {
			us.logger.Error(
				"user with a such username already exists", map[string]interface{}{"username": user.UserName, "id": existingID})
			return fmt.Errorf("user with username %s already exists (ID: %s)", user.UserName, existingID)
		}

		delete(us.indexByUsername, user.UserName)
		us.indexByUsername[*userDTO.UserName] = userID

		user.UserName = *userDTO.UserName
	}

	if userDTO.Email != nil {
		if existingID, ex := us.indexByEmail[user.Email]; ex {
			us.logger.Error(
				"user with a such email already exists", map[string]interface{}{"email": user.Email, "id": existingID})
			return fmt.Errorf("user with email %s already exists (ID: %s)", user.Email, existingID)
		}

		delete(us.indexByEmail, user.Email)
		us.indexByEmail[*userDTO.Email] = userID
		user.Email = *userDTO.Email
	}

	if userDTO.Password != nil {
		user.Password = *userDTO.Password
	}

	us.users[userID] = user
	us.logger.Info("user was updated", nil)

	return nil
}

func (us *UserStorage) GetUsers(ctx context.Context, offset, limit int) ([]models.User, int, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()
	var start, finish int

	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	count := len(us.listIds)

	if offset >= count {
		return make([]models.User, 0), count, nil
	}

	if limit+offset >= count {
		start = offset
		finish = count
	} else {
		start = offset
		finish = offset + limit
	}

	listUsers := make([]models.User, 0)
	for _, id := range us.listIds[start:finish] {
		listUsers = append(listUsers, *us.users[id])
	}

	return listUsers, count, nil
}

func (us *UserStorage) GetOneUserByID(ctx context.Context, userID string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return us.users[userID], nil
}

func (us *UserStorage) GetOneUserByUsername(ctx context.Context, userName string) (*models.User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	existingID, exists := us.indexByUsername[userName]
	if !exists {
		us.logger.Error("user with a such username does not exist", map[string]interface{}{"username": userName})
		return nil, fmt.Errorf("user with username %s does not exist", userName)
	}

	return us.users[existingID], nil
}

func (us *UserStorage) removeID(targetID string) error {
	index := -1
	for i, id := range us.listIds {
		if id == targetID {
			index = i
			break
		}
	}

	if index != -1 {
		us.listIds = append(us.listIds[:index], us.listIds[index+1:]...)
	} else {
		us.logger.Error("there is not such id in listIds", map[string]interface{}{"id": targetID})
		return fmt.Errorf("there is not such id: %s", targetID)
	}

	return nil
}
