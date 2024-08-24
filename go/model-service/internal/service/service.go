package service

import (
	"Hekzory/tg-llm-bot/go/model-service/internal/repository"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"context"
	"time"
)

type ModelService struct {
	userRepo    *repository.UserRepository
	messageRepo *repository.MessageRepository
	logger      *logging.Logger
}

func NewModelService(userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, logger *logging.Logger) *ModelService {
	return &ModelService{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		logger:      logger,
	}
}

func (s *ModelService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *ModelService) GetNewMessages(ctx context.Context) ([]models.Message, error) {
	return s.messageRepo.GetNewMessages()
}

func (s *ModelService) UpdateMessageStatus(ctx context.Context, messageID int, status string) error {
	err := s.messageRepo.UpdateMessageStatus(messageID, status)
	if err != nil {
		return err
	}
	s.logger.Info("Updated message status: messageID=%d, status=%s", messageID, status)
	return nil
}

func (s *ModelService) UpdateMessage(ctx context.Context, message models.Message) error {

	message.Status = "ready"
	err := s.messageRepo.UpdateMessage(&message)

	if err != nil {
		return err
	}

	s.logger.Info("Updated message: %+v", message)
	return nil
}

func (s *ModelService) GetStuckMessages(ctx context.Context, timeout time.Duration) ([]models.Message, error) {
	return s.messageRepo.GetStuckMessages(timeout)
}
