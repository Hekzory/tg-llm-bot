package main

import (
	"Hekzory/tg-llm-bot/go/model-service/app/config"
	"Hekzory/tg-llm-bot/go/model-service/app/handler"
	"Hekzory/tg-llm-bot/go/model-service/app/repository"
	"Hekzory/tg-llm-bot/go/model-service/app/service"
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
	_ "embed"
)

//go:embed sql/init.sql
var initSQLstring string

func main() {

	logger, _ := logging.NewLogger("DEBUG")

	// Load Configuration
	var cfg config.ServiceConfig
	err := cfg.LoadConfig(logger)
	if err != nil {
		logger.Fatal("Error while loading config: %s", err)
	}
	logger.Info("Config loaded successfully")

	// Initialize Database
	db, err := database.NewDatabase(cfg.Config.DatabaseUrl, logger)
	if err != nil {
		logger.Fatal("Error while loading database: %s", err)
	}

	err = db.InitializeTables(initSQLstring)
	if err != nil {
		logger.Fatal("Error initializing database: %s", err)
	}

	// Initialize Repository
	userRepo := repository.NewUserRepository(db, logger)
	messageRepo := repository.NewMessageRepository(db, logger)
	conversationRepo := repository.NewConversationRepository(db, logger)

	// Initialize Service
	svc := service.NewModelService(userRepo, messageRepo, conversationRepo, logger)

	// Initialize Handler
	handler := handler.NewModelHandler(svc, logger, &cfg)

	// Start Server
	handler.StartServer()

}
