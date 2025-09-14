package attendance

import (
	"time"

	"github.com/itsaFan/fleetify-be/internal/helper"
	atdSvc "github.com/itsaFan/fleetify-be/internal/service/attendance"
)

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

type listQueryEmpAtdHistories struct {
	TZ    string `form:"tz" binding:"omitempty"`
	From  string `form:"from" binding:"omitempty"`
	To    string `form:"to" binding:"omitempty"`
	Limit int    `form:"limit"   binding:"omitempty,min=1,max=100"`
	Page  int    `form:"page"    binding:"omitempty,min=1"`
}

type listEmpAtdHistoriesData struct {
	EmployeeID       string                         `json:"employee_id"`
	From             string                         `json:"from"`
	To               string                         `json:"to"`
	TZUsedForRules   string                         `json:"tz_used_for_rules"`
	TZUsedForDisplay string                         `json:"tz_used_for_display"`
	Attendances      []atdSvc.AttendanceHistoryItem `json:"attendances"`
}

type listEmpAtdHistoriesResp struct {
	Message    string                  `json:"message"`
	Data       listEmpAtdHistoriesData `json:"data"`
	Pagination helper.Pagination       `json:"pagination"`
}

type listQueryDeptAtdHistories struct {
	Department *uint64 `form:"dept_id" binding:"omitempty"`
	TZ         string  `form:"tz" binding:"omitempty"`
	From       string  `form:"from" binding:"omitempty"`
	To         string  `form:"to" binding:"omitempty"`
	Limit      int     `form:"limit"   binding:"omitempty,min=1,max=100"`
	Page       int     `form:"page"    binding:"omitempty,min=1"`
}

type listDeptAtdHistoriesData struct {
	DepartmentID     *uint64                        `json:"department_id,omitempty"`
	From             string                         `json:"from"`
	To               string                         `json:"to"`
	TZUsedForRules   string                         `json:"tz_used_for_rules"`
	TZUsedForDisplay string                         `json:"tz_used_for_display"`
	Attendances      []atdSvc.AttendanceHistoryItem `json:"attendances"`
}

type listDeptAtdHistoriesResp struct {
	Message    string                   `json:"message"`
	Data       listDeptAtdHistoriesData `json:"data"`
	Pagination helper.Pagination        `json:"pagination"`
}
