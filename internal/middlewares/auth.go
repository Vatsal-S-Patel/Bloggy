package middlewares

import (
	"errors"

	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func UserAuthMiddleware(app *app.App) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(app.Config.JWTSecret),
		ContextKey: "claims",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var errMessage string
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMessage = "Unauthorized access, token expired"
			} else {
				errMessage = "Unauthorized access, token invalid"
			}

			return models.SendResponse(c, fiber.StatusUnauthorized, models.Response{
				Message: errMessage,
			})
		},
	})
}
