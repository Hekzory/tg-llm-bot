package repository

import (
	"Hekzory/tg-llm-bot/go/shared/logging"
	"database/sql"
)

type ModelRepository struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewModelRepository(db *sql.DB, logger *logging.Logger) *ModelRepository {
	return &ModelRepository{
		db:     db,
		logger: logger,
	}
}
