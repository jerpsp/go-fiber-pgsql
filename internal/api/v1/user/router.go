package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *UserHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.Get("/", handler.GetAllUsers)
	}
}
