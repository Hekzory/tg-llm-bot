package main

import (
	"Hekzory/tg-llm-bot/go/model-service/internal/config"
	"Hekzory/tg-llm-bot/go/model-service/internal/handler"
	"Hekzory/tg-llm-bot/go/model-service/internal/repository"
	"Hekzory/tg-llm-bot/go/model-service/internal/service"
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

func main() {

	logger, _ := logging.NewLogger("DEBUG")

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config loaded successfully")

	// Initialize Database
	db := database.NewDatabase(cfg.DatabaseURL)
	logger.Info("Database loaded successfully")

	// Initialize Repository
	repo := repository.NewModelRepository(db, logger)

	// Initialize Service
	svc := service.NewModelService(repo, logger)

	// Initialize Handler
	handler := handler.NewModelHandler(svc, logger)

	// Start Server
	handler.StartServer(cfg.ServerPort)

}
