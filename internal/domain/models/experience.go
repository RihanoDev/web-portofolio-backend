package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringArray is a custom type for handling PostgreSQL arrays
type StringArray []string

// Value implements the driver.Valuer interface
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return nil, nil
	}
	return json.Marshal(sa)
}

// Scan implements the sql.Scanner interface
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, sa)
	case string:
		return json.Unmarshal([]byte(v), sa)
	}

	return nil
}

// Experience represents a work experience entry
type Experience struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	Title            string `gorm:"not null"`
	Company          string `gorm:"not null"`
	Location         string
	StartDate        time.Time `gorm:"not null"`
	EndDate          *time.Time
	Current          bool              `gorm:"default:false"`
	Description      string            `gorm:"type:text"`
	Responsibilities StringArray       `gorm:"type:jsonb"`                         // Keep as JSON for flexibility
	Technologies     []Tag             `gorm:"many2many:experience_technologies;"` // Use relational for indexing
	Images           []ExperienceImage `gorm:"foreignKey:ExperienceID;constraint:OnDelete:CASCADE;"`
	CompanyURL       string
	LogoURL          string
	Metadata         string `gorm:"type:jsonb;default:'{}'"` // Flexible data like theme, extra info
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ExperienceImage represents an image related to a work experience
type ExperienceImage struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ExperienceID int       `gorm:"not null;index"`
	URL          string    `gorm:"not null"`
	Caption      string    `gorm:"type:text"`
	SortOrder    int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
