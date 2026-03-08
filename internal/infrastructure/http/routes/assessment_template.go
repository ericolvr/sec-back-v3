package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type AssessmentTemplateRoutes struct {
	templateHandler *api.AssessmentTemplateHandler
	versionHandler  *api.AssessmentVersionHandler
	questionHandler *api.QuestionHandler
}

func NewAssessmentTemplateRoutes(templateHandler *api.AssessmentTemplateHandler, versionHandler *api.AssessmentVersionHandler, questionHandler *api.QuestionHandler) *AssessmentTemplateRoutes {
	return &AssessmentTemplateRoutes{
		templateHandler: templateHandler,
		versionHandler:  versionHandler,
		questionHandler: questionHandler,
	}
}

func (r *AssessmentTemplateRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	templates := v1.Group("/assessments")
	templates.Use(middleware.PartnerMiddleware())
	templates.Use(middleware.JWTMiddleware())
	{
		templates.POST("", r.templateHandler.Create)
		templates.GET("", r.templateHandler.List)
		templates.GET("/deleted", r.templateHandler.ListDeleted)
		templates.GET("/:id/versions", r.versionHandler.ListByTemplate)
		templates.GET("/:id/questions", r.questionHandler.List)
		templates.GET("/:id", r.templateHandler.GetByID)
		templates.PUT("/:id", r.templateHandler.Update)
		templates.PATCH("/:id/toggle-active", r.templateHandler.ToggleActive)
		templates.DELETE("/:id", r.templateHandler.Delete)
	}
}
