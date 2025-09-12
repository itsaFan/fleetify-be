package employee

import (
	"context"

	"github.com/itsaFan/fleetify-be/internal/model"
	deptrepo "github.com/itsaFan/fleetify-be/internal/repo/department"
	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"
)

type service struct {
	empRepo  emprepo.Repository
	deptRepo deptrepo.Repository
}

type Service interface {
	Create(ctx context.Context, in CreateInput) (*model.Employee, error)
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	GetByEmployeeID(ctx context.Context, name string) (*model.Employee, error)
	UpdateEmployeeByEmployeeID(ctx context.Context, employeeID string, in UpdateInput) (*model.Employee, error)
	DeleteByEmployeeID(ctx context.Context, employeeID string) error
}

func New(empRepo emprepo.Repository, deptRepo deptrepo.Repository) Service {
	return &service{empRepo: empRepo, deptRepo: deptRepo}
}
