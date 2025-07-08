package user

import (
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

type UserRepository interface {
	FindAllUsers() ([]User, error)
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

func (r *userRepository) FindAllUsers() ([]User, error) {
	var users []User
	if err := r.db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
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
