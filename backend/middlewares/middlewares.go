package middlewares

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

// Stack returns middleware recommended for a production-grade HTTP API:
//   - Recover: turns panics into 500s instead of crashing the process
//   - RequestID: stable id per request for logs and tracing (custom header + generator)
//   - Logger: structured access logging
//   - Timeout: bounds handler work so slow endpoints cannot hang workers indefinitely
//   - CORS: configurable cross-origin rules when a browser client exists
func Stack() []fiber.Handler {
	reqTimeout := 30 * time.Second
	if v := os.Getenv("HTTP_HANDLER_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			reqTimeout = time.Duration(n) * time.Second
		}
	}

	origins := []string{}
	if raw := os.Getenv("CORS_ALLOW_ORIGINS"); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if t := strings.TrimSpace(o); t != "" {
				origins = append(origins, t)
			}
		}
	}
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	corsMW := cors.New(cors.Config{
		AllowOrigins:     strings.Join(origins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowCredentials: len(origins) > 0 && origins[0] != "*",
	})

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
		}, reqTimeout),
	}
}
