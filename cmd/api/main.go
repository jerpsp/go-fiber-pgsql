package main

import (
	"context"
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/auth"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/book"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/jerpsp/go-fiber-beginner/pkg/storage"
)

func main() {
	// Initialize configuration
	cfg := config.InitConfig()

	// Initialize database connection
	db := database.NewGormDB(cfg.PostgresDB)
	redis := database.NewRedisClient(cfg.Redis)

	s3Client := storage.NewS3Client(cfg.AWS)
	s3Repo := storage.NewS3Repo(s3Client)

	fmt.Println(redis.Client.Ping(context.Background()))

	bookRepo := book.NewBookRepository(cfg, db)
	bookService := book.NewBookService(cfg, bookRepo)
	bookHandler := book.NewBookHandler(cfg, bookService)

	userRepo := user.NewUserRepository(cfg, db)
	userService := user.NewUserService(cfg, userRepo, s3Repo)
	userHandler := user.NewUserHandler(cfg, userService)

	tokenRepo := auth.NewAuthRepository(cfg, db, redis)
	authService := auth.NewAuthService(cfg, userRepo, tokenRepo)
	authHandler := auth.NewAuthHandler(cfg, authService)

	// Start the server with handlers and db
	internal.StartServer(cfg, bookHandler, userHandler, authHandler)
}
