package models

import "time"

type Log struct {
	ID          int `gorm:"primaryKey"`
	UserID      *int
	Action      string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
}
