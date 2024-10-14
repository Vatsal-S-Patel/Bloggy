package userapi

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

func (api *api) Register(c *fiber.Ctx) error {

	var userRegistrationRequest *dto.UserRegistrationRequest
	err := c.BodyParser(&userRegistrationRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse user registration request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid user registration request",
		})
	}

	err = api.app.Validator.Struct(userRegistrationRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate user registration request:" + err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Username is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Username cannot be more than 30 characters long",
					})
				}
			case "Email":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Email is required",
					})
				case "email":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Invalid email format",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Email cannot be more than 70 characters long",
					})
				}
			case "Password":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password is required",
					})
				case "min":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password must be atleast 8 characters long",
					})
				case "password":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password must contain atleast one special character",
					})
				}
			case "Bio":
				return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
					Message: "Bio cannot be more than 500 characters",
				})
			case "Avatar":
				return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
					Message: "Invalid avatar format",
				})
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate user registration request",
		})
	}

	hashedPassword, err := utils.HashPassword(userRegistrationRequest.Password)
	if err != nil {
		api.app.Logger.Error("failed to generate hash password:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	user := &models.User{
		ID:          uuid.New(),
		Username:    userRegistrationRequest.Username,
		Email:       userRegistrationRequest.Email,
		Password:    hashedPassword,
		Bio:         userRegistrationRequest.Bio,
		Avatar:      userRegistrationRequest.Avatar,
		Followers:   0,
		Following:   0,
		JoinedAt:    time.Now(),
		LastLoginAt: time.Now(),
	}

	err = api.app.UserService.RegisterUser(user)
	if err != nil {
		api.app.Logger.Error("failed to register user:" + err.Error())
		switch err {
		case errs.ErrUserEmailAlreadyInUse:
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Email is already in use",
			})
		case errs.ErrUsernameAlreadyInUse:
			return models.SendResponse(c, fiber.StatusConflict, models.Response{
				Message: "Username is already in use",
			})
		}

		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	token, err := utils.GenerateJWT(user.ID.String(), []byte(api.app.Config.JWTSecret))
	if err != nil {
		api.app.Logger.Error("failed to generate jwt token for user:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusCreated, models.Response{
		Message: "User registered successfully",
		Data: dto.UserLoginResponse{
			AccessToken: token,
		},
	})
}

func (api *api) Login(c *fiber.Ctx) error {
	var loginRequest *dto.LoginRequest
	err := c.BodyParser(&loginRequest)
	if err != nil {
		api.app.Logger.Error("failed to parse user login request:" + err.Error())
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Invalid login request",
		})
	}

	err = api.app.Validator.Struct(loginRequest)
	if err != nil {
		api.app.Logger.Error("failed to validate user login request:" + err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Username is required",
					})
				case "max":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Username cannot be more than 30 characters long",
					})
				}
			case "Password":
				switch err.Tag() {
				case "required":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password is required",
					})
				case "min":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password must be atleast 8 characters long",
					})
				case "password":
					return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
						Message: "Password must contain atleast one special character",
					})
				}
			}
		}
		return models.SendResponse(c, fiber.StatusBadRequest, models.Response{
			Message: "Failed to validate user login request",
		})
	}

	userID, hashedPassword, err := api.app.UserService.GetIDPasswordByUsername(loginRequest.Username)
	if err != nil {
		api.app.Logger.Error("failed to get password by username:" + err.Error())
		if errors.Is(err, errs.ErrUserNotFound) {
			return models.SendResponse(c, fiber.StatusNotFound, models.Response{
				Message: "User with this username doesnot exist",
			})
		}
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	isValidPassword := utils.VerifyPassword(hashedPassword, loginRequest.Password)
	if !isValidPassword {
		api.app.Logger.Error("user password doesnot matched")
		return models.SendResponse(c, fiber.StatusUnauthorized, models.Response{
			Message: "Invalid credentials provided",
		})
	}

	token, err := utils.GenerateJWT(userID.String(), []byte(api.app.Config.JWTSecret))
	if err != nil {
		api.app.Logger.Error("failed to generate jwt token for user:" + err.Error())
		return models.SendResponse(c, fiber.StatusInternalServerError, models.Response{
			Message: "Internal Server Error",
		})
	}

	return models.SendResponse(c, fiber.StatusOK, models.Response{
		Message: "User login successfully",
		Data: dto.UserLoginResponse{
			AccessToken: token,
		},
	})
}
