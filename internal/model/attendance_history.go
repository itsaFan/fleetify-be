package model

import (
	"time"
)

type AttendanceHistory struct {
	ID             uint      `gorm:"primaryKey;autoIncrement;column:id"`
	EmployeeID     string    `gorm:"size:50;not null;column:employee_id"`
	AttendanceID   string    `gorm:"size:100;not null;column:attendance_id"`
	DateAttendance time.Time `gorm:"not null;column:date_attendance"`
	AttendanceType uint8     `gorm:"type:tinyint;not null;column:attendance_type"` //note: 1=In, 2=Out
	Description    string    `gorm:"type:text;column:description"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`

	// Relations
	Employee   Employee   `gorm:"foreignKey:EmployeeID;references:EmployeeID"`
	Attendance Attendance `gorm:"foreignKey:AttendanceID;references:AttendanceID"`
}


