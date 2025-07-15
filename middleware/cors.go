package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jerpsp/go-fiber-beginner/config"
)

// func UseCorsMiddleware(cfg *config.Config, r *fiber.App) {
// 	// Middleware setup can be done here
// 	// For example, you can set up logging, error handling, etc.

// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins: cfg.Server.AllowOrigins,
// 	}))
// }

func Cors(cfg *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Authorization",
	})
}
