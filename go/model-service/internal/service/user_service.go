package service

import (
	"Hekzory/tg-llm-bot/go/model-service/internal/repository"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"context"
	"fmt"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *logging.Logger
}

func NewModelService(repo *repository.UserRepository, logger *logging.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}
func (s *UserService) ProcessRequests(ctx context.Context) error {
	// Получение всех пользователей из базы данных
	users, err := s.repo.GetAllUsers()
	if err != nil {
		s.logger.Error("Failed to get users: ", err)
		return err
	}

	for _, user := range users {
		s.logger.Info(fmt.Sprintf("Got user: %+v", user))
	}

	return nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllUsers()
}
