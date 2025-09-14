package attendance

import (
	stdhttp "net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/itsaFan/fleetify-be/internal/helper"
	atdSvc "github.com/itsaFan/fleetify-be/internal/service/attendance"
)

type Handler struct {
	svc atdSvc.Service
}

func New(svc atdSvc.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) EmployeeCheckIn(c *gin.Context) {
	raw := c.Param("employee_id")

	empId, err := url.PathUnescape(raw)

	if err != nil {
		helper.BadRequest(c, "Invalid employee_id name in path")
		return
	}

	atd, err := h.svc.CreateEmpAttendance(c.Request.Context(), empId)

	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := attendanceData{
		AttendanceID: atd.AttendanceID,
		EmployeeID:   atd.EmployeeID,
		ClockIn:      atd.ClockIn,
		ClockOut:     atd.ClockOut,
	}

	c.JSON(stdhttp.StatusCreated, checkInResponse{
		Message: "Attendance: Clock In success",
		Data:    data,
	})
}

func (h *Handler) EmployeeCheckOut(c *gin.Context) {
	raw := c.Param("employee_id")

	empId, err := url.PathUnescape(raw)

	if err != nil {
		helper.BadRequest(c, "Invalid employee_id name in path")
		return
	}

	atd, err := h.svc.CloseEmpAttendance(c.Request.Context(), empId)

	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := attendanceData{
		AttendanceID: atd.AttendanceID,
		EmployeeID:   atd.EmployeeID,
		ClockIn:      atd.ClockIn,
		ClockOut:     atd.ClockOut,
	}

	c.JSON(stdhttp.StatusCreated, checkOutResponse{
		Message: "Attendance: Clock Out success",
		Data:    data,
	})

}

func (h *Handler) GetEmpAtdHistories(c *gin.Context) {
	empId := c.Param("employee_id")

	var q listQueryEmpAtdHistories
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.BadRequest(c, "Invalid query parameters")
	}

	if q.Limit == 0 {
		q.Limit = 10
	}
	if q.Page == 0 {
		q.Page = 1
	}
	if q.TZ == "" {
		q.TZ = "UTC"
	}

	res, err := h.svc.ListEmployeeAtdHistories(c.Request.Context(), atdSvc.ListInputEmp{
		EmployeeID: empId,
		FromLocal:  q.From,
		ToLocal:    q.To,
		TZ:         q.TZ,
		Limit:      q.Limit,
		Page:       q.Page,
	})
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	resp := listEmpAtdHistoriesResp{
		Message: "Employee attendances retrieved successfully",
		Data: listEmpAtdHistoriesData{
			EmployeeID:       empId,
			From:             res.FromLocal,
			To:               res.ToLocal,
			TZUsedForRules:   res.TZUsed,
			TZUsedForDisplay: res.TZUsed,
			Attendances:      res.Items,
		},
		Pagination: helper.BuildPagination(res.Total, q.Page, q.Limit),
	}

	c.JSON(stdhttp.StatusOK, resp)
}

func (h *Handler) GetDeptAtdHistories(c *gin.Context) {

	var q listQueryDeptAtdHistories
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.BadRequest(c, "Invalid query parameters")
	}

	if q.Limit == 0 {
		q.Limit = 10
	}
	if q.Page == 0 {
		q.Page = 1
	}
	if q.TZ == "" {
		q.TZ = "UTC"
	}

	res, err := h.svc.ListDeparmentAtdHistories(c.Request.Context(), atdSvc.ListInputDept{
		DepartmentID: q.Department,
		FromLocal:    q.From,
		ToLocal:      q.To,
		TZ:           q.TZ,
		Limit:        q.Limit,
		Page:         q.Page,
	})

	if err != nil {
		helper.WriteError(c, err)
		return
	}

	resp := listDeptAtdHistoriesResp{
		Message: "Department attendance logs retrieved successfully",
		Data: listDeptAtdHistoriesData{
			DepartmentID:     q.Department,
			From:             res.FromLocal,
			To:               res.ToLocal,
			TZUsedForRules:   res.TZUsed,
			TZUsedForDisplay: res.TZUsed,
			Attendances:      res.Items,
		},
		Pagination: helper.BuildPagination(res.Total, q.Page, q.Limit),
	}

	c.JSON(stdhttp.StatusOK, resp)

}
