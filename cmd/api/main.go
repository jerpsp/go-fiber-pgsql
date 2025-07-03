package main

import (
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/book"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
)

func main() {
	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize database connection
	db := database.NewGormDB(cfg.PostgresDB)

	bookRepo := book.NewBookRepository(cfg, db)
	bookService := book.NewBookService(cfg, bookRepo)
	bookHandler := book.NewBookHandler(cfg, bookService)

	// Start the server with the book handler
	internal.StartServer(cfg, bookHandler)
}
