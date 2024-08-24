package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"time"
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

func (repo *MessageRepository) GetNewMessages() ([]models.Message, error) {
	return repo.db.GetMessageByStatus("new")
}

func (repo *MessageRepository) UpdateMessage(message *models.Message) error {
	return repo.db.UpdateMessage(message)
}

func (repo *MessageRepository) GetStuckMessages(timeout time.Duration) ([]models.Message, error) {
	// Получаем текущее время минус таймаут
	cutoffTime := time.Now().Add(-timeout)

	// Используем метод базы данных для получения сообщений по статусу и времени
	return repo.db.GetMessagesByStatusAndTime("processing", cutoffTime)
}
