package user

import (
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
)

type UserService interface {
	GetAllUsers() ([]User, error)
	GetUserByID(userID uuid.UUID) (*User, error)
	CreateUser(userParams UserRequest) (*User, error)
	UpdateUser(userID uuid.UUID, userParams *UserUpdateRequest) error
	DeleteUser(userID uuid.UUID) error
	UpdateUserRole(userID uuid.UUID, role UserRole) error
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

func (s *userService) GetUserByID(userID uuid.UUID) (*User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(userParams UserRequest) (*User, error) {
	var role UserRole
	if userParams.Role == "" {
		role = RoleUser
	} else {
		role = UserRole(userParams.Role)
	}

	newUser := &User{
		Email:     userParams.Email,
		FirstName: userParams.FirstName,
		LastName:  userParams.LastName,
		Role:      role,
	}
	if err := newUser.HashPassword(userParams.Password); err != nil {
		return nil, err
	}
	createdUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *userService) UpdateUser(userID uuid.UUID, userParams *UserUpdateRequest) error {
	user := &User{
		FirstName: userParams.FirstName,
		LastName:  userParams.LastName,
	}

	err := s.repo.UpdateUser(userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUser(userID uuid.UUID) error {
	if err := s.repo.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUserRole(userID uuid.UUID, role UserRole) error {
	user := &User{
		Role: role,
	}
	return s.repo.UpdateUser(userID, user)
}
