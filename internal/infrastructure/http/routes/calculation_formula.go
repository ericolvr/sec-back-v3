package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type CalculationFormulaRoutes struct {
	handler *api.CalculationFormulaHandler
}

func NewCalculationFormulaRoutes(handler *api.CalculationFormulaHandler) *CalculationFormulaRoutes {
	return &CalculationFormulaRoutes{
		handler: handler,
	}
}

func (r *CalculationFormulaRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	formulas := v1.Group("/calculation-formulas")
	formulas.Use(middleware.PartnerMiddleware())
	{
		formulas.GET("/active", r.handler.GetActive)
		formulas.PUT("/active", r.handler.UpdateActive)
	}
}
