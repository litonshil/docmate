package types

type Pagination struct {
	Page     int   `json:"page" query:"page"`
	Limit    int   `json:"limit" query:"limit"`
	Total    int64 `json:"total"`
	LastPage int   `json:"last_page"`
}

type PaginatedResponse[T any] struct {
	Pagination
	Records []T `json:"records"`
}
