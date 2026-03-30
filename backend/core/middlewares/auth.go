package middlewares

import (
	"strings"
	"time"
	"wms/core/utils"
	"wms/internal/model"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	TokenUtil *utils.TokenUtil
}

type AuthMiddlewareConfig struct {
	TokenUtil *utils.TokenUtil
}

func NewAuthMiddleware(cfg AuthMiddlewareConfig) *AuthMiddleware {
	return &AuthMiddleware{
		TokenUtil: cfg.TokenUtil,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := m.extractToken(ctx)
		if err != nil {
			return err
		}

		claims, err := m.TokenUtil.ParseToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorize")
		}

		if err := m.validateExpiration(claims); err != nil {
			return err
		}

		auth := &model.Auth{
			UID:    claims.UID,
			Claims: claims,
		}
		ctx.Locals("auth", auth)

		return ctx.Next()
	}
}

func (m *AuthMiddleware) extractToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Empty token")
	}

	return token, nil
}

func (m *AuthMiddleware) validateExpiration(claims *model.TokenClaims) error {
	if claims.Exp < time.Now().Unix() {
		return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
	}
	return nil
}

// GetAuth is a helper to retrieve the authenticated user from fiber context.
func GetAuth(ctx *fiber.Ctx) *model.Auth {
	auth, ok := ctx.Locals("auth").(*model.Auth)
	if !ok {
		return nil
	}
	return auth
}
