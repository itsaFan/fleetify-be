package employee

import "github.com/gin-gonic/gin"

func (h *Handler) Register(rg *gin.RouterGroup) {
	employee := rg.Group("/employee")

	{
		employee.POST("", h.Create)
		employee.GET("", h.List)
		employee.GET("/:employee_id", h.GetByEmployeeID)
		employee.PATCH("/:employee_id", h.UpdateEmployeeByEmployeeID)
		employee.DELETE("/:employee_id", h.DeleteByEmployeeID)
	}
}
