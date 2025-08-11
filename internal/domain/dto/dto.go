package dto

// PostDTO represents the data transfer object for posts
type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required" validate:"required,min=3,max=255"`
	Content    string   `json:"content" binding:"required" validate:"required,min=10"`
	Excerpt    string   `json:"excerpt" validate:"max=500"`
	Status     string   `json:"status" binding:"required" validate:"required,oneof=draft published"`
	CategoryID int      `json:"category_id" validate:"required,gt=0"`
	Tags       []string `json:"tags"`
	AuthorID   int      `json:"author_id" binding:"required" validate:"required,gt=0"`
}

type UpdatePostRequest struct {
	Title      string   `json:"title" validate:"min=3,max=255"`
	Content    string   `json:"content" validate:"min=10"`
	Excerpt    string   `json:"excerpt" validate:"max=500"`
	Status     string   `json:"status" validate:"oneof=draft published"`
	CategoryID int      `json:"category_id" validate:"gt=0"`
	Tags       []string `json:"tags"`
}

// PageDTO represents the data transfer object for pages
type CreatePageRequest struct {
	Title   string `json:"title" binding:"required" validate:"required,min=3,max=255"`
	Content string `json:"content" binding:"required" validate:"required,min=10"`
	Status  string `json:"status" binding:"required" validate:"required,oneof=draft published"`
}

type UpdatePageRequest struct {
	Title   string `json:"title" validate:"min=3,max=255"`
	Content string `json:"content" validate:"min=10"`
	Status  string `json:"status" validate:"oneof=draft published"`
}

// CommentDTO represents the data transfer object for comments
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required" validate:"required,min=3,max=1000"`
	PostID   int    `json:"post_id" binding:"required" validate:"required,gt=0"`
	AuthorID int    `json:"author_id" binding:"required" validate:"required,gt=0"`
	ParentID *int   `json:"parent_id,omitempty" validate:"omitempty,gt=0"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"min=3,max=1000"`
}

// CategoryDTO represents the data transfer object for categories
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

// UserDTO represents the data transfer object for users
type CreateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=admin editor user"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"email"`
	Role     string `json:"role" validate:"oneof=admin editor user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=6"`
}

// AuthDTO represents the data transfer object for authentication
type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=admin editor user"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Page  int `form:"page" validate:"min=1"`
	Limit int `form:"limit" validate:"min=1,max=100"`
}

// SearchQuery represents search query parameters
type SearchQuery struct {
	Query    string `form:"q"`
	Status   string `form:"status" validate:"omitempty,oneof=draft published"`
	Category string `form:"category"`
	Author   string `form:"author"`
	Tag      string `form:"tag"`
	SortBy   string `form:"sort_by" validate:"omitempty,oneof=created_at updated_at title"`
	SortDir  string `form:"sort_dir" validate:"omitempty,oneof=asc desc"`
}
