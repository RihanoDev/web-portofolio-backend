package models

type Role struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
}

type UserRole struct {
	UserID int `gorm:"primaryKey"`
	RoleID int `gorm:"primaryKey"`
}
