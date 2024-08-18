package repository

import (
	"Hekzory/tg-llm-bot/go/shared/database"
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
)

// UserRepository предоставляет методы для работы с пользователями в базе данных
type UserRepository struct {
	db     *database.DB
	logger *logging.Logger
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *database.DB, logger *logging.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

// GetUserByID получает пользователя по его ID
func (repo *UserRepository) GetUserByID(id int) (*models.User, error) {
	return repo.db.GetUserByID(id)
}

// UpdateUser обновляет существующего пользователя
func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.db.UpdateUser(user)
}

// DeleteUser удаляет пользователя из базы данных
func (repo *UserRepository) DeleteUser(id int) error {
	return repo.db.DeleteUser(id)
}

// GetAllUsers получает всех пользователей из базы данных
func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	return repo.db.GetAllUsers()
}

// AddUser добавляет нового пользователя в базу данных
func (repo *UserRepository) AddUser(user *models.User) error {
	return repo.db.AddUser(user)
}

func (repo *UserRepository) IsUserExist(user *models.User) (bool, error) {
	return repo.db.IsUserExist(user)
}

func (repo *UserRepository) GetUserIdByTgId(tg_id int) (int, error) {
	return repo.db.GetUserIdByTgId(tg_id)
}

func (repo *UserRepository) GetTgIdByUserId(id int) (int, error) {
	return repo.db.GetTgIdByUserId(id)
}
