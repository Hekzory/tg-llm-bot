package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

type ConversationRepository struct {
	db     *database.DB
	logger *logging.Logger
}

func NewConversationRepository(db *database.DB, logger *logging.Logger) *ConversationRepository {
	return &ConversationRepository{
		db:     db,
		logger: logger,
	}
}
