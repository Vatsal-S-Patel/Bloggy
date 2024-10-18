package tagapi

import (
	"errors"

	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/internal/dto"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/go-playground/validator/v10"
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
	var addTagRequest *dto.AddTagRequest
	err := c.BodyParser(&addTagRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse add tag request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid add tag request",
		})
	}

	err = api.app.Validator.Struct(addTagRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate add tag request:" + err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Name is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Name cannot be more than 50 characters long",
					})
				}
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate add tag request",
		})
	}

	tag := &models.Tag{
		ID:   uuid.New(),
		Name: addTagRequest.Name,
	}

	err = api.app.TagService.Add(tag)
	if err != nil {
		if errors.Is(err, errs.ErrTagAlreadyInUse) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Tag already exists",
			})
		}
		api.app.Logger.Error("failed to add tag:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Tag added successfully",
	})
}

func (api *api) Get(c *fiber.Ctx) error {
	tag, err := api.app.TagService.Get(c.Params("tag"))
	if err != nil {
		if errors.Is(err, errs.ErrTagNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Tag not found",
			})
		}
		api.app.Logger.Error("failed to fetch tag:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched tag successfully",
		Data:    tag,
	})
}
