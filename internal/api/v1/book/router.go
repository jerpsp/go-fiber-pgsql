package book

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *BookHandler) {
	bookGroup := router.Group("/books")
	{
		bookGroup.Get("", handler.GetBooks)
		bookGroup.Get("/:id", handler.GetBook)
		bookGroup.Post("", handler.CreateBook)
		bookGroup.Patch("/:id", handler.UpdateBook)
		bookGroup.Delete("/:id", handler.DeleteBook)
		bookGroup.Post("/upload-file", handler.UploadFile)
	}
}
