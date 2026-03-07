package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type AnswerRoutes struct {
	answerHandler *api.AnswerHandler
}

func NewAnswerRoutes(answerHandler *api.AnswerHandler) *AnswerRoutes {
	return &AnswerRoutes{
		answerHandler: answerHandler,
	}
}

func (ar *AnswerRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	// Rota pública para envio em lote (sem middleware)
	v1.POST("/answers/batch", ar.answerHandler.SubmitBatch)

	// Rotas autenticadas
	answers := v1.Group("/answers")
	answers.Use(middleware.TenantMiddleware())
	// answers.Use(middleware.JWTMiddleware())
	{
		answers.POST("", ar.answerHandler.Create)
		answers.GET("", ar.answerHandler.List) // requires ?response_id=X
		answers.GET("/:id", ar.answerHandler.GetByID)
		answers.PUT("/:id", ar.answerHandler.Update)
	}
}
