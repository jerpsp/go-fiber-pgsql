package user_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/mocks"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
	mockRepo    *mocks.UserRepository
	mockStorage *mocks.S3Repository
	service     user.UserService
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
func (s *UserServiceSuite) SetupTest() {
	s.mockRepo = mocks.NewUserRepository(s.T())
	s.mockStorage = mocks.NewS3Repository(s.T())
	s.service = user.NewUserService(&config.Config{}, s.mockRepo, s.mockStorage)
}

func (s *UserServiceSuite) TestGetAllUsers1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get All User Success", 34)

	// Setup Mock
	repoResponse := []user.User{
		{ID: uuid.New(), FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Role: user.RoleUser},
		{ID: uuid.New(), FirstName: "Jane", LastName: "Doe", Email: "jane.doe@example.com", Role: user.RoleUser},
	}
	s.mockRepo.EXPECT().FindAllUsers(mock.Anything, 1, 10).Return(repoResponse, int64(len(repoResponse)), nil)

	// Call the service method
	users, total, err := s.service.GetAllUsers(&fiber.Ctx{}, 1, 10)

	// Assertions
	s.NoError(err)
	s.Equal(2, len(users))
	s.Equal(int64(2), total)
	s.Equal(repoResponse, users)
}

func (s *UserServiceSuite) TestGetAllUsers2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get All User Fail Repo", 34)

	// Setup Mock
	s.mockRepo.EXPECT().FindAllUsers(mock.Anything, 1, 10).Return(nil, 0, fmt.Errorf("repository error"))

	// Call the service method
	users, total, err := s.service.GetAllUsers(&fiber.Ctx{}, 1, 10)

	// Assertions
	s.Error(err)
	s.Nil(users)
	s.Equal(int64(0), total)
}

func (s *UserServiceSuite) TestGetUserByID1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get User By ID Success", 34)

	// Setup Mock
	userID := uuid.New()
	repoResponse := &user.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Role: user.RoleUser}
	s.mockRepo.EXPECT().FindUserByID(mock.Anything, userID).Return(repoResponse, nil)

	// Call the service method
	user, err := s.service.GetUserByID(&fiber.Ctx{}, userID)

	// Assertions
	s.NoError(err)
	s.Equal(repoResponse, user)
}

func (s *UserServiceSuite) TestGetUserByID2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get User By ID Fail Repo", 34)

	// Setup Mock
	userID := uuid.New()
	s.mockRepo.EXPECT().FindUserByID(mock.Anything, userID).Return(nil, fmt.Errorf("repository error"))

	// Call the service method
	user, err := s.service.GetUserByID(&fiber.Ctx{}, userID)

	// Assertions
	s.Error(err)
	s.Nil(user)
}

func (s *UserServiceSuite) TestCreateUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Create User Success", 34)

	// Setup Mock
	userParams := user.UserCreateRequest{
		Email: "john.doe@example.com", FirstName: "John", LastName: "Doe", Password: "password", Role: "user",
	}
	userResponse := &user.User{
		ID: uuid.New(), Email: userParams.Email, FirstName: userParams.FirstName, LastName: userParams.LastName, Role: user.RoleUser,
	}

	s.mockRepo.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("*user.User")).Return(userResponse, nil)

	// Call the service method
	user, err := s.service.CreateUser(&fiber.Ctx{}, userParams, nil)

	// Assertions
	s.NoError(err)
	s.Equal(userResponse, user)
}

func (s *UserServiceSuite) TestCreateUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Create User Fail Repo", 34)

	// Setup Mock
	userParams := user.UserCreateRequest{
		Email: "john.doe@example.com", FirstName: "John", LastName: "Doe", Password: "password", Role: "user",
	}
	s.mockRepo.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("*user.User")).Return(nil, fmt.Errorf("repository error"))

	// Call the service method
	user, err := s.service.CreateUser(&fiber.Ctx{}, userParams, nil)

	// Assertions
	s.Error(err)
	s.Nil(user)
}

func (s *UserServiceSuite) TestUpdateUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Success", 34)

	// Setup Mock
	userID := uuid.New()
	userParams := &user.UserUpdateRequest{FirstName: "John", LastName: "Doe"}
	userResponse := &user.User{
		FirstName: userParams.FirstName, LastName: userParams.LastName,
	}
	s.mockRepo.EXPECT().UpdateUser(mock.Anything, userID, userResponse).Return(nil)

	// Call the service method
	err := s.service.UpdateUser(&fiber.Ctx{}, userID, userParams)

	// Assertions
	s.NoError(err)
}

func (s *UserServiceSuite) TestUpdateUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Fail Repo", 34)

	// Setup Mock
	userID := uuid.New()
	userParams := &user.UserUpdateRequest{FirstName: "John", LastName: "Doe"}
	userResponse := &user.User{
		FirstName: userParams.FirstName, LastName: userParams.LastName,
	}
	s.mockRepo.EXPECT().UpdateUser(mock.Anything, userID, userResponse).Return(fmt.Errorf("repository error"))

	// Call the service method
	err := s.service.UpdateUser(&fiber.Ctx{}, userID, userParams)

	// Assertions
	s.Error(err)
}

func (s *UserServiceSuite) TestDeleteUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Delete User Success", 34)

	// Setup Mock
	userID := uuid.New()
	s.mockRepo.EXPECT().DeleteUser(mock.Anything, userID).Return(nil)

	// Call the service method
	err := s.service.DeleteUser(&fiber.Ctx{}, userID)

	// Assertions
	s.NoError(err)
}

func (s *UserServiceSuite) TestDeleteUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Delete User Fail Repo", 34)

	// Setup Mock
	userID := uuid.New()
	s.mockRepo.EXPECT().DeleteUser(mock.Anything, userID).Return(fmt.Errorf("repository error"))

	// Call the service method
	err := s.service.DeleteUser(&fiber.Ctx{}, userID)

	// Assertions
	s.Error(err)
}

func (s *UserServiceSuite) TestUpdateUserRole1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Role Success", 34)

	// Setup Mock
	userID := uuid.New()
	role := user.RoleUser
	userResponse := &user.User{Role: role}
	s.mockRepo.EXPECT().UpdateUser(mock.Anything, userID, userResponse).Return(nil)

	// Call the service method
	err := s.service.UpdateUserRole(&fiber.Ctx{}, userID, role)

	// Assertions
	s.NoError(err)
}

// TODO: TEST Upload File
