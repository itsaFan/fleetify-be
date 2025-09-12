package department

import deptsvc "github.com/itsaFan/fleetify-be/internal/service/department"

type departmentResp struct {
	ID             uint64 `json:"id"`
	DepartmentName string `json:"department_name"`
	MaxClockIn     string `json:"max_clock_in"`
	MaxClockOut    string `json:"max_clock_out"`
}

type listQuery struct {
	Search  string `form:"search"`
	Limit   int    `form:"limit"   binding:"omitempty,min=1,max=100"`
	Page    int    `form:"page"    binding:"omitempty,min=1"`
	SortBy  string `form:"sortBy"  binding:"omitempty,oneof=id department_name"`
	SortDir string `form:"sortDir" binding:"omitempty,oneof=asc desc"`
}

type createReq struct {
	DepartmentName string `json:"department_name" binding:"required,max=255"`
	MaxClockIn     string `json:"max_clock_in"   binding:"required"`
	MaxClockOut    string `json:"max_clock_out"  binding:"required"`
}
type createResponse struct {
	Message string         `json:"message"`
	Data    departmentResp `json:"data"`
}

type listResponse struct {
	Message    string             `json:"message"`
	Data       []departmentResp   `json:"data"`
	Pagination deptsvc.Pagination `json:"pagination"`
}

type getByNameResponse struct {
	Message string         `json:"message"`
	Data    departmentResp `json:"data"`
}

type updateReq struct {
	DepartmentName *string `json:"department_name,omitempty"`
	MaxClockIn     *string `json:"max_clock_in,omitempty"`
	MaxClockOut    *string `json:"max_clock_out,omitempty"`
}

type updateResponse struct {
	Message string         `json:"message"`
	Data    departmentResp `json:"data"`
}

type deleteResponse struct {
	Message string `json:"message"`
}
