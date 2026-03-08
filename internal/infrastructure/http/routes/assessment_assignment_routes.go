package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type AssessmentAssignmentRoutes struct {
	assignmentHandler *api.AssessmentAssignmentHandler
}

func NewAssessmentAssignmentRoutes(assignmentHandler *api.AssessmentAssignmentHandler) *AssessmentAssignmentRoutes {
	return &AssessmentAssignmentRoutes{
		assignmentHandler: assignmentHandler,
	}
}

func (r *AssessmentAssignmentRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	assignments := v1.Group("/assessment-assignments")
	assignments.Use(middleware.PartnerMiddleware())
	{
		assignments.POST("", r.assignmentHandler.Create)
		assignments.GET("", r.assignmentHandler.List)
		assignments.GET("/:id", r.assignmentHandler.GetByID)
		assignments.POST("/:id/close", r.assignmentHandler.Close)
		assignments.DELETE("/:id", r.assignmentHandler.Delete)
	}
}
