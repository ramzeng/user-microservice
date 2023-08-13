package model

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email" gorm:"uniqueIndex;size:50"`
	Password  string    `json:"password" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}
