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
	token, ok := c.Locals("claims").(*jwt.Token)
	if ok {
		return token.Claims.(jwt.MapClaims)
	}
	return nil
}

func ExtractUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	claims := ExtractJWTClaims(c)
	userID, ok := claims["id"].(string)
	if ok {
		return uuid.Parse(userID)
	}
	return uuid.Nil, nil
}
