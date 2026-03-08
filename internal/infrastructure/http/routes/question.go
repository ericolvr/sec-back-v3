package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type QuestionRoutes struct {
	questionHandler *api.QuestionHandler
}

func NewQuestionRoutes(questionHandler *api.QuestionHandler) *QuestionRoutes {
	return &QuestionRoutes{
		questionHandler: questionHandler,
	}
}

func (qr *QuestionRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	questions := v1.Group("/questions")
	questions.Use(middleware.PartnerMiddleware())
	// questions.Use(middleware.JWTMiddleware())
	{
		questions.POST("", qr.questionHandler.Create)
		questions.GET("", qr.questionHandler.ListAll)
		questions.GET("/:id", qr.questionHandler.GetByID)
		questions.PUT("/:id", qr.questionHandler.Update)
		questions.DELETE("/:id", qr.questionHandler.Delete)
	}
}
