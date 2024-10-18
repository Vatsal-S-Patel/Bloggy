package bookmarkapi

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
	var addBookmarkRequest *dto.AddBookmarkRequest
	err := c.BodyParser(&addBookmarkRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse add bookmark request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid add bookmark request",
		})
	}

	err = api.app.Validator.Struct(addBookmarkRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate add bookmark request:" + err.Error())
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
						Message: "Name cannot be more than 30 characters long",
					})
				}
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate add bookmark request",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	bookmark := &models.Bookmark{
		ID:        uuid.New(),
		Name:      addBookmarkRequest.Name,
		UserID:    userID,
		CreatedAt: time.Now(),
		Visible:   addBookmarkRequest.Visible,
	}

	err = api.app.BookmarkService.Add(bookmark)
	if err != nil {
		if errors.Is(err, errs.ErrBookmarkNameAlreadyInUse) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Bookmark with this name already exists",
			})
		}
		api.app.Logger.Error("failed to add bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Added bookmark successfully",
	})
}

func (api *api) AddBlog(c *fiber.Ctx) error {
	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	err = api.app.BookmarkService.AddBlog(bookmarkID, blogID)
	if err != nil {
		if errors.Is(err, errs.ErrBlogAlreadyInBookmark) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Blog is already in this bookmark",
			})
		} else if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Blog not found",
			})
		}
		api.app.Logger.Error("failed to add blog into bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Added blog into bookmark successfully",
	})
}

func (api *api) GetAll(c *fiber.Ctx) error {
	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	bookmarks, err := api.app.BookmarkService.GetAll(userID)
	if err != nil {
		if errors.Is(err, errs.ErrBookmarkNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "You have no bookmarks",
			})
		}
		api.app.Logger.Error("failed to fetch all bookmarks:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched all bookmarks successfully",
		Data:    bookmarks,
	})
}

func (api *api) Get(c *fiber.Ctx) error {
	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	bookmark, err := api.app.BookmarkService.Get(bookmarkID)
	if err != nil {
		if errors.Is(err, errs.ErrBookmarkNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Bookmark not found",
			})
		}
		api.app.Logger.Error("failed to fetch bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched bookmark successfully",
		Data:    bookmark,
	})
}

func (api *api) GetBlogs(c *fiber.Ctx) error {
	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	blogs, err := api.app.BookmarkService.GetBlogs(bookmarkID)
	if err != nil {
		if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "You have no blogs",
			})
		}
		api.app.Logger.Error("failed to fetch bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Fetched all blogs from bookmark successfully",
		Data:    blogs,
	})
}

func (api *api) Update(c *fiber.Ctx) error {
	var updateBookmarkRequest *dto.AddBookmarkRequest
	err := c.BodyParser(&updateBookmarkRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse update bookmark request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid update bookmark request",
		})
	}

	err = api.app.Validator.Struct(updateBookmarkRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate update bookmark request:" + err.Error())
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
						Message: "Name cannot be more than 30 characters long",
					})
				}
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate update bookmark request",
		})
	}

	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.BookmarkService.Update(updateBookmarkRequest, bookmarkID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrBookmarkNameAlreadyInUse) {
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Bookmark with this name already exists",
			})
		} else if errors.Is(err, errs.ErrBookmarkNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Bookmark not found",
			})
		}
		api.app.Logger.Error("failed to update bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Bookmark updated successfully",
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	err = api.app.BookmarkService.Remove(bookmarkID, userID)
	if err != nil {
		if errors.Is(err, errs.ErrBookmarkNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Bookmark not found",
			})
		}
		api.app.Logger.Error("failed to remove bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}

func (api *api) RemoveBlog(c *fiber.Ctx) error {
	bookmarkID, err := uuid.Parse(c.Params("bookmarkID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Bookmark ID is not valid",
		})
	}

	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	err = api.app.BookmarkService.RemoveBlog(bookmarkID, blogID)
	if err != nil {
		if errors.Is(err, errs.ErrBlogNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Either bookmark or blog not found",
			})
		}
		api.app.Logger.Error("failed to remove blog from bookmark:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
