package analytics

import (
	"time"
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
}

type service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) Service {
	return &service{repo: r}
}

func (s *service) TrackView(v *models.PageView) (*ViewStats, error) {
	if err := s.repo.TrackView(v); err != nil {
		return nil, err
	}
	st, err := s.repo.GetStats(v.Page)
	if err != nil {
		return nil, err
	}
	return &ViewStats{Total: st.Total, Today: st.Today, Week: st.Week, Month: st.Month, Unique: st.Unique}, nil
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
