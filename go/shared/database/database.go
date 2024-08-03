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
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, created_at`
	return db.QueryRowx(query, user.Username).Scan(&user.ID, &user.CreatedAt)
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
