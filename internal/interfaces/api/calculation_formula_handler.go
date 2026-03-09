package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type CalculationFormulaHandler struct {
	formulaService *services.CalculationFormulaService
}

func NewCalculationFormulaHandler(formulaService *services.CalculationFormulaService) *CalculationFormulaHandler {
	return &CalculationFormulaHandler{
		formulaService: formulaService,
	}
}

type UpdateFormulaRequest struct {
	RiskLowMax               float64 `json:"risk_low_max" binding:"required"`
	RiskMediumMax            float64 `json:"risk_medium_max" binding:"required"`
	ReliabilityAcceptableMin float64 `json:"reliability_acceptable_min" binding:"required"`
	ReliabilityGoodMin       float64 `json:"reliability_good_min" binding:"required"`
	ReliabilityExcellentMin  float64 `json:"reliability_excellent_min" binding:"required"`
	Description              string  `json:"description"`
}

func (h *CalculationFormulaHandler) GetActive(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	formula, err := h.formulaService.GetActive(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get formula", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, formula)
}

func (h *CalculationFormulaHandler) UpdateActive(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	var req UpdateFormulaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	formula := &domain.CalculationFormula{
		PartnerID:                partnerID,
		Version:                  "1.0", // Version não muda no update
		RiskLowMax:               req.RiskLowMax,
		RiskMediumMax:            req.RiskMediumMax,
		ReliabilityAcceptableMin: req.ReliabilityAcceptableMin,
		ReliabilityGoodMin:       req.ReliabilityGoodMin,
		ReliabilityExcellentMin:  req.ReliabilityExcellentMin,
		Description:              req.Description,
	}

	updated, err := h.formulaService.UpdateActive(c.Request.Context(), partnerID, formula)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update formula", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
