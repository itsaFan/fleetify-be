package employee

import (
	"context"
	"strings"

	"github.com/itsaFan/fleetify-be/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, d *model.Employee) error
	ListJoinDept(ctx context.Context, p ListParams) ([]model.Employee, int64, error)
	GetEmpByIdJoinDept(ctx context.Context, id uint64) (*model.Employee, error)
	GetByEmployeeIDJoinDept(ctx context.Context, employeeID string) (*model.Employee, error)
	UpdateByEmployeeID(ctx context.Context, employeeID string, p UpdateParams) error
	DeleteByEmployeeID(ctx context.Context, employeeID string) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, d *model.Employee) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *repository) GetEmpByIdJoinDept(ctx context.Context, id uint64) (*model.Employee, error) {
	var out model.Employee
	if err := r.db.WithContext(ctx).
		Preload("Department").
		First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

type ListParams struct {
	Search  string
	Limit   int
	Page    int
	SortBy  string
	SortDir string
}

func (r *repository) ListJoinDept(ctx context.Context, p ListParams) ([]model.Employee, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Employee{})

	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	sortBy := "employees.name"
	switch strings.ToLower(strings.TrimSpace(p.SortBy)) {
	case "id":
		sortBy = "employees.id"
	case "name", "":
		sortBy = "employees.name"
	}

	dir := "ASC"
	if strings.EqualFold(p.SortDir, "desc") {
		dir = "DESC"
	}

	if s := strings.TrimSpace(p.Search); s != "" {
		q = q.Joins("JOIN departments ON employees.department_id = departments.id").
			Where("employees.name LIKE ? OR departments.department_name LIKE ?", "%"+s+"%", "%"+s+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Employee
	if err := q.Preload("Department").
		Order(sortBy + " " + dir).
		Limit(p.Limit).
		Offset((p.Page - 1) * p.Limit).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil

}

type UpdateParams struct {
	Name       *string
	Address    *string
	Department uint64
}

func (r *repository) UpdateByEmployeeID(ctx context.Context, employeeID string, p UpdateParams) error {
	updates := map[string]any{}

	if p.Name != nil {
		updates["name"] = strings.TrimSpace(*p.Name)
	}

	if p.Address != nil {
		updates["address"] = strings.TrimSpace(*p.Address)
	}

	// 0 = "no change"
	if p.Department != 0 {
		updates["department_id"] = p.Department
	}

	if len(updates) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).
		Model(&model.Employee{}).
		Where("employee_id = ?", employeeID).
		Updates(updates)

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) GetByEmployeeIDJoinDept(ctx context.Context, employeeID string) (*model.Employee, error) {
	var out model.Employee
	if err := r.db.WithContext(ctx).
		Preload("Department").
		First(&out, "employee_id = ?", employeeID).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) DeleteByEmployeeID(ctx context.Context, employeeID string) error {
	tx := r.db.WithContext(ctx).
		Where("employee_id = ?", employeeID).
		Delete(&model.Employee{})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
