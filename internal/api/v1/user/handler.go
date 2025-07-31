package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type UserHandler struct {
	config  *config.Config
	service UserService
}

func NewUserHandler(config *config.Config, service UserService) *UserHandler {
	return &UserHandler{config: config, service: service}
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	pagination := PaginationRequest{
		Page:  1,
		Limit: 10,
	}

	if c.Query("page") != "" {
		if err := c.QueryParser(&pagination); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pagination parameters"})
		}
	}

	users, total, err := h.service.GetAllUsers(c, pagination.Page, pagination.Limit)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := int(total) / pagination.Limit
	if int(total)%pagination.Limit > 0 {
		totalPages++
	}

	response := PaginatedResponse{
		Users:      users,
		Total:      total,
		Page:       pagination.Page,
		PerPage:    pagination.Limit,
		TotalPages: totalPages,
	}

	return c.JSON(response)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	user, err := h.service.GetUserByID(c, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	file, _ := c.FormFile("profile_image")
	var user UserCreateRequest
	// fmt.Println("user create request body:", string(c.Body()))
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	createdUser, err := h.service.CreateUser(c, user, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": createdUser})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var user UserUpdateRequest

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	if err := h.service.UpdateUser(c, id, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.DeleteUser(c, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	var roleUpdate UserRoleUpdateRequest

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := c.BodyParser(&roleUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate(&roleUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	if err := h.service.UpdateUserRole(c, id, UserRole(roleUpdate.Role)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User role updated successfully"})
}
