package models

import "time"

type User struct {
	ID         int       `db:"id"`
	TelegramID int       `db:"tg_id"`
	Name       string    `db:"name"`
	Username   string    `db:"username"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IsPremium  bool      `db:"is_premium"`
}
