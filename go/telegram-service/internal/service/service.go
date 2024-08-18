package service

import (
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/repository"
	"context"
	"fmt"
)

type UserService struct {
	userRepo    *repository.UserRepository
	messageRepo *repository.MessageRepository
	logger      *logging.Logger
}

func NewTelegramService(userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, logger *logging.Logger) *UserService {
	return &UserService{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		logger:      logger,
	}
}

func (s *UserService) ProcessRequests(ctx context.Context) error {
	return nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) UpdateMessageStatus(ctx context.Context, messageID int, status string) error {
	err := s.messageRepo.UpdateMessageStatus(messageID, status)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Updated message status: messageID=%d, status=%s", messageID, status))
	return nil
}

func (s *UserService) GetReadyMessages(ctx context.Context) ([]models.Message, error) {
	return s.messageRepo.GetReadyMessages()
}

func (s *UserService) AddMessage(ctx context.Context, userId int, questiont string) error {
	return s.messageRepo.AddMessage(&models.Message{UserID: userId, Question: questiont, Status: "new"})

}

func (s *UserService) UserExists(ctx context.Context, tg_id int) (bool, error) {
	return s.userRepo.UserExists(&models.User{TelegramID: tg_id})
}

func (s *UserService) AddUser(ctx context.Context, tg_id int, name string, username string) error {
	return s.userRepo.AddUser(&models.User{
		TelegramID: tg_id,
		Name: name,
		Username: username,
		IsPremium: false,
	})
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserIdByTgId(ctx context.Context, tg_id int) (int, error) {
	return s.userRepo.GetUserIdByTgId(tg_id)
}

func (s *UserService) GetTgIdByUserId(ctx context.Context, tg_id int) (int, error) {
	return s.userRepo.GetTgIdByUserId(tg_id)
}
