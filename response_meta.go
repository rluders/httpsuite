package httpsuite

// Response represents the structure of an HTTP response, including an optional body and metadata.
type Response[T any] struct {
	Data T   `json:"data"`
	Meta any `json:"meta,omitempty"`
}

// PageMeta provides page-based pagination details.
type PageMeta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

// Meta is kept as a compatibility alias for page-based pagination metadata.
type Meta = PageMeta

// CursorMeta provides cursor-based pagination details.
type CursorMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}

// NewPageMeta builds page-based metadata and derives total pages when possible.
func NewPageMeta(page, pageSize, totalItems int) *PageMeta {
	meta := &PageMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
	}
	if pageSize > 0 && totalItems > 0 {
		meta.TotalPages = (totalItems + pageSize - 1) / pageSize
	}
	return meta
}

// NewCursorMeta builds cursor-based metadata.
func NewCursorMeta(nextCursor, prevCursor string, hasNext, hasPrev bool) *CursorMeta {
	return &CursorMeta{
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}
