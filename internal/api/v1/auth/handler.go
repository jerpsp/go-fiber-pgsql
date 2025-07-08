package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

type AuthHandler struct {
	config  *config.Config
	service AuthService
}

func NewAuthHandler(cfg *config.Config, service AuthService) *AuthHandler {
	return &AuthHandler{config: cfg, service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tokens, err := h.service.Login(req)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to authenticate"})
	}

	return c.JSON(tokens)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tokens, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		status := fiber.StatusInternalServerError
		message := "Failed to refresh token"

		switch err {
		case utils.ErrInvalidToken:
			status = fiber.StatusUnauthorized
			message = "Invalid refresh token"
		case utils.ErrTokenExpired:
			status = fiber.StatusUnauthorized
			message = "Refresh token expired"
		}

		return c.Status(status).JSON(fiber.Map{"error": message})
	}

	return c.JSON(tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err := h.service.Logout(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to logout"})
	}

	return c.SendStatus(fiber.StatusOK)
}
