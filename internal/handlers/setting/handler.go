package setting

import (
	"net/http"
	"strings"
	"web-porto-backend/internal/services/setting"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service setting.Service
}

func NewHandler(service setting.Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetAll(c *gin.Context) {
	// If comma-separated keys are provided in query
	keysQuery := c.Query("keys")

	var data map[string]string
	var err error

	if keysQuery != "" {
		keys := strings.Split(keysQuery, ",")
		data, err = h.service.GetSettings(keys)
	} else {
		data, err = h.service.GetAllSettings()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) Update(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.service.SaveSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings saved successfully", "data": input})
}
