package commentapi

import (
	"errors"
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/internal/dto"
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

	var addCommentRequest *dto.AddCommentRequest
	err = c.BodyParser(&addCommentRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse add comment request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid add comment request",
		})
	}

	if addCommentRequest.Body == "" {
		api.app.Logger.Error("failed to validate add comment request")
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Body is required",
		})
	}

	authorID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	comment := &models.Comment{
		ID:        uuid.New(),
		Body:      addCommentRequest.Body,
		Claps:     0,
		Replies:   0,
		AuthorID:  authorID,
		BlogID:    blogID,
		CreatedAt: time.Now(),
	}

	err = api.app.CommentService.Add(comment)
	if err != nil {
		api.app.Logger.Error("failed to add comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Comment added successfully",
	})
}

func (api *api) AddReply(c *fiber.Ctx) error {

	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	parentID, err := uuid.Parse(c.Params("parentCommentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	var addCommentRequest *dto.AddCommentRequest
	err = c.BodyParser(&addCommentRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse add comment request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid add comment request",
		})
	}

	if addCommentRequest.Body == "" {
		api.app.Logger.Error("failed to validate add comment request")
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Body is required",
		})
	}

	authorID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	comment := &models.Comment{
		ID:        uuid.New(),
		Body:      addCommentRequest.Body,
		Claps:     0,
		Replies:   0,
		AuthorID:  authorID,
		BlogID:    blogID,
		ParentID:  parentID,
		CreatedAt: time.Now(),
	}

	err = api.app.CommentService.AddReply(comment)
	if err != nil {
		if errors.Is(err, errs.ErrCommentNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Comment not found",
			})
		}
		api.app.Logger.Error("failed to add comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "Reply added successfully",
	})
}

func (api *api) Get(c *fiber.Ctx) error {
	blogID, err := uuid.Parse(c.Params("blogID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Blog ID is not valid",
		})
	}

	comments, err := api.app.CommentService.Get(blogID)
	if err != nil {
		if errors.Is(err, errs.ErrCommentNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "No comments! Be first to comment",
			})
		}
		api.app.Logger.Error("failed to get top level comments:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Comments fetched successfully",
		Data:    comments,
	})
}

func (api *api) GetReplies(c *fiber.Ctx) error {
	parentID, err := uuid.Parse(c.Params("parentCommentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Parent comment ID is not valid",
		})
	}

	comments, err := api.app.CommentService.GetReplies(parentID)
	if err != nil {
		if errors.Is(err, errs.ErrCommentNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Comment not found",
			})
		}
		api.app.Logger.Error("failed to get replies of comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Comments fetched successfully",
		Data:    comments,
	})
}

func (api *api) Update(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("commentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Comment ID is not valid",
		})
	}

	var updateCommentRequest *dto.UpdateCommentRequest
	err = c.BodyParser(&updateCommentRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse update comment request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid update comment request",
		})
	}

	if updateCommentRequest.Body == "" {
		api.app.Logger.Error("failed to validate update comment request")
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Body is required",
		})
	}

	userID, err := utils.ExtractUserIDFromContext(c)
	if err != nil {
		api.app.Logger.Error("failed to extract user id from context:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	comment := &models.Comment{
		ID:       commentID,
		Body:     updateCommentRequest.Body,
		AuthorID: userID,
	}

	err = api.app.CommentService.Update(comment)
	if err != nil {
		if errors.Is(err, errs.ErrCommentNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "Comment not found",
			})
		}
		api.app.Logger.Error("failed to update comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "Comment updated successfully",
	})
}

func (api *api) Remove(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("commentID"))
	if err != nil {
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Comment ID is not valid",
		})
	}

	err = api.app.CommentService.Remove(commentID)
	if err != nil {
		api.app.Logger.Error("failed to remove comment:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusNoContent, nil)
}
