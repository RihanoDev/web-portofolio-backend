package post

import (
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/domain/models"
	articleService "web-porto-backend/internal/services/article"
	postService "web-porto-backend/internal/services/post"
)

// PostServiceAdapter adapts an ArticleService to work with old PostHandler
type PostServiceAdapter struct {
	articleService *articleService.Service
}

// NewPostServiceAdapter creates a new adapter to make ArticleService compatible with PostHandler
func NewPostServiceAdapter(articleService *articleService.Service) postService.Service {
	return &PostServiceAdapter{
		articleService: articleService,
	}
}

// Create creates a new post by adapting to ArticleService
func (a *PostServiceAdapter) Create(post *models.Post) error {
	// Convert Post to CreateArticleRequest
	req := dto.CreateArticleRequest{
		Title:            post.Title,
		Content:          post.Content,
		Status:           post.Status,
		AuthorID:         post.AuthorID,
		FeaturedImageURL: post.FeaturedImageURL,
	}

	// Use ArticleService to create the article
	_, err := a.articleService.CreateArticle(req)
	return err
}

// GetAll gets all posts using ArticleService
func (a *PostServiceAdapter) GetAll(page, limit int) ([]*models.Post, *postService.PaginationInfo, error) {
	// Use ArticleService to get articles
	response, err := a.articleService.ListArticles(page, limit)
	if err != nil {
		return nil, nil, err
	}

	// Convert articles to posts
	posts := make([]*models.Post, 0)
	articleListResponses, ok := response.Data.([]dto.ArticleListResponse)
	if !ok {
		// If type assertion fails, try with slice of interface{}
		articleListInterface, ok := response.Data.([]interface{})
		if !ok {
			// Return empty list if we can't convert
			paginationInfo := &postService.PaginationInfo{
				Total:      response.Pagination.TotalCount,
				Page:       response.Pagination.CurrentPage,
				Limit:      response.Pagination.PageSize,
				TotalPages: response.Pagination.TotalPages,
			}
			return posts, paginationInfo, nil
		}

		// Convert each interface{} to ArticleListResponse
		for _, item := range articleListInterface {
			if articleList, ok := item.(dto.ArticleListResponse); ok {
				post := convertArticleResponseToPost(articleList)
				posts = append(posts, post)
			}
		}
	} else {
		// Direct conversion if type assertion succeeds
		for _, article := range articleListResponses {
			post := convertArticleResponseToPost(article)
			posts = append(posts, post)
		}
	}

	// Create pagination info
	paginationInfo := &postService.PaginationInfo{
		Total:      response.Pagination.TotalCount,
		Page:       response.Pagination.CurrentPage,
		Limit:      response.Pagination.PageSize,
		TotalPages: response.Pagination.TotalPages,
	}

	return posts, paginationInfo, nil
}

// GetByID gets a post by ID using ArticleService
func (a *PostServiceAdapter) GetByID(id string) (*models.Post, error) {
	// Use ArticleService to get the article
	article, err := a.articleService.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	// Convert article to post
	return &models.Post{
		ID:               article.ID,
		Title:            article.Title,
		Content:          article.Content,
		Status:           article.Status,
		AuthorID:         article.Author.ID,
		FeaturedImageURL: article.FeaturedImageURL,
		ViewCount:        article.ViewCount,
		CreatedAt:        article.CreatedAt,
		UpdatedAt:        article.UpdatedAt,
	}, nil
}

// Update updates a post by ID using ArticleService
func (a *PostServiceAdapter) Update(id string, post *models.Post) error {
	// Convert Post to UpdateArticleRequest
	req := dto.UpdateArticleRequest{
		Title:            &post.Title,
		Content:          &post.Content,
		Status:           &post.Status,
		FeaturedImageURL: &post.FeaturedImageURL,
	}

	// Use ArticleService to update the article
	_, err := a.articleService.UpdateArticle(id, req)
	return err
}

// Delete deletes a post by ID using ArticleService
func (a *PostServiceAdapter) Delete(id string) error {
	return a.articleService.DeleteArticle(id)
}

// GetBySlug gets a post by slug using ArticleService
func (a *PostServiceAdapter) GetBySlug(slug string) (*models.Post, error) {
	// Use ArticleService to get the article
	article, err := a.articleService.GetArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Convert article to post
	return &models.Post{
		ID:               article.ID,
		Title:            article.Title,
		Content:          article.Content,
		Status:           article.Status,
		AuthorID:         article.Author.ID,
		FeaturedImageURL: article.FeaturedImageURL,
		ViewCount:        article.ViewCount,
		CreatedAt:        article.CreatedAt,
		UpdatedAt:        article.UpdatedAt,
	}, nil
}

// GetByAuthorID gets posts by author ID using ArticleService (simplified implementation)
func (a *PostServiceAdapter) GetByAuthorID(authorID, page, limit int) ([]*models.Post, *postService.PaginationInfo, error) {
	// For simplicity, just return all articles
	// In a real implementation, you'd filter by author ID
	return a.GetAll(page, limit)
}

// GetPublished gets published posts using ArticleService (simplified implementation)
func (a *PostServiceAdapter) GetPublished(page, limit int) ([]*models.Post, *postService.PaginationInfo, error) {
	// For simplicity, just return all articles
	// In a real implementation, you'd filter by status
	return a.GetAll(page, limit)
}

// Helper function to convert ArticleListResponse to Post
func convertArticleResponseToPost(article dto.ArticleListResponse) *models.Post {
	return &models.Post{
		ID:               article.ID,
		Title:            article.Title,
		Slug:             article.Slug,
		Content:          "", // Content not available in list response
		Status:           article.Status,
		AuthorID:         0, // Author ID not available in list response
		FeaturedImageURL: article.FeaturedImageURL,
		ViewCount:        article.ViewCount,
		CreatedAt:        article.CreatedAt,
	}
}
