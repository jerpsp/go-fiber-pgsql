package user

import (
	"fmt"
	"mime/multipart"
	"path"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/storage"
)

type UserService interface {
	GetAllUsers(c *fiber.Ctx, page, limit int) ([]User, int64, error)
	GetUserByID(c *fiber.Ctx, userID uuid.UUID) (*User, error)
	CreateUser(c *fiber.Ctx, userParams UserCreateRequest, file *multipart.FileHeader) (*User, error)
	UpdateUser(c *fiber.Ctx, userID uuid.UUID, userParams *UserUpdateRequest) error
	DeleteUser(c *fiber.Ctx, userID uuid.UUID) error
	UpdateUserRole(c *fiber.Ctx, userID uuid.UUID, role UserRole) error
}

type userService struct {
	config *config.Config
	repo   UserRepository
	s3Repo storage.S3Repository
}

func NewUserService(config *config.Config, repo UserRepository, s3Repo storage.S3Repository) UserService {
	return &userService{config: config, repo: repo, s3Repo: s3Repo}
}

func (s *userService) GetAllUsers(c *fiber.Ctx, page, limit int) ([]User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := s.repo.FindAllUsers(c, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *userService) GetUserByID(c *fiber.Ctx, userID uuid.UUID) (*User, error) {
	user, err := s.repo.FindUserByID(c, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(c *fiber.Ctx, userParams UserCreateRequest, file *multipart.FileHeader) (*User, error) {
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

	if file != nil {
		// validate file type
		fileExt := path.Ext(file.Filename)
		fileType := []string{".png", ".jpg", ".jpeg"}
		if !slices.Contains(fileType, fileExt) {
			return nil, fmt.Errorf("invalid file format: %s", fileExt)
		}
		//validate file size
		if file.Size > 200*1024 {
			return nil, fmt.Errorf("file size exceeds limit: %d", file.Size)
		}

		filePath, err := s.s3Repo.UploadPublicFile(file)
		if err != nil {
			return nil, err
		}
		newUser.ProfileImage = filePath
	}

	createdUser, err := s.repo.CreateUser(c, newUser)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *userService) UpdateUser(c *fiber.Ctx, userID uuid.UUID, userParams *UserUpdateRequest) error {
	user := &User{
		FirstName: userParams.FirstName,
		LastName:  userParams.LastName,
	}

	err := s.repo.UpdateUser(c, userID, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeleteUser(c *fiber.Ctx, userID uuid.UUID) error {
	if err := s.repo.DeleteUser(c, userID); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateUserRole(c *fiber.Ctx, userID uuid.UUID, role UserRole) error {
	user := &User{
		Role: role,
	}
	return s.repo.UpdateUser(c, userID, user)
}
