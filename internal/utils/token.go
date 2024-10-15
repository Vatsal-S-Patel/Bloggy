package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWT(userID string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  userID,
			"exp": time.Now().Add(1 * time.Hour).Unix(),
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

func ExtractUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	claims := ExtractJWTClaims(c)
	return uuid.Parse(claims["id"].(string))
}
