package dto

// CreateTagRequest represents data needed to create a tag
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

// UpdateTagRequest represents data needed to update a tag
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}
