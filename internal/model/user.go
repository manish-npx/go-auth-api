package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // hidden in JSON
	CreatedAt    time.Time `json:"created_at"`
}
