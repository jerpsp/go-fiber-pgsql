package book

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func RegisterRoutes(cfg *config.Config, router fiber.Router, handler *BookHandler) {
	bookGroup := router.Group("/books")
	{
		// Public routes - anyone can access
		bookGroup.Get("", handler.GetBooks)
		bookGroup.Get("/:id", handler.GetBook)

		// User authenticated routes - any authenticated user can access
		bookGroup.Post("", middleware.JWTMiddleware(cfg), handler.CreateBook)

		// Moderator or Admin routes - only moderators and admins can update books
		bookGroup.Patch("/:id", middleware.JWTMiddleware(cfg), middleware.ModeratorOrAdmin(), handler.UpdateBook)

		// Admin only routes - only admins can delete books
		bookGroup.Delete("/:id", middleware.JWTMiddleware(cfg), middleware.AdminOnly(), handler.DeleteBook)
	}
}
