package blogclapapi

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

func (api *api) Add(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.BlogClapService.Add(blogID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyClapped) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Blog already clapped",
			})
		} else if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Blog not found",
			})
		}
		api.app.Logger.Error("failed to add clap for blog:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Clapped blog successfully",
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.BlogClapService.Remove(blogID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyUnClapped) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Blog already unclapped",
			})
		}
		api.app.Logger.Error("failed to remove clap for blog:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
