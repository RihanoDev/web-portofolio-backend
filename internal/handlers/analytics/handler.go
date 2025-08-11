package analytics

import (
	"net/http"
	"time"
	"web-porto-backend/common/response"
	httpAdapter "web-porto-backend/internal/adapters/http"
	"web-porto-backend/internal/domain/models"
	svc "web-porto-backend/internal/services/analytics"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     svc.Service
	httpAdapter *httpAdapter.HTTPAdapter
}

func NewHandler(service svc.Service, httpAdapter *httpAdapter.HTTPAdapter) *Handler {
	return &Handler{service: service, httpAdapter: httpAdapter}
}

type trackRequest struct {
	Page      string `json:"page" binding:"required"`
	VisitorID string `json:"visitorId" binding:"required"`
	UserAgent string `json:"userAgent"`
	SessionID string `json:"sessionId"`
	Referrer  string `json:"referrer"`
}

// POST /api/v1/analytics/track
func (h *Handler) Track(c *gin.Context) {
	var req trackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Derive IP (best effort)
	ip := c.ClientIP()
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}

	v := &models.PageView{
		Page:      req.Page,
		VisitorID: req.VisitorID,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
		IP:        ip,
		Timestamp: time.Now(),
	}

	stats, err := h.service.TrackView(v)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to track view", err.Error()))
		return
	}

	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, stats, "View tracked successfully")
}

// GET /api/v1/analytics/views?page=/path
func (h *Handler) GetViews(c *gin.Context) {
	page := c.Query("page")
	stats, err := h.service.GetStats(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get view stats", err.Error()))
		return
	}
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, stats, "View stats fetched")
}

// GET /api/v1/analytics?startDate=&endDate=&page=
// For now, returns the same aggregate as GetViews; filters can be applied later.
func (h *Handler) GetAnalytics(c *gin.Context) {
	page := c.Query("page")
	start := c.Query("startDate")
	end := c.Query("endDate")
	country := c.Query("country")
	// pass pointers only if non-empty so service can ignore when empty
	var startPtr *string
	var endPtr *string
	if start != "" {
		startPtr = &start
	}
	if end != "" {
		endPtr = &end
	}
	stats, err := h.service.GetStatsWithFilter(page, startPtr, endPtr, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get analytics", err.Error()))
		return
	}
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, stats, "Analytics fetched")
}

// GET /api/v1/analytics/series?page=&startDate=&endDate=&interval=hour|day
func (h *Handler) GetSeries(c *gin.Context) {
	page := c.Query("page")
	start := c.Query("startDate")
	end := c.Query("endDate")
	interval := c.DefaultQuery("interval", "day")
	points, err := h.service.GetTimeSeries(page, start, end, interval)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("Failed to get timeseries", err.Error()))
		return
	}
	h.httpAdapter.SendSuccessResponse(c, http.StatusOK, points, "Timeseries fetched")
}
