package models

import "time"

type PageView struct {
	ID        int    `gorm:"primaryKey"`
	Page      string `gorm:"index;not null"`
	VisitorID string `gorm:"index;not null"`
	UserAgent string
	Referrer  string
	IP        string
	Country   string
	City      string
	Timestamp time.Time `gorm:"index;autoCreateTime"`
}
