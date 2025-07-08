package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *UserHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.Get("", handler.GetAllUsers)
		userGroup.Get("/:id", handler.GetUserByID)
		userGroup.Post("", handler.CreateUser)
		userGroup.Patch("/:id", handler.UpdateUser)
		userGroup.Delete("/:id", handler.DeleteUser)
	}
}
