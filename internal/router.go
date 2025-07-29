package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/auth"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/book"
	"github.com/jerpsp/go-fiber-beginner/internal/api/v1/user"
	"github.com/jerpsp/go-fiber-beginner/middleware"
)

func StartServer(cfg *config.Config, bookHandler *book.BookHandler,
	userHandler *user.UserHandler, authHandler *auth.AuthHandler) {

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		ReadTimeout:   cfg.Server.Timeout * time.Second,
		WriteTimeout:  cfg.Server.Timeout * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	// app.Use(logger.New())
	app.Use(helmet.New())
	app.Use(middleware.Cors(cfg))
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	apiV1 := app.Group("/api/v1")

	auth.RegisterRoutes(cfg, apiV1, authHandler)
	book.RegisterRoutes(cfg, apiV1, bookHandler)
	user.RegisterRoutes(cfg, apiV1, userHandler)

	// Use PORT from environment if available, otherwise use config
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(cfg.Server.Port)
	}

	app.Listen(fmt.Sprintf(":%s", port))
}
