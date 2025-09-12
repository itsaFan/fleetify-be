package department

import "github.com/gin-gonic/gin"

func (h *Handler) Register(rg *gin.RouterGroup) {
	departments := rg.Group("/departments")

	{
		departments.POST("", h.Create)
		departments.GET("", h.List)
		departments.GET("/:name", h.GetByName)
		departments.PATCH("/:name", h.UpdateByName)
		departments.DELETE("/:name", h.DeleteByName)
	}
}
