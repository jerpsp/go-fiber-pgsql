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
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"users": users, "count": len(users)})
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"user": user})
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user UserRequest
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	createdUser, err := h.service.CreateUser(user)
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

	if err := utils.GetValidator().Struct(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	if err := h.service.UpdateUser(id, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
