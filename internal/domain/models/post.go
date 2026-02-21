package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID               string `gorm:"primaryKey;type:uuid"`
	Title            string `gorm:"not null"`
	Slug             string `gorm:"unique;not null"`
	Excerpt          string `gorm:"type:text"`
	Content          string `gorm:"type:text;not null"` // Will store HTML content
	FeaturedImageURL string
	Status           string `gorm:"not null;default:'draft'"`
	AuthorID         int
	Author           User `gorm:"foreignKey:AuthorID"`
	PublishedAt      *time.Time
	ReadTime         int    // Estimated reading time in minutes
	ViewCount        int    `gorm:"default:0"`
	Metadata         string `gorm:"type:json"` // For storing additional metadata as JSON
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PostCategory struct {
	PostID     string   `gorm:"primaryKey;type:uuid"`
	CategoryID int      `gorm:"primaryKey"`
	Post       Post     `gorm:"foreignKey:PostID"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}

type PostTag struct {
	PostID string `gorm:"primaryKey;type:uuid"`
	TagID  int    `gorm:"primaryKey"`
	Post   Post   `gorm:"foreignKey:PostID"`
	Tag    Tag    `gorm:"foreignKey:TagID"`
}

// Additional images for the article
type PostImage struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	PostID    string `gorm:"type:uuid"`
	Post      Post   `gorm:"foreignKey:PostID"`
	URL       string `gorm:"not null"`
	Caption   string
	AltText   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Videos embedded in the article
type PostVideo struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	PostID    string `gorm:"type:uuid"`
	Post      Post   `gorm:"foreignKey:PostID"`
	URL       string `gorm:"not null"` // Can be YouTube/Vimeo URL
	Caption   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to generate UUID for Post
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for PostImage
func (pi *PostImage) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == "" {
		pi.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for PostVideo
func (pv *PostVideo) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == "" {
		pv.ID = uuid.New().String()
	}
	return nil
}
