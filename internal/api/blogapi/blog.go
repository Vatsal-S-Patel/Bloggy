package blogapi

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

func (api *api) Publish(c *fiber.Ctx) error {

	var publishBlogRequest *dto.PublishBlogRequest
	err := c.BodyParser(&publishBlogRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse publish blog request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid publish blog request",
		})
	}

	err = api.app.Validator.Struct(publishBlogRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate publish blog request:" + err.Error())
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
			case "FtImage":
				return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
					Message: "Invalid feature image format",
				})
			case "Tags":
				switch err.Tag() {
				case "dive":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Tags are required",
					})
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Tag is required",
					})
				case "uuid":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Tag ID is not valid",
					})
				}
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate publish blog request",
		})
	}

	if len(publishBlogRequest.Tags) > 10 {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog cannot have more than 10 tags",
		})
	}

	authorID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	blog := &models.Blog{
		ID:        uuid.New(),
		Title:     publishBlogRequest.Title,
		Subtitle:  publishBlogRequest.Subtitle,
		Content:   publishBlogRequest.Content,
		FtImage:   publishBlogRequest.FtImage,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	blogTags := make([]*models.BlogTag, 0, len(publishBlogRequest.Tags))

	for _, tagID := range publishBlogRequest.Tags {
		blogTags = append(blogTags, &models.BlogTag{
			BlogID: blog.ID,
			TagID:  tagID,
		})
	}

	err = api.app.BlogService.Publish(blog, blogTags)
	if err != nil {
		api.app.Logger.Error("failed to publish blog:" + err.Error())
		if errors.Is(err, errs.ErrTagNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Tag not found",
			})
		}
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Blog published successfully",
	})
}

func (api *api) Get(c *fiber.Ctx) error {
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

	blog, err := api.app.BlogService.Get(blogID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Blog not found",
			})
		}
		api.app.Logger.Error("failed to fetch blog by id:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched blog successfully",
		Data:    blog,
	})
}
