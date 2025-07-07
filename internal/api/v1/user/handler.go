package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
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

	return c.JSON(fiber.Map{"users": users})
}
