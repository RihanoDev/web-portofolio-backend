package services

import (
	"errors"
	"time"
	"web-porto-backend/common/helper"
	"web-porto-backend/common/utils"
	"web-porto-backend/internal/domain/models"
)

// PostDomainService encapsulates business rules for posts
type PostDomainService struct {
	helper *helper.Helper
}

// NewPostDomainService creates a new post domain service
func NewPostDomainService() *PostDomainService {
	return &PostDomainService{
		helper: helper.NewHelper(),
	}
}

// PreparePostForCreation prepares a post for creation with business rules
func (s *PostDomainService) PreparePostForCreation(post *models.Post, title, content string) error {
	if err := s.validatePostData(title, content); err != nil {
		return err
	}

	// Generate slug if not provided
	if post.Slug == "" {
		post.Slug = utils.StringToSlug(title)
	}

	// Set timestamps
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	// Set published date if status is published
	if post.Status == "published" && post.PublishedAt == nil {
		post.PublishedAt = &now
	}

	return nil
}

// PreparePostForUpdate prepares a post for update with business rules
func (s *PostDomainService) PreparePostForUpdate(post *models.Post, updates map[string]interface{}) error {
	// Update timestamp
	post.UpdatedAt = time.Now()

	// Handle status change to published
	if status, exists := updates["status"]; exists && status == "published" {
		if post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	// Regenerate slug if title changed
	if title, exists := updates["title"]; exists {
		if titleStr, ok := title.(string); ok && post.Slug == utils.StringToSlug(post.Title) {
			// Only regenerate if current slug is auto-generated
			post.Slug = utils.StringToSlug(titleStr)
		}
	}

	return nil
}

// ValidatePostForPublication validates if a post can be published
func (s *PostDomainService) ValidatePostForPublication(post *models.Post) error {
	if post.Title == "" {
		return errors.New("post title is required for publication")
	}

	if post.Content == "" {
		return errors.New("post content is required for publication")
	}

	if post.AuthorID == 0 {
		return errors.New("post author is required for publication")
	}

	return nil
}

// validatePostData validates basic post data
func (s *PostDomainService) validatePostData(title, content string) error {
	if title == "" {
		return errors.New("post title is required")
	}

	if len(title) > 255 {
		return errors.New("post title cannot exceed 255 characters")
	}

	if content == "" {
		return errors.New("post content is required")
	}

	if len(content) < 10 {
		return errors.New("post content must be at least 10 characters")
	}

	return nil
}

// PageDomainService encapsulates business rules for pages
type PageDomainService struct {
	helper *helper.Helper
}

// NewPageDomainService creates a new page domain service
func NewPageDomainService() *PageDomainService {
	return &PageDomainService{
		helper: helper.NewHelper(),
	}
}

// PreparePageForCreation prepares a page for creation with business rules
func (s *PageDomainService) PreparePageForCreation(page *models.Page, title, content string) error {
	if err := s.validatePageData(title, content); err != nil {
		return err
	}

	// Generate slug if not provided
	if page.Slug == "" {
		page.Slug = utils.StringToSlug(title)
	}

	// Set timestamps
	now := time.Now()
	page.CreatedAt = now
	page.UpdatedAt = now

	return nil
}

// PreparePageForUpdate prepares a page for update with business rules
func (s *PageDomainService) PreparePageForUpdate(page *models.Page, updates map[string]interface{}) error {
	// Update timestamp
	page.UpdatedAt = time.Now()

	// Regenerate slug if title changed
	if title, exists := updates["title"]; exists {
		if titleStr, ok := title.(string); ok && page.Slug == utils.StringToSlug(page.Title) {
			// Only regenerate if current slug is auto-generated
			page.Slug = utils.StringToSlug(titleStr)
		}
	}

	return nil
}

// validatePageData validates basic page data
func (s *PageDomainService) validatePageData(title, content string) error {
	if title == "" {
		return errors.New("page title is required")
	}

	if len(title) > 255 {
		return errors.New("page title cannot exceed 255 characters")
	}

	if content == "" {
		return errors.New("page content is required")
	}

	if len(content) < 10 {
		return errors.New("page content must be at least 10 characters")
	}

	return nil
}

// CommentDomainService encapsulates business rules for comments
type CommentDomainService struct{}

// NewCommentDomainService creates a new comment domain service
func NewCommentDomainService() *CommentDomainService {
	return &CommentDomainService{}
}

// PrepareCommentForCreation prepares a comment for creation with business rules
func (s *CommentDomainService) PrepareCommentForCreation(comment *models.Comment) error {
	if err := s.validateCommentData(comment.Content); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	comment.CreatedAt = now

	return nil
}

// ValidateCommentHierarchy validates comment parent-child relationship
func (s *CommentDomainService) ValidateCommentHierarchy(comment *models.Comment, parentComment *models.Comment) error {
	if comment.ParentID == nil {
		return nil // Top-level comment, no validation needed
	}

	if parentComment == nil {
		return errors.New("parent comment not found")
	}

	if parentComment.PostID != comment.PostID {
		return errors.New("parent comment must belong to the same post")
	}

	// Prevent deeply nested comments (max 3 levels)
	if parentComment.ParentID != nil {
		return errors.New("maximum comment nesting level reached")
	}

	return nil
}

// validateCommentData validates basic comment data
func (s *CommentDomainService) validateCommentData(content string) error {
	if content == "" {
		return errors.New("comment content is required")
	}

	if len(content) < 3 {
		return errors.New("comment content must be at least 3 characters")
	}

	if len(content) > 1000 {
		return errors.New("comment content cannot exceed 1000 characters")
	}

	return nil
}

// UserDomainService encapsulates business rules for users
type UserDomainService struct{}

// NewUserDomainService creates a new user domain service
func NewUserDomainService() *UserDomainService {
	return &UserDomainService{}
}

// PrepareUserForCreation prepares a user for creation with business rules
func (s *UserDomainService) PrepareUserForCreation(user *models.User) error {
	// Set default role if not specified
	if user.Role == "" {
		user.Role = "user"
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// ValidateUserRole validates user role permissions
func (s *UserDomainService) ValidateUserRole(requesterRole, targetRole string) error {
	roleHierarchy := map[string]int{
		"user":   1,
		"editor": 2,
		"admin":  3,
	}

	requesterLevel, exists := roleHierarchy[requesterRole]
	if !exists {
		return errors.New("invalid requester role")
	}

	targetLevel, exists := roleHierarchy[targetRole]
	if !exists {
		return errors.New("invalid target role")
	}

	if requesterLevel <= targetLevel {
		return errors.New("insufficient permissions to assign this role")
	}

	return nil
}
