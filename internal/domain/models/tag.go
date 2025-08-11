package models

type Tag struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}
