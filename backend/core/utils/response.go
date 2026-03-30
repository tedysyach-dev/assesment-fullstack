package utils

type WebResponse[T any] struct {
	Status  bool          `json:"status"`
	Message string        `json:"message,omitempty"`
	Code    int           `json:"code,omitempty"`
	Paging  *PageMetadata `json:"meta,omitempty"`
	Data    T             `json:"resource"`
	Errors  string        `json:"errors,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T           `json:"data,omitempty"`
	PageMetadata *PageMetadata `json:"meta,omitempty"`
}

type PageMetadata struct {
	Page      *int   `json:"page,omitempty"`
	PerPage   *int   `json:"perPage,omitempty"`
	TotalItem *int64 `json:"totalItem,omitempty"`
	TotalPage *int64 `json:"totalPage,omitempty"`
}

func NewPageMetadata(page, limit int, total int64) *PageMetadata {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int64((total + int64(limit) - 1) / int64(limit)) // ceil division
	return &PageMetadata{
		Page:      &page,
		PerPage:   &limit,
		TotalItem: &total,
		TotalPage: &totalPages,
	}
}
