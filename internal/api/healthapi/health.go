package healthapi

import (
	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/gofiber/fiber/v2"
)

type api struct {
	app *app.App
}

func New(app *app.App) *api {
	return &api{
		app: app,
	}
}

func (api *api) Check(c *fiber.Ctx) error {
	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Server Health OK",
	})
}
