package book

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, handler *BookHandler) {
	bookGroup := router.Group("/books")
	{
		// Public routes
		bookGroup.Get("", handler.GetBooks)
		bookGroup.Get("/:id", handler.GetBook)

		// Protected routes
		bookGroup.Post("", middleware.JWTMiddleware(cfg), handler.CreateBook)
		bookGroup.Patch("/:id", middleware.JWTMiddleware(cfg), handler.UpdateBook)
		bookGroup.Delete("/:id", middleware.JWTMiddleware(cfg), handler.DeleteBook)
	}
}
