package utils

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Items []*T           `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

// NormalizePagination validates and normalizes page and limit
func NormalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// CalculatePaginationMeta creates pagination metadata from total count
func CalculatePaginationMeta(total int64, page, limit int) PaginationMeta {
	totalPages := int(total / int64(limit))
	if total%int64(limit) != 0 {
		totalPages++
	}
	return PaginationMeta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// CreatePaginatedResponse creates a paginated response
func CreatePaginatedResponse[T any](items []*T, total int64, page, limit int) *PaginatedResponse[T] {
	if items == nil {
		items = make([]*T, 0)
	}
	return &PaginatedResponse[T]{
		Items: items,
		Meta:  CalculatePaginationMeta(total, page, limit),
	}
}
