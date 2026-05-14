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

func devCORSOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("AUTH0_CORS_ORIGINS"))
	if raw == "" {
		return []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// Stack returns middleware for the HTTP API. CORS policy depends on environment.
func Stack(environment string) []fiber.Handler {
	var corsMW fiber.Handler
	corsHeaders := "Origin, Content-Type, Accept, Authorization, X-Request-ID, Cookie"
	switch environment {
	case initializers.EnvDevelopment:
		origins := devCORSOrigins()
		corsMW = cors.New(cors.Config{
			AllowOrigins:     strings.Join(origins, ","),
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     corsHeaders,
			AllowCredentials: true,
		})
	default:
		origins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
		corsMW = cors.New(cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     corsHeaders,
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
