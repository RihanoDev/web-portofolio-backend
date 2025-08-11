package models

import "time"

type Media struct {
	ID         int    `gorm:"primaryKey"`
	FileName   string `gorm:"not null"`
	FilePath   string `gorm:"not null"`
	FileType   string
	UploadedBy *int
	UploadedAt time.Time
}
