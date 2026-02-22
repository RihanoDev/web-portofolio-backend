package dto

import (
	"time"
)

// Article DTOs
type CreateArticleRequest struct {
	Title            string                 `json:"title" validate:"required"`
	Slug             string                 `json:"slug" validate:"required"`
	Excerpt          string                 `json:"excerpt"`
	Content          string                 `json:"content" validate:"required"`
	FeaturedImageURL string                 `json:"featuredImageUrl"`
	Status           string                 `json:"status" validate:"required,oneof=draft published private"`
	AuthorID         int                    `json:"authorId" validate:"required"`
	Categories       []int                  `json:"categories"`
	CategoryIds      []int                  `json:"categoryIds"`
	CategoryIdStrs   []string               `json:"categoryIdStrs"`
	Tags             []int                  `json:"tags"`
	TagIds           []int                  `json:"tagIds"`
	TagIdStrs        []string               `json:"tagIdStrs"`
	Images           []ArticleImageData     `json:"images"`
	Videos           []ArticleVideoData     `json:"videos"`
	PublishAt        *time.Time             `json:"publishAt"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type UpdateArticleRequest struct {
	Title            *string                `json:"title"`
	Slug             *string                `json:"slug"`
	Excerpt          *string                `json:"excerpt"`
	Content          *string                `json:"content"`
	FeaturedImageURL *string                `json:"featuredImageUrl"`
	Status           *string                `json:"status" validate:"omitempty,oneof=draft published private"`
	Categories       []int                  `json:"categories"`
	CategoryIds      []int                  `json:"categoryIds"`
	CategoryIdStrs   []string               `json:"categoryIdStrs"`
	Tags             []int                  `json:"tags"`
	TagIds           []int                  `json:"tagIds"`
	TagIdStrs        []string               `json:"tagIdStrs"`
	Images           []ArticleImageData     `json:"images"`
	Videos           []ArticleVideoData     `json:"videos"`
	PublishAt        *time.Time             `json:"publishAt"`
	Images           []ArticleImageData     `json:"images"`
	Videos           []ArticleVideoData     `json:"videos"`
	Metadata         map[string]interface{} `json:"metadata"`
}

type ArticleImageData struct {
	URL       string `json:"url" validate:"required"`
	Caption   string `json:"caption"`
	AltText   string `json:"altText"`
	SortOrder int    `json:"sortOrder"`
}

type ArticleVideoData struct {
	URL       string `json:"url" validate:"required"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ArticleResponse struct {
	ID               string                 `json:"id"`
	Title            string                 `json:"title"`
	Slug             string                 `json:"slug"`
	Excerpt          string                 `json:"excerpt"`
	Content          string                 `json:"content"`
	FeaturedImageURL string                 `json:"featuredImageUrl"`
	Status           string                 `json:"status"`
	Author           AuthorResponse         `json:"author"`
	Categories       []CategoryResponse     `json:"categories"`
	Tags             []TagResponse          `json:"tags"`
	PublishedAt      *time.Time             `json:"publishedAt"`
	ReadTime         int                    `json:"readTime"`
	ViewCount        int                    `json:"viewCount"`
	Images           []ArticleImageResponse `json:"images"`
	Videos           []ArticleVideoResponse `json:"videos"`
	Metadata         map[string]interface{} `json:"metadata"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
}

type ArticleImageResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	AltText   string `json:"altText"`
	SortOrder int    `json:"sortOrder"`
}

type ArticleVideoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sortOrder"`
}

type ArticleListResponse struct {
	ID               string                 `json:"id"`
	Title            string                 `json:"title"`
	Slug             string                 `json:"slug"`
	Excerpt          string                 `json:"excerpt"`
	FeaturedImageURL string                 `json:"featuredImageUrl"`
	Status           string                 `json:"status"`
	AuthorName       string                 `json:"authorName"`
	Categories       []string               `json:"categories"`
	CategoryModels   []CategoryResponse     `json:"categoryModels,omitempty"`
	Tags             []string               `json:"tags"`
	TagModels        []TagResponse          `json:"tagModels,omitempty"`
	PublishedAt      *time.Time             `json:"publishedAt"`
	ReadTime         int                    `json:"readTime"`
	ViewCount        int                    `json:"viewCount"`
	Content          string                 `json:"content"`
	Metadata         map[string]interface{} `json:"metadata"`
	Images           []ArticleImageResponse `json:"images"`
	Videos           []ArticleVideoResponse `json:"videos"`
	CreatedAt        time.Time              `json:"createdAt"`
}
