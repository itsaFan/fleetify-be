package model

import (
	"time"
)

type Attendance struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;column:id"`
	EmployeeID   string     `gorm:"size:50;not null;column:employee_id"`
	AttendanceID string     `gorm:"size:100;not null;column:attendance_id"`
	ClockIn      *time.Time `gorm:"column:clock_in"`
	ClockOut     *time.Time `gorm:"column:clock_out"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`

	// Relations
	Employee Employee            `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	History  []AttendanceHistory `gorm:"foreignKey:AttendanceID;references:AttendanceID"`
}


