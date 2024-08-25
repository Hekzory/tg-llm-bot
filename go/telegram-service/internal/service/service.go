package service

import (
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"Hekzory/tg-llm-bot/go/telegram-service/internal/repository"
	"context"
	"database/sql"
)

type UserService struct {
	userRepo    *repository.UserRepository
	messageRepo *repository.MessageRepository
	convRepo    *repository.ConversationRepository
	logger      *logging.Logger
}

func NewTelegramService(userRepo *repository.UserRepository, messageRepo *repository.MessageRepository, convRepo *repository.ConversationRepository, logger *logging.Logger) *UserService {
	return &UserService{
		userRepo:    userRepo,
		messageRepo: messageRepo,
		convRepo:    convRepo,
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
	s.logger.Info("Updated message status: messageID=%d, status=%s", messageID, status)
	return nil
}

func (s *UserService) GetReadyMessages(ctx context.Context) ([]models.Message, error) {
	return s.messageRepo.GetReadyMessages()
}

func (s *UserService) GetProcessingMessages(ctx context.Context) ([]models.Message, error) {
	return s.messageRepo.GetProcessingMessages()
}

func (s *UserService) AddMessage(ctx context.Context, questiont string, convID int, tgQuestionedID int) (error, int) {
	return s.messageRepo.AddMessage(&models.Message{
		Question:       questiont,
		Status:         "new",
		ConversationID: sql.NullInt64{Int64: int64(convID), Valid: true},
		TgAnswerId:     sql.NullInt64{Int64: int64(0), Valid: true},
		TgQuestionId:   sql.NullInt64{Int64: int64(tgQuestionedID), Valid: true},
	})

}

func (s *UserService) UserExists(ctx context.Context, tg_id int) (bool, error) {
	return s.userRepo.UserExists(&models.User{TelegramID: tg_id})
}

func (s *UserService) AddUser(ctx context.Context, tg_id int, name string, username string) error {
	return s.userRepo.AddUser(&models.User{
		TelegramID: tg_id,
		Name:       name,
		Username:   username,
		IsPremium:  false,
	})
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserIdByTgId(ctx context.Context, tg_id int) (int, error) {
	return s.userRepo.GetUserIdByTgId(tg_id)
}

func (s *UserService) GetTgIdByConvId(ctx context.Context, convId int) (int, error) {
	return s.convRepo.GetTgIdByConvId(convId)
}

func (s *UserService) StartNewConversation(ctx context.Context, id int) error {
	return s.convRepo.StartNewConversation(&models.Conversation{
		UserID: id,
	})
}

func (s *UserService) ConvExists(ctx context.Context, userId int) (bool, error) {
	return s.convRepo.ConvExists(&models.Conversation{UserID: userId})
}

func (s *UserService) GetConvIdByUserId(ctx context.Context, userId int) (int, error) {
	return s.convRepo.GetConvIdByUserId(userId)
}
