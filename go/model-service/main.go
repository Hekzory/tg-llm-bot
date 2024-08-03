package main

import (
	"Hekzory/tg-llm-bot/go/model-service/internal/config"
	"Hekzory/tg-llm-bot/go/model-service/internal/handler"
	"Hekzory/tg-llm-bot/go/model-service/internal/repository"
	"Hekzory/tg-llm-bot/go/model-service/internal/service"
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
	_ "embed"
)

//go:embed sql/init.sql
var initSQLstring string

func main() {

	logger, _ := logging.NewLogger("DEBUG")

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config loaded successfully")

	// Initialize Database
	db, err := database.NewDatabase(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("Error while loading database: %s", err)
	}

	err = db.InitializeTables(initSQLstring)
	if err != nil {
		logger.Fatal("Error initializing database: %s", err)
	}

	// Initialize Repository
	repo := repository.NewModelRepository(db, logger)

	// Initialize Service
	svc := service.NewModelService(repo, logger)

	// Initialize Handler
	handler := handler.NewModelHandler(svc, logger)

	// Start Server
	handler.StartServer(cfg.ServerPort)

}
