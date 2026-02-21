package models

type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"type:text"`
}
