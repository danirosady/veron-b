package response

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPagination(page, perPage int, total int64) *PaginationMeta {
	tp := int(total) / perPage
	if int(total)%perPage > 0 {
		tp++
	}
	return &PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: tp,
	}
}
