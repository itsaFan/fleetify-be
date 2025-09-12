package employee

import (
	"time"

	empsvc "github.com/itsaFan/fleetify-be/internal/service/employee"
)

type departmentResp struct {
	ID              uint64 `json:"id"`
	DepartmentName  string `json:"department_name"`
	MaxClockInTime  string `json:"max_clock_in_time"`
	MaxClockOutTime string `json:"max_clock_out_time"`
}

type employeeResp struct {
	ID         uint64         `json:"id"`
	EmployeeID string         `json:"employee_id"`
	Name       string         `json:"name"`
	Address    string         `json:"address"`
	Department departmentResp `json:"department"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type listQuery struct {
	Search  string `form:"search"`
	Limit   int    `form:"limit"   binding:"omitempty,min=1,max=100"`
	Page    int    `form:"page"    binding:"omitempty,min=1"`
	SortBy  string `form:"sortBy"  binding:"omitempty,oneof=id department_name"`
	SortDir string `form:"sortDir" binding:"omitempty,oneof=asc desc"`
}

type listResponse struct {
	Message    string            `json:"message"`
	Data       []employeeResp    `json:"data"`
	Pagination empsvc.Pagination `json:"pagination"`
}

type createReq struct {
	Name       string  `json:"name" binding:"required,max=255"`
	Address    *string `json:"address"`
	Department uint64  `json:"department" binding:"required"`
}

type createResponse struct {
	Message string       `json:"message"`
	Data    employeeResp `json:"data"`
}

type getByEmployeeIDResponse struct {
	Message string       `json:"message"`
	Data    employeeResp `json:"data"`
}

type updateReq struct {
	Name       *string `json:"name,omitempty"`
	Address    *string `json:"address,omitempty"`
	Department *uint64 `json:"department,omitempty"`
}

type updateResponse struct {
	Message string       `json:"message"`
	Data    employeeResp `json:"data"`
}
type deleteResponse struct {
	Message string `json:"message"`
}
