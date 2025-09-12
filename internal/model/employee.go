package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Employee struct {
	ID           uint64      `gorm:"primaryKey;autoIncrement;column:id"`
	EmployeeID   string    `gorm:"size:50;uniqueIndex;not null; column:employee_id"`
	DepartmentID uint64      `gorm:"not null;column:department_id"`
	Name         string    `gorm:"size:255;not null;column:name"`
	Address      string    `gorm:"type:text;column:address"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	// Relations
	Department  Department          `gorm:"foreignKey:DepartmentID;references:ID"`
	Attendances []Attendance        `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	History     []AttendanceHistory `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
}

func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if e.EmployeeID == "" {
		e.EmployeeID = uuid.New().String() 
	}
	return nil
}
