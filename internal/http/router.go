package http

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dpthttp "github.com/itsaFan/fleetify-be/internal/http/department"
	deptrepo "github.com/itsaFan/fleetify-be/internal/repo/department"
	deptsvc "github.com/itsaFan/fleetify-be/internal/service/department"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	v1 := r.Group("/v1")

	dptRepo := deptrepo.New(db)
	dptSvc := deptsvc.New(dptRepo)
	dptHdl := dpthttp.New(dptSvc)
	dptHdl.Register(v1)

	return r
}
