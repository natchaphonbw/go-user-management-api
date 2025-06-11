package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/natchaphonbw/usermanagement/modules/users/entities"
	"gorm.io/gorm"
)

type userPostgresRepository struct {
	db *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) UserRepository {
	return &userPostgresRepository{
		db: db,
	}
}

// Create a new user
func (r *userPostgresRepository) CreateUser(ctx context.Context, user *entities.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return result.Error
	}
	return nil
}

// Get all users
func (r *userPostgresRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {

	var users []entities.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		log.Printf("Error getting users: %v", result.Error)
		return nil, result.Error
	}
	return users, nil

}

// Get user by ID
func (r *userPostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil

}

// Update user by ID
func (r *userPostgresRepository) UpdateUserByID(ctx context.Context, id uuid.UUID, data *entities.User) (*entities.User, error) {
	var user entities.User
	// Find user by ID
	findResult := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	// Check if there are any changes
	if user.Name == data.Name && user.Email == data.Email && user.Age == data.Age {
		return &user, nil
	}

	// update user
	UpdateResult := r.db.WithContext(ctx).Model(&user).Updates(data)
	if UpdateResult.Error != nil {
		log.Printf("Error updating user: %v", UpdateResult.Error)
		return nil, UpdateResult.Error
	}

	return &user, nil

}

// Delete user by ID
func (r *userPostgresRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	findResult := r.db.WithContext(ctx).First(&user, "id = ?", id)

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if findResult.Error != nil {
		return nil, findResult.Error
	}

	// delete
	UpdateResult := r.db.WithContext(ctx).Delete(&user)
	if UpdateResult.Error != nil {
		log.Printf("Error deleting user: %v", UpdateResult.Error)
		return nil, UpdateResult.Error
	}

	return &user, nil

}
