package analytics

import (
	"time"
	"web-porto-backend/internal/adapters/websocket"
	"web-porto-backend/internal/domain/models"
	repo "web-porto-backend/internal/repositories/analytics"
)

type ViewStats struct {
	Total  int64 `json:"total"`
	Today  int64 `json:"today"`
	Week   int64 `json:"week"`
	Month  int64 `json:"month"`
	Unique int64 `json:"unique"`
}

type Service interface {
	TrackView(v *models.PageView) (*ViewStats, error)
	GetStats(page string) (*ViewStats, error)
	GetStatsWithFilter(page string, start, end *string, country string) (*ViewStats, error)
	GetTimeSeries(page string, start, end string, interval string) ([]repo.TimeSeriesPoint, error)
	SetWebsocketManager(wsManager *websocket.Manager)
}

type service struct {
	repo      repo.Repository
	wsManager *websocket.Manager
}

func NewService(r repo.Repository) Service {
	return &service{repo: r, wsManager: nil}
}

func (s *service) SetWebsocketManager(wsManager *websocket.Manager) {
	s.wsManager = wsManager
}

func (s *service) TrackView(v *models.PageView) (*ViewStats, error) {
	if err := s.repo.TrackView(v); err != nil {
		return nil, err
	}
	st, err := s.repo.GetStats(v.Page)
	if err != nil {
		return nil, err
	}

	// Create view stats object to return
	viewStats := &ViewStats{
		Total:  st.Total,
		Today:  st.Today,
		Week:   st.Week,
		Month:  st.Month,
		Unique: st.Unique,
	}

	// Broadcast updated stats via WebSocket if manager exists
	if s.wsManager != nil {
		// Create ViewCountsUpdate struct for websocket broadcast
		viewCountsUpdate := websocket.ViewCountsUpdate{
			Total:  st.Total,
			Today:  st.Today,
			Week:   st.Week,
			Month:  st.Month,
			Unique: st.Unique,
			Page:   v.Page,
		}

		// Update view counts through websocket
		s.wsManager.UpdateViewCounts(viewCountsUpdate, v.Page)
	}

	return viewStats, nil
}

func (s *service) GetStats(page string) (*ViewStats, error) {
	st, err := s.repo.GetStats(page)
	if err != nil {
		return nil, err
	}
	return &ViewStats{Total: st.Total, Today: st.Today, Week: st.Week, Month: st.Month, Unique: st.Unique}, nil
}

func (s *service) GetStatsWithFilter(page string, start, end *string, country string) (*ViewStats, error) {
	var startT *time.Time
	var endT *time.Time
	if start != nil && *start != "" {
		if t, err := time.Parse(time.RFC3339, *start); err == nil {
			startT = &t
		}
	}
	if end != nil && *end != "" {
		if t, err := time.Parse(time.RFC3339, *end); err == nil {
			endT = &t
		}
	}
	st, err := s.repo.GetStatsWithFilter(page, startT, endT, country)
	if err != nil {
		return nil, err
	}
	return &ViewStats{Total: st.Total, Today: st.Today, Week: st.Week, Month: st.Month, Unique: st.Unique}, nil
}

func (s *service) GetTimeSeries(page, start, end, interval string) ([]repo.TimeSeriesPoint, error) {
	if start == "" || end == "" {
		return []repo.TimeSeriesPoint{}, nil
	}
	sT, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, err
	}
	eT, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTimeSeries(page, sT, eT, interval)
}
