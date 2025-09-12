package employee

import (
	"context"
	"strings"

	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"
)

func (in *ListInput) normalize() {
	if in.Limit <= 0 || in.Limit > 100 {
		in.Limit = 10
	}
	if in.Page <= 0 {
		in.Page = 1
	}
	switch strings.ToLower(strings.TrimSpace(in.SortBy)) {
	case "id":
		in.SortBy = "employees.id"
	case "department_name":
		in.SortBy = "departments.department_name"
	default:
		in.SortBy = "employees.name"
	}
	if strings.EqualFold(in.SortDir, "desc") {
		in.SortDir = "desc"
	} else {
		in.SortDir = "asc"
	}
}

func (s *service) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	in.normalize()

	items, total, err := s.empRepo.ListJoinDept(ctx, emprepo.ListParams{
		Search:  in.Search,
		Limit:   in.Limit,
		Page:    in.Page,
		SortBy:  in.SortBy,
		SortDir: in.SortDir,
	})

	if err != nil {
		return nil, err
	}

	totalPages := 0
	if in.Limit > 0 {
		totalPages = int((total + int64(in.Limit) - 1) / int64(in.Limit))
	}

	return &ListOutput{
		Data: items,
		Pagination: Pagination{
			TotalData:   total,
			CurrentPage: in.Page,
			TotalPages:  totalPages,
			HasNextPage: in.Page < totalPages,
			HasPrevPage: in.Page > 1,
		},
	}, nil

}
