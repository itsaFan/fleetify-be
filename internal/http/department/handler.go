package department

import (
	stdhttp "net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/itsaFan/fleetify-be/internal/helper"
	deptSvc "github.com/itsaFan/fleetify-be/internal/service/department"
)

type Handler struct {
	svc deptSvc.Service
}

func New(svc deptSvc.Service) *Handler {
	return &Handler{svc: svc}
}

// POST
func (h *Handler) Create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequest(c, "invalid JSON body")
		return
	}
	input := deptSvc.CreateInput{
		DepartmentName: helper.NormalizeStringField(req.DepartmentName),
		MaxClockIn:     req.MaxClockIn,
		MaxClockOut:    req.MaxClockOut,
	}

	dept, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		helper.WriteError(c, err)
		return
	}
	data := departmentResp{
		DepartmentName: dept.DepartmentName,
		MaxClockIn:     dept.MaxClockInTime,
		MaxClockOut:    dept.MaxClockOutTime,
	}

	c.JSON(stdhttp.StatusCreated, createResponse{
		Message: "Department created successfully",
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

	out, err := h.svc.List(c.Request.Context(), deptSvc.ListInput{
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

	data := make([]departmentResp, len(out.Data))
	for i := range out.Data {
		d := out.Data[i]
		data[i] = departmentResp{
			DepartmentName: d.DepartmentName,
			MaxClockIn:     d.MaxClockInTime,
			MaxClockOut:    d.MaxClockOutTime,
		}
	}

	c.JSON(stdhttp.StatusOK, listResponse{
		Message:    "Department retrieved successfully",
		Data:       data,
		Pagination: out.Pagination,
	})
}

// GET by Name
func (h *Handler) GetByName(c *gin.Context) {
	raw := c.Param("name")

	name, err := url.PathUnescape(raw)
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "bad_request", "message": "invalid department name in path"})
		return
	}

	dept, err := h.svc.GetByName(c.Request.Context(), name)
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	data := departmentResp{
		DepartmentName: dept.DepartmentName,
		MaxClockIn:     dept.MaxClockInTime,
		MaxClockOut:    dept.MaxClockOutTime,
	}
	c.JSON(stdhttp.StatusOK, getByNameResponse{
		Message: "Department retrieved successfully",
		Data:    data,
	})
}

// Update dpt
func (h *Handler) UpdateByName(c *gin.Context) {
	raw := c.Param("name")

	name, err := url.PathUnescape(raw)

	if err != nil {
		helper.BadRequest(c, "Invalid department name in path")
		return
	}

	var req updateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequest(c, "invalid JSON body")
		return
	}

	in := deptSvc.UpdateInput{}
	if req.DepartmentName != nil {
		n := helper.NormalizeStringField(*req.DepartmentName)
		in.DepartmentName = &n
	}
	if req.MaxClockIn != nil {
		in.MaxClockIn = req.MaxClockIn
	}
	if req.MaxClockOut != nil {
		in.MaxClockOut = req.MaxClockOut
	}

	dept, err := h.svc.UpdateByName(c.Request.Context(), name, in)
	if err != nil {
		helper.WriteError(c, err)
		return
	}

	c.JSON(stdhttp.StatusOK, updateResponse{
		Message: "Department updated successfully",
		Data: departmentResp{
			DepartmentName: dept.DepartmentName,
			MaxClockIn:     dept.MaxClockInTime,
			MaxClockOut:    dept.MaxClockOutTime,
		},
	})

}

func (h *Handler) DeleteByName(c *gin.Context) {
	raw := c.Param("name")
	name, err := url.PathUnescape(raw)
	if err != nil {
		helper.BadRequest(c, "invalid department name in path")
		return
	}

	if err := h.svc.DeleteByName(c.Request.Context(), name); err != nil {
		helper.WriteError(c, err)
		return
	}

	c.JSON(stdhttp.StatusOK, deleteResponse{
		Message: "Department deleted successfully",
	})
}
