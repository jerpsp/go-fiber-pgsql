package user

import "github.com/jerpsp/go-fiber-beginner/config"

type UserService interface {
	GetAllUsers() ([]User, error)
}

type userService struct {
	config *config.Config
	repo   UserRepository
}

func NewUserService(config *config.Config, repo UserRepository) UserService {
	return &userService{config: config, repo: repo}
}

func (s *userService) GetAllUsers() ([]User, error) {
	users, err := s.repo.FindAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
