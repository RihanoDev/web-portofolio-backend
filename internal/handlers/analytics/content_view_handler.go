package analytics

import (
	"net/http"
	"strconv"

	applog "web-porto-backend/common/logger"
	"web-porto-backend/common/response"
	svc "web-porto-backend/internal/services/analytics"

	"github.com/gin-gonic/gin"
)

// Add these to the existing Handler struct
type trackContentRequest struct {
	ContentID string `json:"contentId" binding:"required"`
	Type      string `json:"contentType" binding:"required"`
	VisitorID string `json:"visitorId"`
	UserAgent string `json:"userAgent"`
	Referrer  string `json:"referrer"`
}

type viewCountResponse struct {
	Count int `json:"count"`
}

// TrackContentView handles POST /api/v1/views/track
func (h *Handler) TrackContentView(c *gin.Context) {
	log := applog.GetLogger().WithFields(applog.Fields{"handler": "analytics.TrackContentView"})

	var req trackContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid track content payload", applog.Fields{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Get visitor IP
	ip := c.ClientIP()
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}

	// Get visitor ID or generate one based on IP and user agent
	visitorID := req.VisitorID
	if visitorID == "" {
		visitorID = ip + "-" + req.UserAgent
	}

	// Convert type string to ContentType
	var contentType svc.ContentType
	switch req.Type {
	case "article":
		contentType = svc.ContentTypeArticle
	case "project":
		contentType = svc.ContentTypeProject
	case "page":
		contentType = svc.ContentTypePage
	default:
		contentType = svc.ContentTypePage
	}

	// Track view
	cvs, ok := h.service.(interface {
		TrackContentView(contentID string, contentType svc.ContentType, visitorID string, userAgent string, referrer string, ip string) error
	})

	if !ok {
		log.Error("service does not implement TrackContentView")
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Service not available", "Internal error"))
		return
	}

	err := cvs.TrackContentView(req.ContentID, contentType, visitorID, req.UserAgent, req.Referrer, ip)

	if err != nil {
		log.Error("track content view failed", applog.Fields{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to track content view", err.Error()))
		return
	}

	log.Info("track content view ok", applog.Fields{"contentId": req.ContentID, "type": req.Type})
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, nil, "Content view tracked successfully")
}

// GetContentViewCount handles GET /api/v1/views/count
func (h *Handler) GetContentViewCount(c *gin.Context) {
	log := applog.GetLogger().WithFields(applog.Fields{"handler": "analytics.GetContentViewCount"})

	contentID := c.Query("contentId")
	contentTypeStr := c.Query("contentType")

	if contentID == "" || contentTypeStr == "" {
		log.Warn("missing required parameters")
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Missing parameters", "contentId and contentType are required"))
		return
	}

	// Convert type string to ContentType
	var contentType svc.ContentType
	switch contentTypeStr {
	case "article":
		contentType = svc.ContentTypeArticle
	case "project":
		contentType = svc.ContentTypeProject
	case "page":
		contentType = svc.ContentTypePage
	default:
		contentType = svc.ContentTypePage
	}

	// Get view count
	cvs, ok := h.service.(interface {
		GetContentViewCount(contentID string, contentType svc.ContentType) (int, error)
	})

	if !ok {
		log.Error("service does not implement GetContentViewCount")
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Service not available", "Internal error"))
		return
	}

	count, err := cvs.GetContentViewCount(contentID, contentType)

	if err != nil {
		log.Error("get content view count failed", applog.Fields{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get content view count", err.Error()))
		return
	}

	log.Info("get content view count ok", applog.Fields{"contentId": contentID, "type": contentTypeStr, "count": count})
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, viewCountResponse{Count: count}, "Content view count retrieved successfully")
}

// GetContentViewAnalytics handles GET /api/v1/views/analytics
func (h *Handler) GetContentViewAnalytics(c *gin.Context) {
	log := applog.GetLogger().WithFields(applog.Fields{"handler": "analytics.GetContentViewAnalytics"})

	contentID := c.Query("contentId")
	contentTypeStr := c.Query("contentType")
	period := c.DefaultQuery("period", "day")
	limitStr := c.DefaultQuery("limit", "30")

	if contentID == "" || contentTypeStr == "" {
		log.Warn("missing required parameters")
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Missing parameters", "contentId and contentType are required"))
		return
	}

	// Convert limit to int
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 30
	}

	// Convert type string to ContentType
	var contentType svc.ContentType
	switch contentTypeStr {
	case "article":
		contentType = svc.ContentTypeArticle
	case "project":
		contentType = svc.ContentTypeProject
	case "page":
		contentType = svc.ContentTypePage
	default:
		contentType = svc.ContentTypePage
	}

	// Get analytics data
	cvs, ok := h.service.(interface {
		GetContentViewAnalytics(contentID string, contentType svc.ContentType, period string, limit int) ([]svc.AnalyticsDataPoint, error)
	})

	if !ok {
		log.Error("service does not implement GetContentViewAnalytics")
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Service not available", "Internal error"))
		return
	}

	dataPoints, err := cvs.GetContentViewAnalytics(contentID, contentType, period, limit)

	if err != nil {
		log.Error("get content view analytics failed", applog.Fields{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get content view analytics", err.Error()))
		return
	}

	log.Info("get content view analytics ok", applog.Fields{"contentId": contentID, "type": contentTypeStr, "points": len(dataPoints)})
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, dataPoints, "Content view analytics retrieved successfully")
}
