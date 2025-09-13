package http

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dpthttp "github.com/itsaFan/fleetify-be/internal/http/department"
	deptrepo "github.com/itsaFan/fleetify-be/internal/repo/department"
	deptsvc "github.com/itsaFan/fleetify-be/internal/service/department"

	emphttp "github.com/itsaFan/fleetify-be/internal/http/employee"
	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"
	empsvc "github.com/itsaFan/fleetify-be/internal/service/employee"

	atdhttp "github.com/itsaFan/fleetify-be/internal/http/attendance"
	atdrepo "github.com/itsaFan/fleetify-be/internal/repo/attendance"
	atdsvc "github.com/itsaFan/fleetify-be/internal/service/attendance"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	v1 := r.Group("/v1")

	dptRepo := deptrepo.New(db)
	dptSvc := deptsvc.New(dptRepo)
	dptHdl := dpthttp.New(dptSvc)
	dptHdl.Register(v1)

	empRepo := emprepo.New(db)
	empSvc := empsvc.New(empRepo, dptRepo)
	empHdl := emphttp.New(empSvc)
	empHdl.Register(v1)

	atdRepo := atdrepo.New(db)
	atdSvc := atdsvc.New(atdRepo, empRepo)
	atdHdl := atdhttp.New(atdSvc)
	atdHdl.Register(v1)

	return r
}
