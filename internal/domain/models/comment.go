package models

import "time"

type Comment struct {
	ID        int `gorm:"primaryKey"`
	PostID    int `gorm:"index"`
	UserID    *int
	Content   string `gorm:"not null"`
	ParentID  *int
	CreatedAt time.Time
}
