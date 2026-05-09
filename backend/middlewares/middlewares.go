package middlewares

import (
	"os"
	"strings"
	"time"

	"library-management-system/backend/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

const handlerTimeout = 30 * time.Second

// Stack returns middleware for the HTTP API. CORS policy depends on environment.
func Stack(environment string) []fiber.Handler {
	var corsMW fiber.Handler
	switch environment {
	case initializers.EnvDevelopment:
		corsMW = cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
			AllowCredentials: false,
		})
	default:
		origins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
		corsMW = cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
			AllowCredentials: origins != "" && origins != "*",
		})
	}

	return []fiber.Handler{
		recover.New(),
		requestid.New(requestid.Config{
			Header:     fiber.HeaderXRequestID,
			ContextKey: "requestid",
		}),
		corsMW,
		logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path} | ${error}\n",
		}),
		timeout.NewWithContext(func(c *fiber.Ctx) error {
			return c.Next()
		}, handlerTimeout),
	}
}
