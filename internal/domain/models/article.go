package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Article represents a blog article
type Article struct {
	ID               string `gorm:"primaryKey;type:uuid"`
	Title            string `gorm:"not null"`
	Slug             string `gorm:"unique;not null"`
	Excerpt          string `gorm:"type:text"`
	Content          string `gorm:"type:text;not null"`
	FeaturedImageURL string
	Status           string `gorm:"not null;default:'draft'"`
	AuthorID         int
	Author           User `gorm:"foreignKey:AuthorID"`
	PublishedAt      *time.Time
	ReadTime         int            `gorm:"default:0"`
	ViewCount        int            `gorm:"default:0"`
	Metadata         string         `gorm:"type:jsonb"`
	Categories       []Category     `gorm:"many2many:article_categories;"`
	Tags             []Tag          `gorm:"many2many:article_tags;"`
	Images           []ArticleImage `gorm:"foreignKey:ArticleID"`
	Videos           []ArticleVideo `gorm:"foreignKey:ArticleID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ArticleImage for multiple images in an article
type ArticleImage struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	ArticleID string `gorm:"type:uuid"`
	URL       string `gorm:"not null"`
	Caption   string
	AltText   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ArticleVideo for videos in an article
type ArticleVideo struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	ArticleID string `gorm:"type:uuid"`
	URL       string `gorm:"not null"`
	Caption   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to generate UUID
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for ArticleImage
func (ai *ArticleImage) BeforeCreate(tx *gorm.DB) error {
	if ai.ID == "" {
		ai.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for ArticleVideo
func (av *ArticleVideo) BeforeCreate(tx *gorm.DB) error {
	if av.ID == "" {
		av.ID = uuid.New().String()
	}
	return nil
}
