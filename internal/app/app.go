package app

import (
	"fmt"
	"log"
	"os"

	"github.com/Vatsal-S-Patel/Bloggy/internal/consts"
	"github.com/Vatsal-S-Patel/Bloggy/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type App struct {
	DB *sqlx.DB

	Config *models.Config
	Logger *zap.Logger
}

func New() (*App, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Println("ERROR failed to initialize logger:", err.Error())
		return &App{}, err
	}

	configFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		logger.Error("failed to read config.yaml:", zap.String("error", err.Error()))
		return &App{}, err
	}

	config := &models.Config{}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		logger.Error("failed to unmarshal config file:", zap.String("error", err.Error()))
		return &App{}, err
	}

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", config.PostgresConfig.Host, config.PostgresConfig.User, config.PostgresConfig.Password, config.PostgresConfig.Database, config.PostgresConfig.SSLMode)

	db, err := sqlx.Open(consts.DB_DRIVER, psqlInfo)
	if err != nil {
		logger.Error("failed to establish postgresql database connection", zap.String("error", err.Error()))
		return &App{}, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("failed to ping database connection", zap.String("error", err.Error()))
		return &App{}, err
	}

	logger.Info("postgresql database connection established")

	app := &App{
		DB:     db,
		Config: config,
		Logger: logger,
	}

	return app, nil
}
