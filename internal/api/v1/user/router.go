package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, handler *UserHandler) {
	userRouter := router.Group("/users")
	{
		// Public routes
		userRouter.Post("", handler.CreateUser)

		// Protected routes
		userRouter.Get("", middleware.JWTMiddleware(cfg), handler.GetAllUsers)
		userRouter.Get("/:id", middleware.JWTMiddleware(cfg), handler.GetUserByID)
		userRouter.Patch("/:id", middleware.JWTMiddleware(cfg), handler.UpdateUser)
		userRouter.Delete("/:id", middleware.JWTMiddleware(cfg), handler.DeleteUser)
	}
}
