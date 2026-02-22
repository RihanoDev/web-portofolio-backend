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
	Images           []ArticleImage `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"`
	Videos           []ArticleVideo `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BeforeCreate hook to generate UUID
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

type ArticleImage struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ArticleID string `gorm:"type:uuid;index"`
	URL       string `gorm:"not null"`
	Caption   string
	AltText   string
	SortOrder int `gorm:"default:0"`
}

type ArticleVideo struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ArticleID string `gorm:"type:uuid;index"`
	URL       string `gorm:"not null"`
	Caption   string
	SortOrder int `gorm:"default:0"`
}
