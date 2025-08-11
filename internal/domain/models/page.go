package models

import "time"

type Page struct {
	ID        int    `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Slug      string `gorm:"unique;not null"`
	Content   string `gorm:"not null"`
	Status    string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
