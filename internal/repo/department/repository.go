package department

import (
	"context"
	"strings"

	"github.com/itsaFan/fleetify-be/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, d *model.Department) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, p ListParams) ([]model.Department, int64, error)
	GetByName(ctx context.Context, name string) (*model.Department, error)
	UpdateByName(ctx context.Context, name string, p UpdateParams) error
	DeleteByName(ctx context.Context, name string) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, d *model.Department) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Department{}).Where("department_name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}

type ListParams struct {
	Search  string
	Limit   int
	Page    int
	SortBy  string
	SortDir string
}

func (r *repository) List(ctx context.Context, p ListParams) ([]model.Department, int64, error) {
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	sortBy := "department_name"
	switch strings.ToLower(strings.TrimSpace(p.SortBy)) {
	case "id":
		sortBy = "id"
	case "department_name", "":
	default:
		sortBy = "department_name"
	}

	dir := "ASC"
	if strings.EqualFold(p.SortDir, "desc") {
		dir = "DESC"
	}

	q := r.db.WithContext(ctx).Model(&model.Department{})
	if s := strings.TrimSpace(p.Search); s != "" {
		q = q.Where("department_name LIKE ?", "%"+s+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Department
	if err := q.
		Order(sortBy + " " + dir).
		Limit(p.Limit).
		Offset((p.Page - 1) * p.Limit).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*model.Department, error) {
	var d model.Department

	err := r.db.WithContext(ctx).Where("department_name = ?", name).First(&d).Error

	if err != nil {
		return nil, err
	}

	return &d, nil
}

type UpdateParams struct {
	DepartmentName  *string
	MaxClockInTime  *string
	MaxClockOutTime *string
}

func (r *repository) UpdateByName(ctx context.Context, name string, p UpdateParams) error {
	updates := map[string]any{}

	if p.DepartmentName != nil {
		updates["department_name"] = strings.TrimSpace(*p.DepartmentName)
	}

	if p.MaxClockInTime != nil {
		updates["max_clock_in_time"] = *p.MaxClockInTime
	}

	if p.MaxClockOutTime != nil {
		updates["max_clock_out_time"] = *p.MaxClockOutTime
	}

	if len(updates) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).
		Model(&model.Department{}).
		Where("department_name = ?", name).
		Updates(updates)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}

func (r *repository) DeleteByName(ctx context.Context, name string) error {
	tx := r.db.WithContext(ctx).
		Where("department_name = ?", name).
		Delete(&model.Department{})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
