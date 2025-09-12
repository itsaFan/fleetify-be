package model

import "time"

type Department struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id"`
	DepartmentName  string    `gorm:"size:255;not null;column:department_name"`
	MaxClockInTime  time.Time `gorm:"type:time;column:max_clock_in_time"`
	MaxClockOutTime time.Time `gorm:"type:time;column:max_clock_out_time"`

	// Relations
	Employees []Employee `gorm:"foreignKey:DepartmentID;references:ID"`
}
