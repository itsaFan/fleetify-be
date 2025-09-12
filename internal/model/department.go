package model

type Department struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement;column:id"`
	DepartmentName  string `gorm:"size:255;not null;column:department_name"`
	MaxClockInTime  string `gorm:"type:time;not null;column:max_clock_in_time"` 
	MaxClockOutTime string `gorm:"type:time;not null;column:max_clock_out_time"`

	Employees []Employee `gorm:"foreignKey:DepartmentID;references:ID"`
}
