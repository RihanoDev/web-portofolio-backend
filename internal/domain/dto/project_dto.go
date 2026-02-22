package dto

import (
	"time"
)

// Project DTOs
type CreateProjectRequest struct {
	Title           string                 `json:"title" validate:"required"`
	Slug            string                 `json:"slug" validate:"required"`
	Description     string                 `json:"description"`
	Content         string                 `json:"content" validate:"required"`
	ThumbnailURL    string                 `json:"thumbnailUrl"`
	Status          string                 `json:"status" validate:"required,oneof=draft published private"`
	CategoryID      *int                   `json:"categoryId"`
	CategoryIds     []int                  `json:"categoryIds"`
	CategoryIdStrs  []string               `json:"categoryIdStrs"`
	AuthorID        int                    `json:"authorId"` // Not required anymore
	GitHubURL       string                 `json:"githubUrl"`
	LiveDemoURL     string                 `json:"liveDemoUrl"`
	Images          []ProjectImageData     `json:"images"`
	Videos          []ProjectVideoData     `json:"videos"`
	Categories      []int                  `json:"categories"`      // Category IDs
	CategoryNames   []string               `json:"categoryNames"`   // Category names for auto-creation
	Technologies    []int                  `json:"technologies"`    // Tag IDs (preferred)
	TechnologyNames []string               `json:"technologyNames"` // Alternative: Tag names
	Tags            []int                  `json:"tags"`            // General tags IDs
	TagIds          []int                  `json:"tagIds"`
	TagIdStrs       []string               `json:"tagIdStrs"`
	TagNames        []string               `json:"tagNames"` // General tags names
	Metadata        map[string]interface{} `json:"metadata"`
}

type UpdateProjectRequest struct {
	Title           *string                `json:"title"`
	Slug            *string                `json:"slug"`
	Description     *string                `json:"description"`
	Content         *string                `json:"content"`
	ThumbnailURL    *string                `json:"thumbnailUrl"`
	Status          *string                `json:"status" validate:"omitempty,oneof=draft published private"`
	CategoryID      *int                   `json:"categoryId"`
	CategoryIds     []int                  `json:"categoryIds"`
	CategoryIdStrs  []string               `json:"categoryIdStrs"`
	Categories      []int                  `json:"categories"`
	GitHubURL       *string                `json:"githubUrl"`
	LiveDemoURL     *string                `json:"liveDemoUrl"`
	Technologies    []int                  `json:"technologies"`
	TechnologyNames []string               `json:"technologyNames"`
	Tags            []int                  `json:"tags"`
	TagIds          []int                  `json:"tagIds"`
	TagIdStrs       []string               `json:"tagIdStrs"`
	TagNames        []string               `json:"tagNames"`
	Images          []ProjectImageData     `json:"images"`
	Videos          []ProjectVideoData     `json:"videos"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type ProjectImageData struct {
	URL       string `json:"url" validate:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ProjectVideoData struct {
	URL       string `json:"url" validate:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ProjectResponse struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Slug         string                 `json:"slug"`
	Description  string                 `json:"description"`
	Content      string                 `json:"content"`
	ThumbnailURL string                 `json:"thumbnailUrl"`
	Status       string                 `json:"status"`
	Category     *CategoryResponse      `json:"category,omitempty"`
	Categories   []CategoryResponse     `json:"categories"`
	Author       AuthorResponse         `json:"author"`
	GitHubURL    string                 `json:"githubUrl"`
	LiveDemoURL  string                 `json:"liveDemoUrl"`
	Images       []ProjectImageResponse `json:"images"`
	Videos       []ProjectVideoResponse `json:"videos"`
	Technologies []TagResponse          `json:"technologies"`
	Tags         []TagResponse          `json:"tags"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type ProjectImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ProjectVideoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ProjectListResponse struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	Slug           string                 `json:"slug"`
	Description    string                 `json:"description"`
	ThumbnailURL   string                 `json:"thumbnailUrl"`
	Status         string                 `json:"status"`
	CategoryID     *int                   `json:"categoryId,omitempty"`
	Category       string                 `json:"category,omitempty"`
	Categories     []string               `json:"categories"`
	CategoryModels []CategoryResponse     `json:"categoryModels,omitempty"`
	AuthorName     string                 `json:"authorName"`
	Tags           []TagResponse          `json:"tags,omitempty"`
	Technologies   []TagResponse          `json:"technologies,omitempty"`
	TechnologyStrs []string               `json:"technologyStrs,omitempty"`
	TagStrs        []string               `json:"tagStrs,omitempty"`
	GitHubURL      string                 `json:"githubUrl"`
	LiveDemoURL    string                 `json:"liveDemoUrl"`
	Content        string                 `json:"content"`
	Metadata       map[string]interface{} `json:"metadata"`
	Images         []ProjectImageResponse `json:"images"`
	Videos         []ProjectVideoResponse `json:"videos"`
	CreatedAt      time.Time              `json:"createdAt"`
}
