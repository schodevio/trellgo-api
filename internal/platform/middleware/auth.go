package middleware

import (
	"strings"

	paseto "aidanwoods.dev/go-paseto"
	"github.com/gofiber/fiber/v3"
	"github.com/schodevio/trellgo/internal/platform/apierrors"
)

func Authenticate(key paseto.V4SymmetricKey) fiber.Handler {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	return func(c fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization)
		if header == "" {
			return apierrors.Unauthorized("missing authorization header")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return apierrors.Unauthorized("invalid authorization header format")
		}

		token, err := parser.ParseV4Local(key, parts[1], nil)
		if err != nil {
			return apierrors.Unauthorized("invalid or expired token")
		}

		typ, err := token.GetString("typ")
		if err != nil || typ != "access" {
			return apierrors.Unauthorized("invalid token type")
		}

		userID, err := token.GetSubject()
		if err != nil {
			return apierrors.Unauthorized("invalid token claims")
		}

		c.Locals("user_id", userID)

		return c.Next()
	}
}
