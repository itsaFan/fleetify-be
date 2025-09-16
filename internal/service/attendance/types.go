package attendance

import "time"

type ListInputEmp struct {
	EmployeeID string
	FromLocal  string
	ToLocal    string
	TZ         string
	Limit      int
	Page       int
}

type ListInputDept struct {
	DepartmentID *uint64
	FromLocal    string
	ToLocal      string
	TZ           string
	Limit        int
	Page         int
}

type AttendanceHistoryItem struct {
	EmployeeID      string     `json:"employee_id"`
	EmployeeName    string     `json:"employee_name,omitempty"`
	DepartmentName  *string    `json:"department_name,omitempty"`
	DateLocal       string     `json:"date_local"`
	ClockInLocal    *string    `json:"clock_in_local"`
	ClockInUTC      *time.Time `json:"clock_in_utc"`
	StatusIn        string     `json:"status_in"` // on_time | late | early | missing_in
	DeltaInMinutes  *int       `json:"delta_in_minutes"`
	ClockOutLocal   *string    `json:"clock_out_local"`
	ClockOutUTC     *time.Time `json:"clock_out_utc"`
	StatusOut       string     `json:"status_out"` // normal | overtime | no_out
	DeltaOutMinutes *int       `json:"delta_out_minutes"`
	AttendanceID    string     `json:"attendance_id,omitempty"`
}

type AttendanceHistoryOutput struct {
	Items     []AttendanceHistoryItem
	Total     int64
	FromLocal string
	ToLocal   string
	TZUsed    string
}
