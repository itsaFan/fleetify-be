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
	ListHistoryByEmpId(ctx context.Context, p ListParamsEmp) ([]model.AttendanceHistory, int64, error)
	ListHistoryByDepartment(ctx context.Context, p ListParamsDept) ([]model.AttendanceHistory, int64, error)

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
	Limit      int
	Page       int
	FromUtc    time.Time
	ToUtc      time.Time
}

func (r *repository) ListHistoryByEmpId(ctx context.Context, p ListParamsEmp) ([]model.AttendanceHistory, int64, error) {

	empId := strings.TrimSpace(p.EmployeeID)
	if empId == "" {
		return nil, 0, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	offset := (p.Page - 1) * p.Limit

	q := r.db.WithContext(ctx).Model(&model.AttendanceHistory{}).
		Where("employee_id = ?", empId)

	if !p.FromUtc.IsZero() {
		q = q.Where("date_attendance >= ?", p.FromUtc)
	}

	if !p.ToUtc.IsZero() {
		q = q.Where("date_attendance <= ?", p.ToUtc)
	}

	var total int64
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.AttendanceHistory
	if err := q.
		Order("date_attendance ASC, id ASC").
		Limit(p.Limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

type ListParamsDept struct {
	Limit   int
	Page    int
	FromUtc time.Time
	ToUtc   time.Time
}


func (r *repository) ListHistoryByDepartment(ctx context.Context, p ListParamsDept) ([]model.AttendanceHistory, int64, error) {
	
}