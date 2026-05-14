package middlewares

import (
	"strings"
	"fmt"

	"library-management-system/backend/auth"

	"github.com/gofiber/fiber/v2"

	gojwt "gopkg.in/go-jose/go-jose.v2/jwt"
)

// RequireAuth validates the access JWT (cookie or Authorization: Bearer) on each request.
// Handlers that need the raw access JWT (e.g. after silent refresh) must use auth.AccessTokenFromCtx
// because httpOnly cookies on the wire are unchanged until the next browser request.
func RequireAuth(c *fiber.Ctx) error {
	cfg := auth.Conf
	if cfg == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "auth not configured"})
	}

	token := strings.TrimSpace(c.Cookies(cfg.AccessCookieName))
	if token == "" {
		h := strings.TrimSpace(c.Get("Authorization"))
		if len(h) > 7 && strings.EqualFold(h[:7], "bearer ") {
			token = strings.TrimSpace(h[7:])
		}
	}
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authentication"})
	}

	validated, err := auth.ValidateAccessToken(c.UserContext(), token)
	
	if err != nil {
		rt := strings.TrimSpace(c.Cookies(cfg.RefreshCookieName))
		if rt == "" {
			fmt.Println("error 2: ", err, gojwt.ErrExpired)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		tr, rerr := auth.RefreshTokens(c.UserContext(), cfg, rt)
		fmt.Println("@@@@@", tr)
		if rerr != nil {
			fmt.Println("error 3: ", rerr, gojwt.ErrExpired)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		validated, err = auth.ValidateAccessToken(c.UserContext(), tr.AccessToken)
		if err != nil {
			fmt.Println("error 4: ", err, gojwt.ErrExpired)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		c.Locals("claims", validated)
		c.Locals("sub", validated.RegisteredClaims.Subject)
		c.Locals(auth.AccessTokenCtxKey, strings.TrimSpace(tr.AccessToken))
		
		auth.SetAuthCookies(c, cfg, tr)
		return c.Next()
	}

	c.Locals("claims", validated)
	c.Locals("sub", validated.RegisteredClaims.Subject)
	c.Locals(auth.AccessTokenCtxKey, token)
	return c.Next()
}
