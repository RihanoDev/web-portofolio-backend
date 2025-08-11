package models

type Menu struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	ParentID *int
	URL      string
	OrderNum int
}
