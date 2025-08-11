package models

import "time"

type User struct {
	ID           int    `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
