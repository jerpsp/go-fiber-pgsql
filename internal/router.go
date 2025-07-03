package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/book"
)

func StartServer(cfg *config.Config, bookHandler *book.BookHandler) {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ReadTimeout:   cfg.Server.Timeout * time.Second,
		WriteTimeout:  cfg.Server.Timeout * time.Second,
	})

	SetupMiddleware(cfg, app)

	apiV1 := app.Group("/api/v1")
	{
		book.RegisterRoutes(apiV1, bookHandler)
	}

	app.Listen(fmt.Sprintf(":%s", strconv.Itoa(cfg.Server.Port)))
}
