package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

// MessageRepository предоставляет методы для работы с сообщениями в базе данных
type MessageRepository struct {
	db     *database.DB
	logger *logging.Logger
}

// NewMessageRepository создает новый экземпляр MessageRepository
func NewMessageRepository(db *database.DB, logger *logging.Logger) *MessageRepository {
	return &MessageRepository{
		db:     db,
		logger: logger,
	}
}

// GetMessageByID получает сообщение по его ID
func (repo *MessageRepository) GetMessageByID(id int) (*models.Message, error) {
	return repo.db.GetMessageByID(id)
}

func (repo *MessageRepository) UpdateMessageStatus(id int, status string) error {
	return repo.db.UpdateMessageStatus(id, status)
}

func (repo *MessageRepository) GetReadyMessages() ([]models.Message, error) {
	return repo.db.GetMessageByStatus("ready")
}

func (repo *MessageRepository) GetProcessingMessages() ([]models.Message, error) {
	return repo.db.GetMessageByStatus("processing")
}

func (repo *MessageRepository) UpdateMessage(message *models.Message) error {
	return repo.db.UpdateMessage(message)
}

func (repo *MessageRepository) AddMessage(message *models.Message) (error, int) {
	return repo.db.AddMessage(message)
}
