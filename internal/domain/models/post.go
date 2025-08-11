package models

import "time"

type Post struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Slug        string `gorm:"unique;not null"`
	Content     string `gorm:"not null"`
	Status      string `gorm:"not null"`
	AuthorID    int
	Author      User `gorm:"foreignKey:AuthorID"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostCategory struct {
	PostID     int `gorm:"primaryKey"`
	CategoryID int `gorm:"primaryKey"`
}

type PostTag struct {
	PostID int `gorm:"primaryKey"`
	TagID  int `gorm:"primaryKey"`
}
