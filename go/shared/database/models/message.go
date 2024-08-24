package models

import (
	"database/sql"
	"time"
)

// Message from message_queue table
type Message struct {
	ID             int           `db:"id"`
	UserID         int           `db:"user_id"`
	ConversationID sql.NullInt64 `db:"conversation_id"`
	Question       string        `db:"question"`
	Answer         string        `db:"answer"`
	Status         string        `db:"status"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
	TgQuestionId   sql.NullInt64 `db:"tg_question_id"`
	TgAnswerId     sql.NullInt64 `db:"tg_answer_id"`
}
