package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AccessTokenFromRequest reads Bearer token or access-token cookies (primary or passwordless app).
func AccessTokenFromRequest(c *fiber.Ctx) string {
	if h := strings.TrimSpace(c.Get("Authorization")); len(h) > 7 && strings.EqualFold(h[:7], "bearer ") {
		if t := strings.TrimSpace(h[7:]); t != "" {
			return t
		}
	}
	for _, cfg := range AuthConfigs() {
		if t := strings.TrimSpace(c.Cookies(cfg.AccessCookieName)); t != "" {
			return t
		}
	}
	return ""
}

// RefreshTokenFromRequest reads the refresh cookie for cfg, or scans all Auth0 app configs when cfg is nil.
func RefreshTokenFromRequest(c *fiber.Ctx, cfg *Config) string {
	if cfg != nil {
		return strings.TrimSpace(c.Cookies(cfg.RefreshCookieName))
	}
	for _, app := range AuthConfigs() {
		if t := strings.TrimSpace(c.Cookies(app.RefreshCookieName)); t != "" {
			return t
		}
	}
	return ""
}
