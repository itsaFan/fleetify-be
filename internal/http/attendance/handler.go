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
