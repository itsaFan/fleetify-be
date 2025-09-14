package attendance

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	WithTx(ctx context.Context, fn func(txRepo Repository) error) error

	FindEmpOpenAttendanceForUpdate(ctx context.Context, employeeID string) (*model.Attendance, error)
	ListHistoryByEmpId(ctx context.Context, p ListParamsEmp) ([]model.AttendanceHistory, error)
	ListHistoryByDepartment(ctx context.Context, p ListParamsDept) ([]model.AttendanceHistory, error)

	CreateEmpAttendanceByEmpId(ctx context.Context, d *model.Attendance) error
	CreateAttendanceHistory(ctx context.Context, d *model.AttendanceHistory) error
	UpdateAttendanceOutByAttendanceID(tx context.Context, attendanceID string, clockOut time.Time) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Transaction boundary
func (r *repository) WithTx(ctx context.Context, fn func(txRepo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &repository{db: tx}
		return fn(txRepo)
	})
}

// Functional ops
func (r *repository) FindEmpOpenAttendanceForUpdate(ctx context.Context, employeeID string) (*model.Attendance, error) {
	var att model.Attendance
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("employee_id = ? AND clock_out IS NULL", employeeID).
		Order("id DESC").
		First(&att).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &att, err
}

func (r *repository) CreateEmpAttendanceByEmpId(ctx context.Context, d *model.Attendance) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *repository) CreateAttendanceHistory(ctx context.Context, d *model.AttendanceHistory) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *repository) UpdateAttendanceOutByAttendanceID(ctx context.Context, attendanceID string, clockOut time.Time) error {
	tx := r.db.WithContext(ctx).
		Model(&model.Attendance{}).
		Where("attendance_id = ? AND clock_out IS NULL", attendanceID).
		Updates(map[string]any{"clock_out": clockOut})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type ListParamsEmp struct {
	EmployeeID string
	FromUtc    time.Time
	ToUtc      time.Time
}

func (r *repository) ListHistoryByEmpId(ctx context.Context, p ListParamsEmp) ([]model.AttendanceHistory, error) {

	empId := strings.TrimSpace(p.EmployeeID)
	if empId == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	q := r.db.WithContext(ctx).Model(&model.AttendanceHistory{}).
		Where("employee_id = ?", empId)

	if !p.FromUtc.IsZero() {
		q = q.Where("date_attendance >= ?", p.FromUtc)
	}

	if !p.ToUtc.IsZero() {
		q = q.Where("date_attendance <= ?", p.ToUtc)
	}

	var items []model.AttendanceHistory
	if err := q.
		Order("date_attendance ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

type ListParamsDept struct {
	DepartmentID *uint64
	FromUtc      time.Time
	ToUtc        time.Time
}

func (r *repository) ListHistoryByDepartment(ctx context.Context, p ListParamsDept) ([]model.AttendanceHistory, error) {

	q := r.db.WithContext(ctx).
		Model(&model.AttendanceHistory{}).
		Joins("JOIN employees e ON e.employee_id = attendance_histories.employee_id")

	if p.DepartmentID != nil {
		q = q.Where("e.department_id = ?", *p.DepartmentID)
	}

	if !p.FromUtc.IsZero() {
		q = q.Where("attendance_histories.date_attendance >= ?", p.FromUtc)
	}

	if !p.ToUtc.IsZero() {
		q = q.Where("attendance_histories.date_attendance <= ?", p.ToUtc)
	}

	var items []model.AttendanceHistory
	if err := q.
		Select("attendance_histories.*").
		Order("attendance_histories.date_attendance ASC, attendance_histories.id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil

}
