package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetCompanyDashboard retorna dashboard para gestor de empresa
func (h *DashboardHandler) GetCompanyDashboard(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	dashboard, err := h.dashboardService.GetCompanyDashboard(c.Request.Context(), partnerID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get company dashboard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetPartnerDashboard retorna dashboard para consultoria (Partner)
func (h *DashboardHandler) GetPartnerDashboard(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	dashboard, err := h.dashboardService.GetPartnerDashboard(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get partner dashboard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetDepartmentDashboard retorna dashboard para supervisor de departamento
func (h *DashboardHandler) GetDepartmentDashboard(c *gin.Context) {
	partnerID, err := strconv.ParseInt(c.GetHeader("X-Partner-ID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partner ID"})
		return
	}

	companyID, err := strconv.ParseInt(c.Param("companyId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	departmentID, err := strconv.ParseInt(c.Param("departmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	dashboard, err := h.dashboardService.GetDepartmentDashboard(c.Request.Context(), partnerID, companyID, departmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get department dashboard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
