package department

import (
	"context"

	"github.com/itsaFan/fleetify-be/internal/model"
	deptrepo "github.com/itsaFan/fleetify-be/internal/repo/department"
)

type service struct {
	repo deptrepo.Repository
}

type Service interface {
	Create(ctx context.Context, in CreateInput) (*model.Department, error)
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	GetByName(ctx context.Context, name string) (*model.Department, error)
	UpdateByName(ctx context.Context, currentName string, in UpdateInput) (*model.Department, error)
	DeleteByName(ctx context.Context, name string) error
}

func New(repo deptrepo.Repository) Service {
	return &service{repo: repo}
}
