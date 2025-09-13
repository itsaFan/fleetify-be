package attendance

import "time"

type attendanceData struct {
	AttendanceID string     `json:"attendance_id"`
	EmployeeID   string     `json:"employee_id"`
	ClockIn      *time.Time `json:"clock_in,omitempty"`
	ClockOut     *time.Time `json:"clock_out,omitempty"`
}

type checkInResponse struct {
	Message string         `json:"message"`
	Data    attendanceData `json:"data"`
}

type checkOutResponse struct {
	Message string         `json:"message"`
	Data    attendanceData `json:"data"`
}
