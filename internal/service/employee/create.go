package employee

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/model"
)

func (s *service) Create(ctx context.Context, in CreateInput) (*model.Employee, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, appErr.ErrRequiredField
	}
	if len(name) > 255 {
		return nil, appErr.ErrInvalidRange
	}
	if in.Department == 0 {
		return nil, appErr.ErrRequiredField
	}

	exists, err := s.deptRepo.ExistsByID(ctx, in.Department)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, appErr.ErrNotFound
	}

	var addr string
	if in.Address != nil {
		addr = strings.TrimSpace(*in.Address)
	}

	emp := &model.Employee{
		Name:         name,
		Address:      addr,
		DepartmentID: in.Department,
	}

	if err := s.empRepo.Create(ctx, emp); err != nil {
		if isDuplicateKey(err) {
			return nil, appErr.ErrAlreadyExists
		}

		if isForeignKeyConstraint(err) {
			return nil, appErr.ErrInvalidInput
		}
		return nil, err
	}

	out, err := s.empRepo.GetEmpByIdJoinDept(ctx, emp.ID)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func isDuplicateKey(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}
	return false
}

func isForeignKeyConstraint(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1452 || me.Number == 1451
	}
	return false
}
