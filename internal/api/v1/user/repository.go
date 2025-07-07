package user

import (
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

type UserRepository interface {
	FindAllUsers() ([]User, error)
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
