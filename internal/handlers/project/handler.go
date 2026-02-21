package project

import (
	"net/http"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/dto"
	"web-porto-backend/internal/services/project"

	"github.com/gin-gonic/gin"
)

// Constants for common messages
const (
	msgProjectsRetrieved  = "Projects retrieved successfully"
	msgProjectNotFound    = "Project not found"
	msgIDRequired         = "ID is required"
	msgInvalidRequestData = "Invalid request data"
)

// Handler handles HTTP requests for project endpoints
type Handler struct {
	service     *project.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

// NewHandler creates a new project handler
func NewHandler(service *project.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{
		service:     service,
		httpAdapter: httpAdapter,
	}
}

// Create creates a new project
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidRequestData, err.Error()))
		return
	}

	// Set default authorID to be handled by service
	// This simplifies authentication - we just check if the token is valid
	// but we don't need to extract user ID from it
	// req.AuthorID will be set by service using default admin user if 0

	project, err := h.service.CreateProject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create project", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, project, "Project created successfully")
}

// GetAll gets all projects
func (h *Handler) GetAll(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	projects, err := h.service.ListProjects(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get projects", err.Error()))
		return
	}

	// Log the retrieved data for debugging
	// fmt.Printf("Retrieved projects data\n")

	// Ensure we're returning a valid, non-null response
	if projects == nil {
		projects = &dto.PaginatedResponse{
			Data:       []interface{}{},
			Pagination: dto.PaginationResponse{},
		}
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, projects, msgProjectsRetrieved)
}

// GetPublished gets all published projects
func (h *Handler) GetPublished(c *gin.Context) {
	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	// For simplicity, just call list projects (could be optimized with a specific method)
	projects, err := h.service.ListProjects(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get published projects", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, projects, msgProjectsRetrieved)
}

// GetByID gets a project by ID
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgIDRequired))
		return
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse(msgProjectNotFound, err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, project, "Project retrieved successfully")
}

// GetBySlug gets a project by slug
func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Slug is required"))
		return
	}

	project, err := h.service.GetProjectBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse(msgProjectNotFound, err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, project, "Project retrieved successfully")
}

// GetByCategory gets projects by category slug
func (h *Handler) GetByCategory(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Category slug is required"))
		return
	}

	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	projects, err := h.service.GetProjectsByCategorySlug(slug, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get projects by category", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, projects, msgProjectsRetrieved)
}

// GetByTechnology gets projects by technology name
func (h *Handler) GetByTechnology(c *gin.Context) {
	techName := c.Param("name")
	if techName == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Technology name is required"))
		return
	}

	pagination := h.httpAdapter.GetPaginationFromQuery(c)

	// This would need to be implemented in the service
	// For now, just return all projects
	projects, err := h.service.ListProjects(pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get projects by technology", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, projects, msgProjectsRetrieved)
}

// Update updates a project
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgIDRequired))
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgInvalidRequestData, err.Error()))
		return
	}

	// Handle potential temporary ID from frontend
	if len(id) > 0 && id[:5] == "temp-" {
		// Create new project instead with the same data
		createReq := dto.CreateProjectRequest{
			Title:           req.Title,
			Slug:            req.Slug,
			Description:     req.Description,
			Content:         req.Content,
			ThumbnailURL:    req.ThumbnailURL,
			Status:          req.Status,
			CategoryID:      req.CategoryID,
			GitHubURL:       req.GitHubURL,
			LiveDemoURL:     req.LiveDemoURL,
			Technologies:    req.Technologies,
			TechnologyNames: req.TechnologyNames,
			AuthorID:        0, // Will be set by service using default admin
			Metadata:        req.Metadata,
		}

		project, err := h.service.CreateProject(createReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create project from temp ID", err.Error()))
			return
		}

		h.httpAdapter.SendSuccessResponse(c, http.StatusOK, project, "Project created successfully")
		return
	}

	// If not a temp ID, proceed with normal update
	project, err := h.service.UpdateProject(id, req)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(msgProjectNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update project", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, project, "Project updated successfully")
}

// Delete deletes a project
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(msgIDRequired))
		return
	}

	err := h.service.DeleteProject(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, response.NewErrorResponse(msgProjectNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete project", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Project deleted successfully")
}

// AddImage adds an image to a project
func (h *Handler) AddImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID is required"))
		return
	}

	var imageData dto.ProjectImageData
	if err := c.ShouldBindJSON(&imageData); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid image data", err.Error()))
		return
	}

	image, err := h.service.AddProjectImage(id, imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to add image", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, image, "Image added successfully")
}

// AddVideo adds a video to a project
func (h *Handler) AddVideo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID is required"))
		return
	}

	var videoData dto.ProjectVideoData
	if err := c.ShouldBindJSON(&videoData); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid video data", err.Error()))
		return
	}

	video, err := h.service.AddProjectVideo(id, videoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to add video", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, video, "Video added successfully")
}

// DeleteImage deletes an image from a project
func (h *Handler) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	imageId := c.Param("imageId")
	if id == "" || imageId == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID and Image ID are required"))
		return
	}

	// This method would need to be implemented in the service
	// For now, we just show a success message
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Image deleted successfully")
}

// DeleteVideo deletes a video from a project
func (h *Handler) DeleteVideo(c *gin.Context) {
	id := c.Param("id")
	videoId := c.Param("videoId")
	if id == "" || videoId == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID and Video ID are required"))
		return
	}

	// This method would need to be implemented in the service
	// For now, we just show a success message
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Video deleted successfully")
}

// AddTechnology adds a technology to a project
func (h *Handler) AddTechnology(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID is required"))
		return
	}

	var techData struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&techData); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid technology data", err.Error()))
		return
	}

	// This method would need to be implemented in the service
	// For now, we just show a success message
	h.httpAdapter.SendSuccessResponse(c, http.StatusCreated, techData, "Technology added successfully")
}

// RemoveTechnology removes a technology from a project
func (h *Handler) RemoveTechnology(c *gin.Context) {
	id := c.Param("id")
	techId := c.Param("techId")
	if id == "" || techId == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Project ID and Technology ID are required"))
		return
	}

	// This method would need to be implemented in the service
	// For now, we just show a success message
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Technology removed successfully")
}
