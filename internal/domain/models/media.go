package models

import "time"

// Media represents an uploaded file stored on the server
type Media struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FileName     string    `gorm:"not null" json:"fileName"`
	OriginalName string    `gorm:"not null" json:"originalName"`
	FilePath     string    `gorm:"not null" json:"filePath"`
	FileURL      string    `gorm:"not null" json:"fileUrl"`
	FileType     string    `gorm:"not null" json:"fileType"`
	FileSize     int64     `gorm:"not null;default:0" json:"fileSize"`
	MimeType     string    `gorm:"not null" json:"mimeType"`
	UploadedBy   *uint     `json:"uploadedBy,omitempty"`
	UploadedAt   time.Time `gorm:"autoCreateTime" json:"uploadedAt"`
}
