package department

import (
	"context"
	"errors"
	"fmt"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	deptrepo "github.com/itsaFan/fleetify-be/internal/repo/department"
	"gorm.io/gorm"
)

func (in UpdateInput) isEmpty() bool {
	return in.DepartmentName == nil && in.MaxClockIn == nil && in.MaxClockOut == nil
}

func (s *service) UpdateByName(ctx context.Context, currentName string, in UpdateInput) (*model.Department, error) {
	if helper.NormalizeStringField(currentName) == "" {
		return nil, fmt.Errorf("%w: department_name is required", appErr.ErrRequiredField)
	}
	ident := helper.NormalizeStringField(currentName)

	cur, err := s.repo.GetByName(ctx, ident)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department %q", appErr.ErrNotFound, ident)
		}
		return nil, err
	}

	finalName := cur.DepartmentName
	if in.DepartmentName != nil {
		nm := helper.NormalizeStringField(*in.DepartmentName)
		if nm == "" {
			return nil, fmt.Errorf("%w: department_name is required", appErr.ErrRequiredField)
		}

		if nm != cur.DepartmentName {
			exists, err := s.repo.ExistsByName(ctx, nm)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, fmt.Errorf("%w: department %q", appErr.ErrAlreadyExists, nm)
			}
		}
		finalName = nm
	}

	finalIn := cur.MaxClockInTime
	finalOut := cur.MaxClockOutTime

	if in.MaxClockIn != nil {
		if _, err := helper.ParseTimeOfDay(*in.MaxClockIn); err != nil {
			return nil, fmt.Errorf("%w: max_clock_in_time invalid: %v", appErr.ErrInvalidInput, err)
		}
		finalIn = *in.MaxClockIn
	}
	if in.MaxClockOut != nil {
		if _, err := helper.ParseTimeOfDay(*in.MaxClockOut); err != nil {
			return nil, fmt.Errorf("%w: max_clock_out_time invalid: %v", appErr.ErrInvalidInput, err)
		}
		finalOut = *in.MaxClockOut
	}

	if in.isEmpty() {
		return cur, nil
	}

	inT, err := helper.ParseTimeOfDay(finalIn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", appErr.ErrInvalidInput, err)
	}
	outT, err := helper.ParseTimeOfDay(finalOut)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", appErr.ErrInvalidInput, err)
	}
	if !inT.Before(outT) {
		return nil, fmt.Errorf("%w: max_clock_in_time must be earlier than max_clock_out_time", appErr.ErrInvalidTimeRange)
	}

	up := deptrepo.UpdateParams{}
	if in.DepartmentName != nil && finalName != cur.DepartmentName {
		up.DepartmentName = &finalName
	}
	if in.MaxClockIn != nil {
		up.MaxClockInTime = &finalIn
	}
	if in.MaxClockOut != nil {
		up.MaxClockOutTime = &finalOut
	}

	if err := s.repo.UpdateByName(ctx, ident, up); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: department %q", appErr.ErrNotFound, ident)
		}
		return nil, err
	}

	updatedName := finalName
	d, err := s.repo.GetByName(ctx, updatedName)
	if err != nil {
		return nil, err
	}
	return d, nil
}
