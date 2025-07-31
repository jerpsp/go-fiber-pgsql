package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

type UserHandlerSuite struct {
	suite.Suite
	mockUserSvc *mocks.UserService
	userHandler user.UserHandler
	router      *fiber.App
	httptest    *httptest.ResponseRecorder
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

func (s *UserHandlerSuite) SetupTest() {
	// Setup Test Engine
	s.mockUserSvc = mocks.NewUserService(s.T())
	s.userHandler = *user.NewUserHandler(&config.Config{}, s.mockUserSvc)
	s.router = fiber.New()
	s.httptest = httptest.NewRecorder()

	// Setup Router
	s.router = fiber.New()

	// User
	userGroup := s.router.Group("api/v1/users")
	{
		userGroup.Get("", s.userHandler.GetAllUsers)
		userGroup.Get("/:id", s.userHandler.GetUserByID)
		userGroup.Post("", s.userHandler.CreateUser)
		userGroup.Patch("/:id", s.userHandler.UpdateUser)
		userGroup.Delete("/:id", s.userHandler.DeleteUser)
		userGroup.Patch("/:id/role", s.userHandler.UpdateUserRole)
	}

}

func (s *UserHandlerSuite) TestGetAllUsers1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get All User Success", 34)

	// Setup Mock
	searchId := uuid.New()
	handlerResponse := user.PaginatedResponse{Users: []user.User{{ID: searchId, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}}, Total: 1, Page: 1, PerPage: 10, TotalPages: 1}
	serviceResponse := []user.User{{ID: searchId, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}}
	s.mockUserSvc.EXPECT().GetAllUsers(mock.Anything, mock.Anything, mock.Anything).Return(serviceResponse, 1, nil)

	// Setup Request
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(handlerResponse)
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestGetUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get User Success", 34)

	// Setup Mock
	userID := uuid.New()
	handlerResponse := fiber.Map{"user": user.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}}
	serviceResponse := user.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	s.mockUserSvc.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(&serviceResponse, nil)

	// Setup Request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", serviceResponse.ID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(handlerResponse)
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestGetUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get User Parse UUID Fail", 31)

	// Setup Mock
	userID := "invalid-uuid"

	// Setup Request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", userID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "Invalid ID format"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestGetUser3() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Get User Not Found", 31)

	// Setup Mock
	userID := uuid.New()
	s.mockUserSvc.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(nil, fmt.Errorf("user not found"))

	// Setup Request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", userID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "user not found"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusNotFound, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestCreateUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Create User Success", 34)

	// Setup Mock
	userID := uuid.New()
	handlerResponse := user.User{ID: userID, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Role: user.RoleUser}
	s.mockUserSvc.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything).Return(&handlerResponse, nil)

	// Setup Request
	reqBody := user.UserCreateRequest{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		Role:      "user",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"user": handlerResponse})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusCreated, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestCreateUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Create User Validation Fail", 31)

	// Setup Request
	reqBody := user.UserCreateRequest{
		Email:     "invalid-email",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{
		"details": "validation error: field 'Email' failed on the 'email' tag",
		"error":   "Validation failed",
	})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestCreateUser3() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Create User Internal Server Error", 31)

	// Setup Mock
	s.mockUserSvc.EXPECT().CreateUser(mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("internal server error"))

	// Setup Request
	reqBody := user.UserCreateRequest{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		Role:      "user",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "internal server error"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusInternalServerError, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Success", 34)

	// Setup Mock
	userID := uuid.New()
	handlerResponse := fiber.Map{"message": "User updated successfully"}
	s.mockUserSvc.EXPECT().UpdateUser(mock.Anything, userID, mock.Anything).Return(nil)

	// Setup Request
	reqBody := user.UserUpdateRequest{
		FirstName: "John",
		LastName:  "Doe",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(handlerResponse)
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Parse UUID Fail", 31)

	// Setup Request
	userID := "invalid-uuid"
	reqBody := user.UserUpdateRequest{
		FirstName: "John",
		LastName:  "Doe",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "Invalid ID format"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUser3() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Validation Fail", 31)

	// Setup Request
	userID := uuid.New()
	reqBody := user.UserUpdateRequest{
		FirstName: "",
		LastName:  "Doe", // Last name is required
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{
		"details": "validation error: field 'FirstName' failed on the 'required' tag",
		"error":   "Validation failed",
	})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUser4() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Internal Server Error", 31)

	// Setup Mock
	userID := uuid.New()
	s.mockUserSvc.EXPECT().UpdateUser(mock.Anything, userID, mock.Anything).Return(fmt.Errorf("internal server error"))

	// Setup Request
	reqBody := user.UserUpdateRequest{
		FirstName: "John",
		LastName:  "Doe",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "internal server error"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusInternalServerError, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestDeleteUser1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Delete User Success", 34)

	// Setup Mock
	userID := uuid.New()
	handlerResponse := fiber.Map{"message": "User deleted successfully"}
	s.mockUserSvc.EXPECT().DeleteUser(mock.Anything, userID).Return(nil)

	// Setup Request
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", userID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(handlerResponse)
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestDeleteUser2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Delete User Parse UUID Fail", 31)

	// Setup Request
	userID := "invalid-uuid"
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", userID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "Invalid ID format"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestDeleteUser3() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Delete User Internal Server Error", 31)

	// Setup Mock
	userID := uuid.New()
	s.mockUserSvc.EXPECT().DeleteUser(mock.Anything, userID).Return(fmt.Errorf("internal Server Error"))

	// Setup Request
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", userID), nil)

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "internal Server Error"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusInternalServerError, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUserRole1() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Role Success", 34)

	// Setup Mock
	userID := uuid.New()
	roleUpdate := user.UserRoleUpdateRequest{Role: "admin"}
	handlerResponse := fiber.Map{"message": "User role updated successfully"}
	s.mockUserSvc.EXPECT().UpdateUserRole(mock.Anything, userID, user.RoleAdmin).Return(nil)

	// Setup Request
	reqBodyBytes, _ := json.Marshal(roleUpdate)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s/role", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(handlerResponse)
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUserRole2() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Role Parse UUID Fail", 31)

	// Setup Request
	userID := "invalid-uuid"
	roleUpdate := user.UserRoleUpdateRequest{Role: "admin"}
	reqBodyBytes, _ := json.Marshal(roleUpdate)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s/role", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "Invalid ID format"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUserRole3() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Role Validation Fail", 31)

	// Setup Request
	userID := uuid.New()
	roleUpdate := user.UserRoleUpdateRequest{Role: ""}
	reqBodyBytes, _ := json.Marshal(roleUpdate)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s/role", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{
		"details": "validation error: field 'Role' failed on the 'required' tag",
		"error":   "Validation failed",
	})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}

func (s *UserHandlerSuite) TestUpdateUserRole4() {
	// Case Name Print In Test
	utils.ConsolePrintColoredText("CASE: Update User Role Internal Server Error", 31)

	// Setup Mock
	userID := uuid.New()
	roleUpdate := user.UserRoleUpdateRequest{Role: "admin"}
	s.mockUserSvc.EXPECT().UpdateUserRole(mock.Anything, userID, user.RoleAdmin).Return(fmt.Errorf("internal server error"))

	// Setup Request
	reqBodyBytes, _ := json.Marshal(roleUpdate)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s/role", userID), bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Run Test Request
	resp, _ := s.router.Test(req)

	expectedResp, _ := json.Marshal(fiber.Map{"error": "internal server error"})
	actualResp, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusInternalServerError, resp.StatusCode)
	s.Equal(string(expectedResp), string(actualResp))
}
