package authorization

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TokenContainer struct {
	Sender string
}

type RefreshTokenContainer struct {
	TC TokenContainer
}

type AuthToken interface {
	GenerateToken(tc TokenContainer) (string, error)
	ValidateToken(tokenString string) (*TokenContainer, error)
	GenerateRefreshToken(rtc RefreshTokenContainer) (string, error)
	ValidateRefreshToken(rtString string) (*RefreshTokenContainer, error)
}

func AuthMiddleware(authToken AuthToken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authorization header is missing"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token format"})
			}

			payload, err := authToken.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
			}

			c.Set("id", payload.Sender)
			return next(c)
		}
	}
}
