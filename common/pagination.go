package common

type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNextPage  bool  `json:"has_next_page"`
	HasPrevPage  bool  `json:"has_prev_page"`
	Limit        int   `json:"limit"`
}

func CalculatePagination(totalRecords, currentPage, limit int) *PaginationInfo {
	totalPages := (totalRecords + limit - 1) / limit

	if totalPages == 0 {
		totalPages = 1
	}

	hasNextPage := currentPage < totalPages
	hasPrevPage := currentPage > 1

	return &PaginationInfo{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		TotalRecords: int64(totalRecords),
		HasNextPage:  hasNextPage,
		HasPrevPage:  hasPrevPage,
		Limit:        limit,
	}
}
