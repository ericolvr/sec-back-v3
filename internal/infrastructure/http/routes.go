package http

import (
	"github.com/gin-gonic/gin"
)

type RouteSetup interface {
	SetupRoutes(v1 *gin.RouterGroup)
}

type Router struct {
	// Core routes - mantidas
	AnswerRoutes     RouteSetup
	AuthRoutes       RouteSetup
	DepartmentRoutes RouteSetup
	EmployeeRoutes   RouteSetup
	QuestionRoutes   RouteSetup
	UserRoutes       RouteSetup

	// New routes - refatoradas
	PartnerRoutes                 RouteSetup
	CompanyRoutes                 RouteSetup
	EmployeeSubmissionRoutes      RouteSetup
	AssessmentTemplateRoutes      RouteSetup
	RiskMetricsRoutes             RouteSetup
	AnalyticsRoutes               RouteSetup
	ActionPlanRoutes              RouteSetup
	AssessmentAssignmentRoutes RouteSetup
}

func NewRouter() *Router {
	return &Router{
		// Rotas serão configuradas posteriormente
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	v1 := engine.Group("/api/v1")
	{
		// Rotas básicas
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Rotas serão configuradas quando handlers forem criados
		if r.AuthRoutes != nil {
			r.AuthRoutes.SetupRoutes(v1)
		}
		if r.PartnerRoutes != nil {
			r.PartnerRoutes.SetupRoutes(v1)
		}
		if r.CompanyRoutes != nil {
			r.CompanyRoutes.SetupRoutes(v1)
		}
		if r.DepartmentRoutes != nil {
			r.DepartmentRoutes.SetupRoutes(v1)
		}
		if r.EmployeeRoutes != nil {
			r.EmployeeRoutes.SetupRoutes(v1)
		}
		if r.UserRoutes != nil {
			r.UserRoutes.SetupRoutes(v1)
		}
		if r.QuestionRoutes != nil {
			r.QuestionRoutes.SetupRoutes(v1)
		}
		if r.AnswerRoutes != nil {
			r.AnswerRoutes.SetupRoutes(v1)
		}
		if r.AssessmentTemplateRoutes != nil {
			r.AssessmentTemplateRoutes.SetupRoutes(v1)
		}
		if r.EmployeeSubmissionRoutes != nil {
			r.EmployeeSubmissionRoutes.SetupRoutes(v1)
		}
		if r.RiskMetricsRoutes != nil {
			r.RiskMetricsRoutes.SetupRoutes(v1)
		}
		if r.AnalyticsRoutes != nil {
			r.AnalyticsRoutes.SetupRoutes(v1)
		}
		if r.ActionPlanRoutes != nil {
			r.ActionPlanRoutes.SetupRoutes(v1)
		}
		if r.AssessmentAssignmentRoutes != nil {
			r.AssessmentAssignmentRoutes.SetupRoutes(v1)
		}
	}
}
