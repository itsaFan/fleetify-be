package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"gorm.io/gorm"
)

func (s *service) DeleteByEmployeeID(ctx context.Context, name string) error {
	norm := helper.NormalizeStringField(name)
	if norm == "" {
		return fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	if err := s.empRepo.DeleteByEmployeeID(ctx, norm); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: employee %q", appErr.ErrNotFound, norm)
		}

		return err
	}
	return nil
}
