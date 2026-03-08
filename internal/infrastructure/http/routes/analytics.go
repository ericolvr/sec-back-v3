package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type AnalyticsRoutes struct {
	analyticsHandler *api.AnalyticsHandler
}

func NewAnalyticsRoutes(analyticsHandler *api.AnalyticsHandler) *AnalyticsRoutes {
	return &AnalyticsRoutes{
		analyticsHandler: analyticsHandler,
	}
}

func (r *AnalyticsRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	analytics := v1.Group("/analytics")
	analytics.Use(middleware.PartnerMiddleware())
	{
		// Real-time analytics
		analytics.GET("/questionnaire/:template_id", r.analyticsHandler.GetTemplateReport)
		analytics.GET("/department/:department_id", r.analyticsHandler.GetDepartmentReport)

		// Snapshot/Reports (department-based)
		analytics.POST("/department/:department_id/snapshot", r.analyticsHandler.CreateSnapshot)
		analytics.GET("/department/:department_id/reports", r.analyticsHandler.ListReportsByDepartment)
		analytics.GET("/questionnaire/:template_id/reports", r.analyticsHandler.ListReportsByTemplate)
		analytics.GET("/reports", r.analyticsHandler.ListReports)
		analytics.GET("/snapshots", r.analyticsHandler.ListSnapshots)
		analytics.GET("/snapshots/:id", r.analyticsHandler.GetSnapshot)
		analytics.GET("/snapshots/:id/risk-categories", r.analyticsHandler.GetRiskCategoriesBySnapshot)
	}

	// Statistics endpoints (compatibilidade com frontend)
	statistics := v1.Group("/statistics")
	statistics.Use(middleware.PartnerMiddleware())
	{
		statistics.GET("/questionnaires/in-progress", r.analyticsHandler.GetInProgressQuestionnaires)
		statistics.GET("/questionnaires/closed/risk-overview", r.analyticsHandler.GetClosedDepartmentSnapshots)
	}
}
