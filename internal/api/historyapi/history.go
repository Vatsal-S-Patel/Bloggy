package historyapi

import (
	"errors"

	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/internal/utils"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type api struct {
	app *app.App
}

func New(app *app.App) *api {
	return &api{
		app: app,
	}
}

func (api *api) Get(c *fiber.Ctx) error {
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	history, err := api.app.HistoryService.Get(userID)
	if err != nil {
		if errors.Is(err, errs.ErrHistoryNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "You have no history",
			})
		}
		api.app.Logger.Error("failed to fetch history:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched reading history",
		Data:    history,
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "BlogID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.HistoryService.Remove(userID, blogID)
	if err != nil {
		if errors.Is(err, errs.ErrHistoryNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "History not found",
			})
		}
		api.app.Logger.Error("failed to remove blog from history:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, models.Response{
		Message: "Removed blog from reading history",
	})
}

func (api *api) RemoveAll(c *fiber.Ctx) error {
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.HistoryService.RemoveAll(userID)
	if err != nil {
		if errors.Is(err, errs.ErrHistoryNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "You have no history",
			})
		}
		api.app.Logger.Error("failed to remove all history:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, models.Response{
		Message: "Removed all reading history",
	})
}
