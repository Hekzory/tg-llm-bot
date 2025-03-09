package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/database/models"
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

func (repo *ConversationRepository) StartNewConversation(conversation *models.Conversation) (error) {
	return repo.db.StartNewConversation(conversation)
}

func (repo *ConversationRepository) ConvExists(conversation *models.Conversation) (bool, error) {
	return repo.db.ConvExists(conversation)
}

func (repo *ConversationRepository) GetConvIdByUserId(userId int) (int, error) {
	return repo.db.GetConvIdByUserId(userId)
}

func (repo *ConversationRepository) GetTgIdByConvId(convId int) (int, error) {
	return repo.db.GetTgIdByConvId(convId)
}