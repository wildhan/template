package token

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(maker *PasetoMaker) echo.MiddlewareFunc {
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

			payload, err := maker.DecryptToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
			}

			expiredAt, err := payload.GetExpiration()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token expiration"})
			}

			if expiredAt.Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "token has expired"})
			}

			sub, err := payload.GetSubject()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token subject"})
			}

			c.Set("id", sub)
			return next(c)
		}
	}
}
