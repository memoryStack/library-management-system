package middlewares

import (
	"library-management-system/backend/helpers/client"

	"github.com/gofiber/fiber/v2"
)

// DetectAuthClient resolves web vs native and sets on c.Locals:
//   - client.LocalsClientKnown (bool)
//   - client.LocalsIsWebClient (bool)
//   - client.LocalsIsNativeClient (bool)
func DetectAuthClient() fiber.Handler {
	uaParser, parseErr := client.Parser()
	return func(c *fiber.Ctx) error {
		var info client.Info
		if parseErr != nil {
			info = client.Resolve(c, nil)
		} else {
			info = client.Resolve(c, uaParser)
		}
		client.SetLocals(c, info)
		return c.Next()
	}
}
