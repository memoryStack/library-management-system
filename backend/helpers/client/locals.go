package client

import "github.com/gofiber/fiber/v2"

// Locals keys set by DetectAuthClient middleware. Handlers read with:
//
//	if c.Locals("clientKnown") != true { ... }
//	if c.Locals("isWebClient") == true { ... }   // browser → httpOnly cookies
//	if c.Locals("isNativeClient") == true { ... } // native → tokens in JSON
const (
	localsClientKnown    = "clientKnown"
	localsIsWebClient    = "isWebClient"
	localsIsNativeClient = "isNativeClient"
)

// SetLocals stores boolean flags on the request context (middleware only).
func SetLocals(c *fiber.Ctx, info Info) {
	c.Locals(localsClientKnown, info.OK)
	c.Locals(localsIsWebClient, info.OK && info.Kind == Web)
	c.Locals(localsIsNativeClient, info.OK && info.Kind == Native)
}
