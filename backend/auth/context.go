package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AccessTokenCtxKey is the fiber Locals key for the validated access JWT string
// (same value the client sent via cookie or Authorization: Bearer). Set by RequireAuth
// after validation or silent refresh so handlers do not read stale cookies mid-request.
const AccessTokenCtxKey = "auth_access_token"

// AccessTokenFromCtx returns the access JWT placed by RequireAuth, or "" if missing.
func AccessTokenFromCtx(c *fiber.Ctx) string {
	v := c.Locals(AccessTokenCtxKey)
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}
