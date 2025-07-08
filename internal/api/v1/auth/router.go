package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, authHandler *AuthHandler) {
	authRouter := router.Group("/auth")
	{
		// Public routes
		authRouter.Post("/login", authHandler.Login)
		authRouter.Post("/refresh", authHandler.RefreshToken)

		// Protected routes
		authRouter.Post("/logout", middleware.JWTMiddleware(cfg), authHandler.Logout)
	}
}
