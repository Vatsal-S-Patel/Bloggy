package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  userID,
			"exp": time.Now().Add(20 * time.Minute).Unix(),
		},
	)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ExtractJWTClaims(c *fiber.Ctx) map[string]interface{} {
	token := c.Locals("claims").(*jwt.Token)
	return token.Claims.(jwt.MapClaims)
}
