package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

type ModelRepository struct {
	db     *database.DB
	logger *logging.Logger
}

func NewModelRepository(db *database.DB, logger *logging.Logger) *ModelRepository {
	return &ModelRepository{
		db:     db,
		logger: logger,
	}
}
