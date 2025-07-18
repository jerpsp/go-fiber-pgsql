package user

import (
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

type UserRepository interface {
	FindAllUsers(page, perPage int) ([]User, int64, error)
	FindUserByID(userID uuid.UUID) (*User, error)
	FindUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(userID uuid.UUID, user *User) error
	DeleteUser(userID uuid.UUID) error
}

type userRepository struct {
	config *config.Config
	db     *database.GormDB
}

func NewUserRepository(cfg *config.Config, db *database.GormDB) UserRepository {
	return &userRepository{config: cfg, db: db}
}

func (r *userRepository) FindAllUsers(page, perPage int) ([]User, int64, error) {
	var users []User
	var total int64

	// Count total records
	if err := r.db.DB.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * perPage
	if err := r.db.DB.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindUserByID(userID uuid.UUID) (*User, error) {
	var user User
	if err := r.db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *User) (*User, error) {
	if err := r.db.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(userID uuid.UUID, user *User) error {
	if err := r.db.DB.Where("id = ?", userID).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(userID uuid.UUID) error {
	if err := r.db.DB.Delete(&User{}, "id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
