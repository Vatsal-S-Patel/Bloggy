package models

import "github.com/gofiber/fiber/v2"

// Response
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendResponse will send response with mentioned statusCode and data
func SendResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}
