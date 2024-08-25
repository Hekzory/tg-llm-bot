package models

import (
	"time"
)

type Conversation struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
