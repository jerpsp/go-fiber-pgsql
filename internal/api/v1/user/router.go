package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, handler *UserHandler) {
	userRouter := router.Group("/users")
	{
		userRouter.Get("", middleware.JWTMiddleware(cfg), middleware.AdminOnly(), handler.GetAllUsers)
		userRouter.Get("/:id", middleware.JWTMiddleware(cfg), handler.GetUserByID)
		userRouter.Post("", middleware.JWTMiddleware(cfg), middleware.AdminOnly(), handler.CreateUser)
		userRouter.Patch("/:id", middleware.JWTMiddleware(cfg), handler.UpdateUser)
		userRouter.Patch("/:id/role", middleware.JWTMiddleware(cfg), middleware.AdminOnly(), handler.UpdateUserRole)
		userRouter.Delete("/:id", middleware.JWTMiddleware(cfg), middleware.AdminOnly(), handler.DeleteUser)
	}
}
