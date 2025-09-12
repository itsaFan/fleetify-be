package department

import "github.com/itsaFan/fleetify-be/internal/model"

type Pagination struct {
	TotalData   int64 `json:"totalData"`
	CurrentPage int   `json:"currentPage"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

type CreateInput struct {
	DepartmentName string
	// "HH:MM:SS"
	MaxClockIn  string
	MaxClockOut string
}

type ListInput struct {
	Search  string
	Limit   int
	Page    int
	SortBy  string
	SortDir string
}

type ListOutput struct {
	Data       []model.Department `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type UpdateInput struct {
	DepartmentName *string
	MaxClockIn     *string
	MaxClockOut    *string
}
