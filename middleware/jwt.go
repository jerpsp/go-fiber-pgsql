package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/pkg/utils"
)

// Protected is a middleware that checks if the user is authenticated
func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate the token
		userInfo, err := utils.ValidateToken(cfg, tokenString, utils.AccessToken)
		if err != nil {
			status := fiber.StatusUnauthorized
			message := "Unauthorized"

			switch err {
			case utils.ErrInvalidToken:
				message = "Invalid token"
			case utils.ErrTokenExpired:
				message = "Token expired"
			}

			return c.Status(status).JSON(fiber.Map{
				"error": message,
			})
		}

		// Set user information in context
		c.Locals("userID", userInfo.ID)
		c.Locals("email", userInfo.Email)
		c.Locals("role", userInfo.Role)

		return c.Next()
	}
}
