package attendance

import "github.com/gin-gonic/gin"

func (h *Handler) Register(rg *gin.RouterGroup) {
	attendance := rg.Group("/attendance")

	{
		attendance.POST("/:employee_id", h.EmployeeCheckIn)
		attendance.PUT("/:employee_id", h.EmployeeCheckOut)

		attendance.GET("/histories", h.GetDeptAtdHistories)
		attendance.GET("/employee/:employee_id/histories", h.GetEmpAtdHistories)
	}
}
