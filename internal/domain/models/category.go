package models

type Category struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Slug        string `gorm:"unique;not null"`
	Description string
}
