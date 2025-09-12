package department

import (
	"context"
	"errors"
	"fmt"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"gorm.io/gorm"
)

func (s *service) DeleteByName(ctx context.Context, name string) error {
	norm := helper.NormalizeStringField(name)
	if norm == "" {
		return fmt.Errorf("%w: department_name is required", appErr.ErrRequiredField)
	}

	if err := s.repo.DeleteByName(ctx, norm); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: department %q", appErr.ErrNotFound, norm)
		}

		return err
	}
	return nil
}
