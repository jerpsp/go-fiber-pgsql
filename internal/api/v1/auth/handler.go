package auth

import (
	"github.com/gofiber/fiber/v2"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	tokens, err := h.service.Login(c, req)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid email or password"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to authenticate"})
	}

	return c.JSON(tokens)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	tokens, err := h.service.RefreshToken(c, req.RefreshToken)
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

		return c.Status(status).JSON(fiber.Map{"message": message})
	}

	return c.JSON(tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	err := h.service.LogoutWithToken(c, req.RefreshToken)
	if err != nil {
		status := fiber.StatusInternalServerError
		message := "Failed to logout"

		switch err {
		case utils.ErrInvalidToken:
			status = fiber.StatusUnauthorized
			message = "Invalid refresh token"
		case utils.ErrTokenExpired:
			status = fiber.StatusUnauthorized
			message = "Refresh token expired"
		}

		return c.Status(status).JSON(fiber.Map{"message": message})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	err := h.service.Register(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}
