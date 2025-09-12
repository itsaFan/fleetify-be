package employee

import (
	"errors"
	stdhttp "net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	empSvc "github.com/itsaFan/fleetify-be/internal/service/employee"
)

type Handler struct {
	svc empSvc.Service
}

func New(svc empSvc.Service) *Handler {
	return &Handler{svc: svc}
}

// POST
func (h *Handler) Create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequest(c, "invalid JSON body")
		return
	}

	input := empSvc.CreateInput{
		Name:       helper.NormalizeStringField(req.Name),
		Address:    req.Address,
		Department: req.Department,
	}

	emp, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, appErr.ErrNotFound) {
			helper.NotFound(c, "department not found")
			return
		}
		helper.WriteError(c, err)
		return
	}

	data := employeeResp{
		ID:         emp.ID,
		EmployeeID: emp.EmployeeID,
		Name:       emp.Name,
		Address:    emp.Address,
		Department: departmentResp{
			ID:              emp.Department.ID,
			DepartmentName:  emp.Department.DepartmentName,
			MaxClockInTime:  emp.Department.MaxClockInTime,
			MaxClockOutTime: emp.Department.MaxClockOutTime,
		},
		CreatedAt: emp.CreatedAt,
		UpdatedAt: emp.UpdatedAt,
	}
	c.JSON(stdhttp.StatusCreated, createResponse{
		Message: "Employee created successfully",
		Data:    data,
	})

}

// GET List with flex q
func (h *Handler) List(c *gin.Context) {
	var q listQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.BadRequest(c, "invalid query parameters")
		return
	}
	out, err := h.svc.List(c.Request.Context(), empSvc.ListInput{
		Search:  q.Search,
		Limit:   q.Limit,
		Page:    q.Page,
		SortBy:  q.SortBy,
		SortDir: q.SortDir,
	})
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := make([]employeeResp, 0, len(out.Data))
	for _, e := range out.Data {
		data = append(data, employeeResp{
			ID:         e.ID,
			EmployeeID: e.EmployeeID,
			Name:       e.Name,
			Address:    e.Address,
			Department: departmentResp{
				ID:              e.Department.ID,
				DepartmentName:  e.Department.DepartmentName,
				MaxClockInTime:  e.Department.MaxClockInTime,
				MaxClockOutTime: e.Department.MaxClockOutTime,
			},
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	c.JSON(stdhttp.StatusOK, listResponse{
		Message:    "Employees retrieved successfully",
		Data:       data,
		Pagination: out.Pagination,
	})
}

// GET by EmployeeID
func (h *Handler) GetByEmployeeID(c *gin.Context) {
	raw := c.Param("employee_id")

	empId, err := url.PathUnescape(raw)
	if err != nil {
		helper.BadRequest(c, "Invalid employee_id name in path")
		return
	}

	emp, err := h.svc.GetByEmployeeID(c.Request.Context(), empId)
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := employeeResp{
		ID:         emp.ID,
		EmployeeID: emp.EmployeeID,
		Name:       emp.Name,
		Address:    emp.Address,
		Department: departmentResp{
			ID:              emp.Department.ID,
			DepartmentName:  emp.Department.DepartmentName,
			MaxClockInTime:  emp.Department.MaxClockInTime,
			MaxClockOutTime: emp.Department.MaxClockOutTime,
		},
		CreatedAt: emp.CreatedAt,
		UpdatedAt: emp.UpdatedAt,
	}
	c.JSON(stdhttp.StatusOK, getByEmployeeIDResponse{
		Message: "Employee retrieved successfully",
		Data:    data,
	})
}

func (h *Handler) UpdateEmployeeByEmployeeID(c *gin.Context) {
	raw := c.Param("employee_id")
	empId, err := url.PathUnescape(raw)
	if err != nil {
		helper.BadRequest(c, "Invalid employee_id name in path")
		return
	}

	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequest(c, "invalid JSON body")
		return
	}

	in := empSvc.UpdateInput{}
	if req.Name != nil {
		n := helper.NormalizeStringField(*req.Name)
		in.Name = n
	}

	if req.Address != nil {
		n := helper.NormalizeStringField(*req.Address)
		in.Address = &n
	}
	if req.Department != nil {
		in.Department = *req.Department
	}

	emp, err := h.svc.UpdateEmployeeByEmployeeID(c.Request.Context(), empId, in)
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := employeeResp{
		ID:         emp.ID,
		EmployeeID: emp.EmployeeID,
		Name:       emp.Name,
		Address:    emp.Address,
		Department: departmentResp{
			ID:              emp.Department.ID,
			DepartmentName:  emp.Department.DepartmentName,
			MaxClockInTime:  emp.Department.MaxClockInTime,
			MaxClockOutTime: emp.Department.MaxClockOutTime,
		},
		CreatedAt: emp.CreatedAt,
		UpdatedAt: emp.UpdatedAt,
	}

	c.JSON(stdhttp.StatusOK, updateResponse{
		Message: "Employee updated successfully",
		Data:    data,
	})

}

func (h *Handler) DeleteByEmployeeID(c *gin.Context) {
	raw := c.Param("employee_id")
	name, err := url.PathUnescape(raw)
	if err != nil {
		helper.BadRequest(c, "invalid employee_id name in path")
		return
	}

	if err := h.svc.DeleteByEmployeeID(c.Request.Context(), name); err != nil {
		helper.WriteError(c, err)
		return
	}

	c.JSON(stdhttp.StatusOK, deleteResponse{
		Message: "Employee deleted successfully",
	})
}
