package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type DepartmentRoutes struct {
	departmentHandler *api.DepartmentHandler
}

func NewDepartmentRoutes(departmentHandler *api.DepartmentHandler) *DepartmentRoutes {
	return &DepartmentRoutes{
		departmentHandler: departmentHandler,
	}
}

func (dr *DepartmentRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	departments := v1.Group("/departments")
	departments.Use(middleware.PartnerMiddleware())
	// departments.Use(middleware.JWTMiddleware())
	{
		departments.POST("", dr.departmentHandler.Create)
		departments.GET("", dr.departmentHandler.List)
		departments.GET("/deleted", dr.departmentHandler.ListDeleted)
		departments.GET("/:id", dr.departmentHandler.GetByID)
		departments.PATCH("/:id", dr.departmentHandler.Update)
		departments.DELETE("/:id", dr.departmentHandler.Delete)
		departments.PATCH("/:id/toggle-active", dr.departmentHandler.ToggleActive)
	}
}
