package middlewares

import (
	"strings"

	"library-management-system/backend/auth"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth validates access JWT (primary or passwordless app) from Bearer or httpOnly cookie.
func RequireAuth(c *fiber.Ctx) error {
	token := auth.AccessTokenFromRequest(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authentication"})
	}

	_, cfg, err := auth.ValidateAccessTokenAny(c.UserContext(), token)
	if err != nil {
		rt := auth.RefreshTokenFromRequest(c, nil)
		if rt == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		var tr *auth.TokenResponse
		var refreshCfg *auth.Config
		for _, app := range auth.AuthConfigs() {
			if strings.TrimSpace(c.Cookies(app.RefreshCookieName)) == "" {
				continue
			}
			tr, err = auth.RefreshTokens(c.UserContext(), app, rt)
			if err == nil {
				refreshCfg = app
				break
			}
		}
		if tr == nil || refreshCfg == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		if _, err = auth.ValidateAccessTokenForConfig(c.UserContext(), refreshCfg, tr.AccessToken); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		cfg = refreshCfg
		auth.SetAuthCookies(c, cfg, tr)
		token = strings.TrimSpace(tr.AccessToken)
	}

	c.Locals(auth.AccessTokenCtxKey, token)
	return c.Next()
}
