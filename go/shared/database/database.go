package database

import (
	"Hekzory/tg-llm-bot/go/shared/database/models"
	"Hekzory/tg-llm-bot/go/shared/logging"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB represents the database connection
type DB struct {
	*sqlx.DB
	logger *logging.Logger
}

// NewDatabase establishes a new database connection
func NewDatabase(databaseURL string, logger *logging.Logger) (*DB, error) {
	db, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	logger.Info("Connection created successfully")

	return &DB{DB: db, logger: logger}, nil
}

// InitializeTables runs the SQL file to create tables, expected to run on start
func (db *DB) InitializeTables(code string) error {
	_, err := db.Exec(code)
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	db.logger.Info("Tables initialized successfully")
	return nil
}

// AddUser adds a new user to the database
func (db *DB) AddUser(user *models.User) error {
	query := `INSERT INTO users (tg_id, name, username, is_premium) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return db.QueryRowx(query, user.TelegramID, user.Name, user.Username, user.IsPremium).Scan(&user.ID, &user.CreatedAt)
}

// GetUserByID retrieves a user by their ID
func (db *DB) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserByTgID(tg_id int) (*models.User, error) {
	user := &models.User{}
	err := db.Get(user, "SELECT * FROM users WHERE tg_id = $1", tg_id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) GetUserIdByTgId(tg_id int) (int, error) {
	var id int
	err := db.QueryRowx("SELECT id FROM users WHERE tg_id = $1", tg_id).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, err
}

func (db *DB) GetTgIdByUserId(id int) (int, error) {
	var tg_id int
	err := db.QueryRowx("SELECT tg_id FROM users WHERE id = $1", id).Scan(&tg_id)
	if err != nil {
		return -1, err
	}
	return tg_id, err
}

// UpdateUser updates an existing user
func (db *DB) UpdateUser(user *models.User) error {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	_, err := db.Exec(query, user.Username, user.ID)
	return err
}

// DeleteUser removes a user from the database
func (db *DB) DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// GetAllUsers retrieves all users from the database
func (db *DB) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := db.Select(&users, "SELECT * FROM users")

	return users, err
}


func (db *DB) UserExists(user  *models.User) (bool, error) {
	var exist bool
	err := db.QueryRowx("SELECT EXISTS(SELECT 1 FROM users WHERE tg_id = $1);", user.TelegramID).Scan(&exist)
	return exist, err
}

// AddMessage adds a new message to the message_queue
func (db *DB) AddMessage(message *models.Message) error {
	query := `INSERT INTO message_queue (user_id, question, status) VALUES ($1, $2, $3) RETURNING id, created_at`
	return db.QueryRowx(query, message.UserID, message.Question, message.Status).Scan(&message.ID, &message.CreatedAt)
}

// GetMessageByID retrieves a message by its ID
func (db *DB) GetMessageByID(id int) (*models.Message, error) {
	message := &models.Message{}
	err := db.Get(message, "SELECT * FROM message_queue WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// GetMessageByStatus retrieves a message by its status
func (db *DB) GetMessageByStatus(status string) ([]models.Message, error) {
	var messages []models.Message
	err := db.Select(&messages, "SELECT * FROM message_queue WHERE status = $1", status)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// UpdateMessageStatus updates the status of a message
func (db *DB) UpdateMessageStatus(id int, status string) error {
	query := `UPDATE message_queue SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := db.Exec(query, status, id)
	return err
}

// DeleteMessage removes a message from the message_queue
func (db *DB) DeleteMessage(id int) error {
	_, err := db.Exec("DELETE FROM message_queue WHERE id = $1", id)
	return err
}

// GetAllMessages retrieves all messages from the message_queue
func (db *DB) GetAllMessages() ([]models.Message, error) {
	var messages []models.Message
	err := db.Select(&messages, "SELECT * FROM message_queue")
	return messages, err
}

// UpdateMessage updates an existing message
func (db *DB) UpdateMessage(message *models.Message) error {
	query := `UPDATE message_queue SET user_id = $1, question = $2, answer = $3, status = $4 WHERE id = $5`
	_, err := db.Exec(query, message.UserID, message.Question, message.Answer, message.Status, message.ID)
	return err
}
