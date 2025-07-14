package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, authHandler *AuthHandler) {
	authRouter := router.Group("/auth")
	{
		// Public routes
		authRouter.Post("/signin", authHandler.Login)
		authRouter.Post("/refresh", authHandler.RefreshToken)
		authRouter.Post("/signout", authHandler.Logout)
		authRouter.Post("/signup", authHandler.Register)
	}
}
