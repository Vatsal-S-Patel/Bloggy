package draftapi

import (
	"errors"
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/internal/dto"
	"github.com/Vatsal-S-Patel/Bloggy/internal/errs"
	"github.com/Vatsal-S-Patel/Bloggy/internal/utils"
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
	var addDraftRequest *dto.AddDraftRequest
	err := c.BodyParser(&addDraftRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse add draft request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid add draft request",
		})
	}

	err = api.app.Validator.Struct(addDraftRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate add draft request:" + err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Title":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Title is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Title cannot be more than 130 characters long",
					})
				}
			case "Subtitle":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Subtitle is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Subtitle cannot be more than 170 characters long",
					})
				}
			case "Content":
				return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
					Message: "Content cannot be empty",
				})
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate add draft request",
		})
	}

	authorID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	draft := &models.Draft{
		ID:        uuid.New(),
		Title:     addDraftRequest.Title,
		Subtitle:  addDraftRequest.Subtitle,
		Content:   addDraftRequest.Content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = api.app.DraftService.Add(draft)
	if err != nil {
		api.app.Logger.Error("failed to add draft:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Draft saved successfully",
	})
}

func (api *api) GetAll(c *fiber.Ctx) error {
	authorID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	drafts, err := api.app.DraftService.GetAll(authorID)
	if err != nil {
		if errors.Is(err, errs.ErrDraftNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "You have no drafts",
			})
		}
		api.app.Logger.Error("failed to fetch drafts:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched drafts successfully",
		Data:    drafts,
	})
}

func (api *api) Get(c *fiber.Ctx) error {
	draftID, err := uuid.Parse(c.Params("draftID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Draft ID is not valid",
		})
	}

	draft, err := api.app.DraftService.Get(draftID)
	if err != nil {
		if errors.Is(err, errs.ErrDraftNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Draft not found",
			})
		}
		api.app.Logger.Error("failed to fetch draft by id:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched draft successfully",
		Data:    draft,
	})
}

func (api *api) Update(c *fiber.Ctx) error {
	draftID, err := uuid.Parse(c.Params("draftID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Draft ID is not valid",
		})
	}

	var updateDraftRequest *dto.AddDraftRequest
	err = c.BodyParser(&updateDraftRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse update draft request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid update draft request",
		})
	}

	err = api.app.Validator.Struct(updateDraftRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate update draft request:" + err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Title":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Title is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Title cannot be more than 130 characters long",
					})
				}
			case "Subtitle":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Subtitle is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Subtitle cannot be more than 170 characters long",
					})
				}
			case "Content":
				return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
					Message: "Content cannot be empty",
				})
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate update draft request",
		})
	}

	draft := &models.Draft{
		ID:        draftID,
		Title:     updateDraftRequest.Title,
		Subtitle:  updateDraftRequest.Subtitle,
		Content:   updateDraftRequest.Content,
		UpdatedAt: time.Now(),
	}

	err = api.app.DraftService.Update(draft)
	if err != nil {
		if errors.Is(err, errs.ErrDraftNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Draft not found",
			})
		}
		api.app.Logger.Error("failed to add draft:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Draft updated successfully",
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	draftID, err := uuid.Parse(c.Params("draftID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Draft ID is not valid",
		})
	}

	err = api.app.DraftService.Remove(draftID)
	if err != nil {
		if errors.Is(err, errs.ErrDraftNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Draft not found",
			})
		}
		api.app.Logger.Error("failed to remove draft:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
