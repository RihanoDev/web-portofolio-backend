package dto

import "time"

// Experience DTOs
type CreateExperienceRequest struct {
	Title            string                 `json:"title" validate:"required"`
	Company          string                 `json:"company" validate:"required"`
	Location         string                 `json:"location"`
	StartDate        string                 `json:"startDate" validate:"required"` // Format: "YYYY-MM-DD"
	EndDate          string                 `json:"endDate"`                       // Format: "YYYY-MM-DD"
	Current          bool                   `json:"current"`
	Description      string                 `json:"description"`
	Responsibilities []string               `json:"responsibilities"` // Keep as JSON for flexibility
	TechnologyIDs    []int                  `json:"technologyIds"`    // Tag IDs (preferred)
	TechnologyNames  []string               `json:"technologyNames"`  // Alternative: Tag names
	CompanyURL       string                 `json:"companyUrl"`
	LogoURL          string                 `json:"logoUrl"`
	Images           []ExperienceImageData  `json:"images"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type UpdateExperienceRequest struct {
	Title            *string                `json:"title"`
	Company          *string                `json:"company"`
	Location         *string                `json:"location"`
	StartDate        *string                `json:"startDate"` // Format: "YYYY-MM-DD"
	EndDate          *string                `json:"endDate"`   // Format: "YYYY-MM-DD"
	Current          *bool                  `json:"current"`
	Description      *string                `json:"description"`
	Responsibilities *[]string              `json:"responsibilities"` // Keep as JSON for flexibility
	TechnologyIDs    []int                  `json:"technologyIds"`    // Tag IDs (preferred)
	TechnologyNames  []string               `json:"technologyNames"`  // Alternative: Tag names
	CompanyURL       *string                `json:"companyUrl"`
	LogoURL          *string                `json:"logoUrl"`
	Images           []ExperienceImageData  `json:"images"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type ExperienceResponse struct {
	ID               int                       `json:"id"`
	Title            string                    `json:"title"`
	Company          string                    `json:"company"`
	Location         string                    `json:"location"`
	StartDate        string                    `json:"startDate"` // Format: "YYYY-MM-DD"
	EndDate          *string                   `json:"endDate"`   // Format: "YYYY-MM-DD", nullable
	Current          bool                      `json:"current"`
	Description      string                    `json:"description"`
	Responsibilities []string                  `json:"responsibilities"` // Keep as JSON for flexibility
	Technologies     []TagResponse             `json:"technologies"`     // Return as tag objects
	CompanyURL       string                    `json:"companyUrl"`
	LogoURL          string                    `json:"logoUrl"`
	Images           []ExperienceImageResponse `json:"images"`
	Metadata         map[string]interface{}    `json:"metadata"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
}

type ExperienceImageData struct {
	URL       string `json:"url" validate:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ExperienceImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}
