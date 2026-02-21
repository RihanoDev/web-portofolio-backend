package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"web-porto-backend/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Handler handles media upload and management
type Handler struct {
	db        *gorm.DB
	uploadDir string
	baseURL   string
}

// NewHandler creates a new media handler
func NewHandler(db *gorm.DB, uploadDir string, baseURL string) *Handler {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	return &Handler{
		db:        db,
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// Upload handles file upload, saves to disk, and records to DB
func (h *Handler) Upload(c *gin.Context) {
	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided: " + err.Error()})
		return
	}
	defer file.Close()

	// Validate file type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		// Detect from content
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		// Reset reader
		file.Seek(0, io.SeekStart)
	}

	allowed := isAllowedMIME(mimeType)
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed. Allowed: images, PDF, video"})
		return
	}

	// Validate file size
	// Images max 10MB, Videos max 50MB
	const maxImageSize int64 = 10 * 1024 * 1024 // 10 MB
	const maxVideoSize int64 = 50 * 1024 * 1024 // 50 MB
	var maxSize int64 = maxImageSize

	if strings.HasPrefix(mimeType, "video/") || strings.HasPrefix(header.Header.Get("Content-Type"), "video/") {
		maxSize = maxVideoSize
	}

	if header.Size > maxSize {
		maxSizeMB := maxSize / (1024 * 1024)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large. Maximum allowed size is %dMB.", maxSizeMB)})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueName := uuid.New().String() + ext

	// Handle specific folder paths (e.g., profiles, projects/projectABC)
	folderPath := c.PostForm("folder")
	finalUploadDir := h.uploadDir
	urlPrefix := "uploads"

	if folderPath != "" {
		cleanFolder := filepath.Clean(folderPath)
		// Ensure it doesn't try to escape into absolute roots or parents
		if !strings.HasPrefix(cleanFolder, "..") && cleanFolder != "." {
			finalUploadDir = filepath.Join(h.uploadDir, cleanFolder)
			urlPrefix = fmt.Sprintf("uploads/%s", strings.ReplaceAll(cleanFolder, "\\", "/"))
		}
	}

	// Ensure structural directory exists
	if err := os.MkdirAll(finalUploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory structure"})
		return
	}

	// Create destination
	destPath := filepath.Join(finalUploadDir, uniqueName)
	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file on server"})
		return
	}
	defer out.Close()

	// Copy content
	size, err := io.Copy(out, file)
	if err != nil {
		os.Remove(destPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Build accessible URL
	fileURL := fmt.Sprintf("%s/%s/%s", strings.TrimRight(h.baseURL, "/"), urlPrefix, uniqueName)

	// Save to database
	media := models.Media{
		FileName:     uniqueName,
		OriginalName: header.Filename,
		FilePath:     destPath,
		FileURL:      fileURL,
		FileType:     getFileCategory(mimeType),
		FileSize:     size,
		MimeType:     mimeType,
		UploadedAt:   time.Now(),
	}

	if err := h.db.Create(&media).Error; err != nil {
		// Clean up file if DB insert fails
		os.Remove(destPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           media.ID,
			"fileName":     media.FileName,
			"originalName": media.OriginalName,
			"fileUrl":      media.FileURL,
			"fileType":     media.FileType,
			"fileSize":     media.FileSize,
			"mimeType":     media.MimeType,
			"uploadedAt":   media.UploadedAt,
		},
		"message": "File uploaded successfully",
	})
}

// GetAll returns all media records
func (h *Handler) GetAll(c *gin.Context) {
	var mediaList []models.Media
	if err := h.db.Order("uploaded_at desc").Find(&mediaList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mediaList})
}

// Delete removes a media file from disk and database
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	var media models.Media
	if err := h.db.First(&media, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// Remove file from disk
	os.Remove(media.FilePath)

	// Remove from database
	if err := h.db.Delete(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

func isAllowedMIME(mime string) bool {
	allowed := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
		"application/pdf",
		"video/mp4",
		"video/webm",
		"video/ogg",
	}
	mime = strings.ToLower(strings.TrimSpace(mime))
	for _, a := range allowed {
		if strings.HasPrefix(mime, a) {
			return true
		}
	}
	return false
}

func getFileCategory(mime string) string {
	if strings.HasPrefix(mime, "image/") {
		return "image"
	}
	if strings.HasPrefix(mime, "video/") {
		return "video"
	}
	if strings.HasPrefix(mime, "application/pdf") {
		return "document"
	}
	return "file"
}
