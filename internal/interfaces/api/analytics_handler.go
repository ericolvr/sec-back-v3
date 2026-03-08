package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// CreateSnapshot cria um snapshot de analytics
func (h *AnalyticsHandler) CreateSnapshot(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	// Buscar department_id da URL
	departmentID, err := strconv.ParseInt(c.Param("department_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department_id"})
		return
	}

	// Buscar template_id do query parameter (compatibilidade com frontend)
	templateIDStr := c.Query("questionnaire_id")
	if templateIDStr == "" {
		templateIDStr = c.Query("template_id")
	}
	if templateIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "questionnaire_id or template_id is required"})
		return
	}

	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template_id"})
		return
	}

	// Buscar company_id do departamento
	// TODO: implementar busca real do company_id
	companyID := int64(1)

	snapshot, err := h.analyticsService.CreateSnapshot(
		c.Request.Context(),
		partnerID,
		companyID,
		departmentID,
		templateID,
		nil, // created_by
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, snapshot)
}

// GetSnapshot busca um snapshot específico
func (h *AnalyticsHandler) GetSnapshot(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	snapshotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snapshot ID"})
		return
	}

	snapshot, err := h.analyticsService.GetReport(c.Request.Context(), partnerID, snapshotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Snapshot not found"})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// ListReportsByDepartment lista snapshots de um departamento específico
func (h *AnalyticsHandler) ListReportsByDepartment(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	departmentID := c.Param("department_id")

	deptID, err := strconv.ParseInt(departmentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department_id must be a valid number"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	reports, err := h.analyticsService.ListReportsByDepartment(c.Request.Context(), partnerID, deptID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// ListReportsByTemplate lista snapshots de um template específico
func (h *AnalyticsHandler) ListReportsByTemplate(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	templateID := c.Param("template_id")

	qID, err := strconv.ParseInt(templateID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id must be a valid number"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	reports, err := h.analyticsService.ListReportsByTemplate(c.Request.Context(), partnerID, qID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// ListReports lista todos os snapshots
func (h *AnalyticsHandler) ListReports(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	reports, err := h.analyticsService.ListReports(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// ListSnapshots alias para ListReports
func (h *AnalyticsHandler) ListSnapshots(c *gin.Context) {
	h.ListReports(c)
}

// GetDepartmentReport retorna analytics em tempo real de um departamento
func (h *AnalyticsHandler) GetDepartmentReport(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	departmentID := c.Param("department_id")
	templateIDStr := c.Query("template_id")

	if templateIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id is required"})
		return
	}

	deptID, err := strconv.ParseInt(departmentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "department_id must be a valid number"})
		return
	}

	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id must be a valid number"})
		return
	}

	// Buscar companyID do departamento
	companyID := int64(1) // TODO: buscar do departamento

	report, err := h.analyticsService.GetDepartmentReport(c.Request.Context(), partnerID, companyID, deptID, templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetTemplateReport retorna analytics de todos os departamentos de um template
func (h *AnalyticsHandler) GetTemplateReport(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	templateID := c.Param("template_id")

	qID, err := strconv.ParseInt(templateID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "template_id must be a valid number"})
		return
	}

	// TODO: Implementar GetTemplateReport no service
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":       "Not implemented yet",
		"partner_id":  partnerID,
		"template_id": qID,
		"message":     "This endpoint will return analytics for all departments of a questionnaire",
	})
}

// GetRiskCategoriesBySnapshot retorna risk categories de um snapshot
func (h *AnalyticsHandler) GetRiskCategoriesBySnapshot(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	snapshotID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid snapshot ID"})
		return
	}

	categories, err := h.analyticsService.GetRiskCategoriesBySnapshot(c.Request.Context(), partnerID, snapshotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetClosedDepartmentSnapshots retorna todos os snapshots de departamentos encerrados
func (h *AnalyticsHandler) GetClosedDepartmentSnapshots(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	snapshots, err := h.analyticsService.GetClosedDepartmentSnapshots(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapshots)
}

// GetInProgressQuestionnaires retorna todos os questionários em andamento
func (h *AnalyticsHandler) GetInProgressQuestionnaires(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	// TODO: Buscar company_id do usuário logado
	// Por enquanto usa company_id = 1
	companyID := int64(1)

	questionnaires, err := h.analyticsService.GetInProgressTemplates(c.Request.Context(), partnerID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questionnaires": questionnaires,
	})
}
