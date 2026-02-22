package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID           string `gorm:"primaryKey;type:uuid"`
	Title        string `gorm:"not null"`
	Slug         string `gorm:"unique;not null"`
	Description  string `gorm:"type:text"`
	Content      string `gorm:"type:text"`
	ThumbnailURL string
	Status       string `gorm:"not null;default:'published'"`
	CategoryID   *int
	Category     *Category  `gorm:"foreignKey:CategoryID"`
	Categories   []Category `gorm:"many2many:project_categories;"`
	AuthorID     int
	Author       User           `gorm:"foreignKey:AuthorID"`
	Technologies []Tag          `gorm:"many2many:project_technologies;"`
	Tags         []Tag          `gorm:"many2many:project_tags;"`
	Metadata     string         `gorm:"type:jsonb;default:'{}'"`
	GitHubURL    string         `gorm:"column:github_url"`
	LiveDemoURL  string         `gorm:"column:live_demo_url"`
	Images       []ProjectImage `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE;"`
	Videos       []ProjectVideo `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Images for the project
type ProjectImage struct {
	ID        string  `gorm:"primaryKey;type:uuid"`
	ProjectID string  `gorm:"type:uuid"`
	Project   Project `gorm:"foreignKey:ProjectID"`
	URL       string  `gorm:"not null"`
	Caption   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Videos for the project (e.g. YouTube embeds)
type ProjectVideo struct {
	ID        string  `gorm:"primaryKey;type:uuid"`
	ProjectID string  `gorm:"type:uuid"`
	Project   Project `gorm:"foreignKey:ProjectID"`
	URL       string  `gorm:"not null"` // Can be YouTube/Vimeo URL
	Caption   string
	SortOrder int `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to generate UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for ProjectImage
func (pi *ProjectImage) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == "" {
		pi.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate hook for ProjectVideo
func (pv *ProjectVideo) BeforeCreate(tx *gorm.DB) error {
	if pv.ID == "" {
		pv.ID = uuid.New().String()
	}
	return nil
}
