package department

import (
	"context"
	"fmt"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
)

func (in CreateInput) validate() error {
	if in.DepartmentName == "" {
		return fmt.Errorf("%w: department_name is required", appErr.ErrRequiredField)
	}
	if _, err := helper.ParseTimeOfDay(in.MaxClockIn); err != nil {
		return fmt.Errorf("%w: max_clock_in_time invalid: %v", appErr.ErrInvalidInput, err)
	}
	if _, err := helper.ParseTimeOfDay(in.MaxClockOut); err != nil {
		return fmt.Errorf("%w: max_clock_out_time invalid: %v", appErr.ErrInvalidInput, err)
	}
	return nil
}

func (s *service) Create(ctx context.Context, in CreateInput) (*model.Department, error) {
	// Validate payload
	// fmt.Println("payload", in)
	if err := in.validate(); err != nil {
		return nil, err
	}

	name := helper.NormalizeStringField(in.DepartmentName)

	// validate
	clockIn, err := helper.ParseTimeOfDay(in.MaxClockIn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", appErr.ErrInvalidInput, err)
	}
	clockOut, err := helper.ParseTimeOfDay(in.MaxClockOut)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", appErr.ErrInvalidInput, err)
	}
	if !clockIn.Before(clockOut) {
		return nil, fmt.Errorf("%w: max_clock_in_time must be earlier than max_clock_out_time", appErr.ErrInvalidRange)
	}

	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: department %q", appErr.ErrAlreadyExists, name)
	}

	dept := &model.Department{
		DepartmentName:  helper.NormalizeStringField(in.DepartmentName),
		MaxClockInTime:  in.MaxClockIn,
		MaxClockOutTime: in.MaxClockOut,
	}

	if err := s.repo.Create(ctx, dept); err != nil {
		return nil, err
	}

	return dept, nil
}
