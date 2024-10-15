package api

import (
	"os"
	"os/signal"
	"time"

	"github.com/Vatsal-S-Patel/Bloggy/internal/api/blogapi"
	"github.com/Vatsal-S-Patel/Bloggy/internal/api/draftapi"
	"github.com/Vatsal-S-Patel/Bloggy/internal/api/healthapi"
	"github.com/Vatsal-S-Patel/Bloggy/internal/api/userapi"
	"github.com/Vatsal-S-Patel/Bloggy/internal/app"
	"github.com/Vatsal-S-Patel/Bloggy/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ListenAndServe(app *app.App) {
	fiberApp := fiber.New()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan

		// Clear logger buffer	// Not needed for stdout
		// err := app.Logger.Sync()
		// if err != nil {
		// 	app.Logger.Error("failed to sync logger", zap.String("error", err.Error()))
		// 	return
		// }

		err := app.DB.Close()
		if err != nil {
			app.Logger.Error("failed to close postgresql database connection", zap.String("error", err.Error()))
			return
		}
		app.Logger.Info("postgresql database connection closed")

		app.Logger.Info("server stopped gracefully")

		err = fiberApp.ShutdownWithTimeout(15 * time.Second)
		if err != nil {
			app.Logger.Error("failed to shutdown http server", zap.String("error", err.Error()))
			return
		}
	}()

	RegisterRoutes(fiberApp, app)

	err := fiberApp.Listen(":" + app.Config.ServerConfig.Port)
	if err != nil {
		app.Logger.Error("failed to initialze server", zap.String("error", err.Error()))
		return
	}
}

func RegisterRoutes(fiberApp *fiber.App, app *app.App) {
	healthAPI := healthapi.New(app)
	userAPI := userapi.New(app)
	blogAPI := blogapi.New(app)
	draftAPI := draftapi.New(app)

	userAuthMiddleware := middlewares.UserAuthMiddleware(app)

	router := fiberApp.Group("/v1")

	router.Get("/health", healthAPI.Check)

	userRouter := router.Group("/users")
	userRouter.Post("/register", userAPI.Register)
	userRouter.Post("/login", userAPI.Login)

	blogRouter := router.Group("/blogs")
	blogRouter.Post("/", userAuthMiddleware, blogAPI.Publish)

	draftRouter := router.Group("/drafts")
	draftRouter.Post("/", userAuthMiddleware, draftAPI.Add)
	draftRouter.Get("/", userAuthMiddleware, draftAPI.GetAll)
	draftRouter.Get("/:draftID", userAuthMiddleware, draftAPI.Get)
	draftRouter.Put("/:draftID", userAuthMiddleware, draftAPI.Update)
	draftRouter.Delete("/:draftID", userAuthMiddleware, draftAPI.Remove)
}
