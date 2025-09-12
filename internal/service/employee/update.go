package employee

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"

	"gorm.io/gorm"
)

func (s *service) UpdateEmployeeByEmployeeID(ctx context.Context, employeeID string, in UpdateInput) (*model.Employee, error) {
	if helper.NormalizeStringField(employeeID) == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", appErr.ErrRequiredField)
	}

	if in.Address != nil {
		a := strings.TrimSpace(*in.Address)
		in.Address = &a
	}

	if in.Department != 0 {
		exists, err := s.deptRepo.ExistsByID(ctx, in.Department)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, appErr.ErrNotFound
		}
	}

	if err := s.empRepo.UpdateByEmployeeID(ctx, employeeID, emprepo.UpdateParams{
		Name:       &name,
		Address:    in.Address,
		Department: in.Department,
	}); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, appErr.ErrNotFound
		case isForeignKeyConstraint(err):
			return nil, appErr.ErrInvalidInput
		default:
			return nil, err
		}
	}

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, employeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.ErrNotFound
		}
		return nil, err
	}
	return emp, nil
}
