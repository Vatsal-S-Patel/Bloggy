package readlaterapi

import (
	"errors"
	"time"

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

	readLaterBlog := &models.ReadLater{
		UserID:    userID,
		BlogID:    blogID,
		CreatedAt: time.Now(),
	}

	err = api.app.ReadLaterService.Add(readLaterBlog)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyInReadLater) {
			return models.SendResponse(c, fiber.StatusOK, models.Response{
				Message: "Blog already in read later",
			})
		}
		api.app.Logger.Error("failed to add blog into read later:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Added to read later successfully",
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

	err = api.app.ReadLaterService.Remove(userID, blogID)
	if err != nil {
		if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Blog not found",
			})
		}
		api.app.Logger.Error("failed to remove blog from read later:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
