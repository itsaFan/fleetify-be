package employee

import "github.com/itsaFan/fleetify-be/internal/model"

type Pagination struct {
	TotalData   int64 `json:"totalData"`
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

type CreateInput struct {
	Name       string
	Address    *string
	Department uint64
}

type ListInput struct {
	Search  string
	Limit   int
	Page    int
	SortBy  string
	SortDir string
}
type ListOutput struct {
	Data       []model.Employee `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

type UpdateInput struct {
	Name       string
	Address    *string
	Department uint64
}
