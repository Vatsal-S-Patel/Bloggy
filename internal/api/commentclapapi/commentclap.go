package commentclapapi

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
	commentID, err := uuid.Parse(c.Params("commentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Comment ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.CommentClapService.Add(commentID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyClapped) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Comment already clapped",
			})
		} else if errors.Is(err, errs.ErrCommentNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Comment not found",
			})
		}
		api.app.Logger.Error("failed to add clap for comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Comment clapped successfully",
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("commentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Comment ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.CommentClapService.Remove(commentID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyUnClapped) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Comment already unclapped",
			})
		}
		api.app.Logger.Error("failed to remove clap for comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
