package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type DashboardRoutes struct {
	handler *api.DashboardHandler
}

func NewDashboardRoutes(handler *api.DashboardHandler) *DashboardRoutes {
	return &DashboardRoutes{
		handler: handler,
	}
}

func (r *DashboardRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	dashboard := v1.Group("/dashboard")
	dashboard.Use(middleware.PartnerMiddleware())
	{
		// Dashboard do Partner (consultoria)
		dashboard.GET("/partner", r.handler.GetPartnerDashboard)
		
		// Dashboard da Company (gestor)
		dashboard.GET("/company/:companyId", r.handler.GetCompanyDashboard)
		
		// Dashboard do Department (supervisor)
		dashboard.GET("/company/:companyId/department/:departmentId", r.handler.GetDepartmentDashboard)
	}
}
