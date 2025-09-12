package department

import (
	"context"
	"fmt"
	"errors"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	"gorm.io/gorm"
)

func (s *service) GetByName(ctx context.Context, name string) (*model.Department, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: department_name is required", appErr.ErrRequiredField)
	}
	normalized := helper.NormalizeStringField(name)

	d, err := s.repo.GetByName(ctx, normalized)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department %q", appErr.ErrNotFound, normalized)
		}
		return nil, err
	}
	return d, nil
}
