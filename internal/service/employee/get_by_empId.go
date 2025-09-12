package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	"gorm.io/gorm"
)

func (s *service) GetByEmployeeID(ctx context.Context, employeeID string) (*model.Employee, error) {
	if employeeID == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}
	normalized := helper.NormalizeStringField(employeeID)

	d, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, normalized)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee %q", appErr.ErrNotFound, normalized)
		}
		return nil, err
	}
	return d, nil
}
