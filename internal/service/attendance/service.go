package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	atdrepo "github.com/itsaFan/fleetify-be/internal/repo/attendance"
	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"
	"gorm.io/gorm"
)

type service struct {
	atdRepo atdrepo.Repository
	empRepo emprepo.Repository
}

type Service interface {
	CreateEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error)
	CloseEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error)
}

func New(atdRepo atdrepo.Repository, empRepo emprepo.Repository) Service {
	return &service{atdRepo: atdRepo, empRepo: empRepo}
}

func (s *service) CreateEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error) {
	if employeeID == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	normalizedEmpId := helper.NormalizeStringField(employeeID)

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, normalizedEmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee %q", appErr.ErrNotFound, normalizedEmpId)
		}
		return nil, err
	}

	_ = emp

	now := time.Now().UTC()
	attID := uuid.New().String()

	att := &model.Attendance{
		EmployeeID:   normalizedEmpId,
		AttendanceID: attID,
		ClockIn:      &now,
	}

	hist := &model.AttendanceHistory{
		EmployeeID:     normalizedEmpId,
		AttendanceID:   attID,
		DateAttendance: now,
		AttendanceType: 1,
		Description:    "Clock in",
	}

	if err := s.atdRepo.WithTx(ctx, func(tx atdrepo.Repository) error {
		open, err := tx.FindEmpOpenAttendanceForUpdate(ctx, normalizedEmpId)
		if err != nil {
			return err
		}
		if open != nil {
			return fmt.Errorf("%w: already clocked in with attendance_id=%s", appErr.ErrAlreadyExists, open.AttendanceID)
		}
		if err := tx.CreateEmpAttendanceByEmpId(ctx, att); err != nil {
			return err
		}
		if err := tx.CreateAttendanceHistory(ctx, hist); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return att, nil
}

func (s *service) CloseEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error) {
	if employeeID == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	normalizedEmpId := helper.NormalizeStringField(employeeID)

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, normalizedEmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee %q", appErr.ErrNotFound, normalizedEmpId)
		}
		return nil, err
	}

	_ = emp

	now := time.Now().UTC()
	var updated *model.Attendance

	if err := s.atdRepo.WithTx(ctx, func(tx atdrepo.Repository) error {
		open, err := tx.FindEmpOpenAttendanceForUpdate(ctx, normalizedEmpId)

		if err != nil {
			return err
		}

		if open == nil {
			return fmt.Errorf("%w: no open attendance for employee %q", appErr.ErrNotFound, normalizedEmpId)
		}

		if err := tx.UpdateAttendanceOutByAttendanceID(ctx, open.AttendanceID, now); err != nil {
			return err
		}

		hist := &model.AttendanceHistory{
			EmployeeID:     normalizedEmpId,
			AttendanceID:   open.AttendanceID,
			DateAttendance: now,
			AttendanceType: 2,
			Description:    "Clock out",
		}

		if err := tx.CreateAttendanceHistory(ctx, hist); err != nil {
			return err
		}

		open.ClockOut = &now
		updated = open
		return nil

	}); err != nil {
		return nil, err
	}

	return updated, nil

}
