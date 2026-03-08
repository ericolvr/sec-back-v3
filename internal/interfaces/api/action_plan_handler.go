package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type ActionPlanHandler struct {
	actionPlanService *services.ActionPlanService
}

func NewActionPlanHandler(actionPlanService *services.ActionPlanService) *ActionPlanHandler {
	return &ActionPlanHandler{
		actionPlanService: actionPlanService,
	}
}

func (h *ActionPlanHandler) Create(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	var req struct {
		CompanyID       int64    `json:"company_id" binding:"required"`
		TemplateID int64    `json:"template_id" binding:"required"`
		DepartmentID    int64    `json:"department_id" binding:"required"`
		SnapshotID      *int64   `json:"snapshot_id"`
		Title           string   `json:"title" binding:"required"`
		Description     string   `json:"description" binding:"required"`
		RiskLevel       string   `json:"risk_level" binding:"required"`
		Priority        string   `json:"priority" binding:"required"`
		Category        string   `json:"category"`
		ResponsibleName string   `json:"responsible_name" binding:"required"`
		ResponsibleID   *int64   `json:"responsible_id"`
		Status          string   `json:"status"`
		DueDate         *string  `json:"due_date"`
		EvidenceURLs    []string `json:"evidence_urls"`
		Notes           string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use YYYY-MM-DD"})
			return
		}
		dueDate = &parsed
	}

	if req.Status == "" {
		req.Status = "pending"
	}

	actionPlan := &domain.ActionPlan{
		PartnerID:       partnerID,
		CompanyID:       req.CompanyID,
		TemplateID: req.TemplateID,
		DepartmentID:    req.DepartmentID,
		SnapshotID:      req.SnapshotID,
		Title:           req.Title,
		Description:     req.Description,
		RiskLevel:       req.RiskLevel,
		Priority:        req.Priority,
		Category:        req.Category,
		ResponsibleName: req.ResponsibleName,
		ResponsibleID:   req.ResponsibleID,
		Status:          req.Status,
		DueDate:         dueDate,
		EvidenceURLs:    req.EvidenceURLs,
		Notes:           req.Notes,
	}

	if err := h.actionPlanService.Create(c.Request.Context(), actionPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, actionPlan)
}

func (h *ActionPlanHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	actionPlans, err := h.actionPlanService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionPlans)
}

func (h *ActionPlanHandler) GetByID(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	actionPlan, err := h.actionPlanService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action plan not found"})
		return
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (h *ActionPlanHandler) Update(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Title           *string  `json:"title"`
		Description     *string  `json:"description"`
		Priority        *string  `json:"priority"`
		ResponsibleName *string  `json:"responsible_name"`
		ResponsibleID   *int64   `json:"responsible_id"`
		Status          *string  `json:"status"`
		DueDate         *string  `json:"due_date"`
		EvidenceURLs    []string `json:"evidence_urls"`
		Notes           *string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	actionPlan, err := h.actionPlanService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Action plan not found"})
		return
	}

	if req.Title != nil {
		actionPlan.Title = *req.Title
	}
	if req.Description != nil {
		actionPlan.Description = *req.Description
	}
	if req.Priority != nil {
		actionPlan.Priority = *req.Priority
	}
	if req.ResponsibleName != nil {
		actionPlan.ResponsibleName = *req.ResponsibleName
	}
	if req.ResponsibleID != nil {
		actionPlan.ResponsibleID = req.ResponsibleID
	}
	if req.Status != nil {
		actionPlan.Status = *req.Status
	}
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format"})
			return
		}
		actionPlan.DueDate = &parsed
	}
	if req.EvidenceURLs != nil {
		actionPlan.EvidenceURLs = req.EvidenceURLs
	}
	if req.Notes != nil {
		actionPlan.Notes = *req.Notes
	}

	if err := h.actionPlanService.Update(c.Request.Context(), actionPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (h *ActionPlanHandler) Delete(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.actionPlanService.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Action plan deleted successfully"})
}

func (h *ActionPlanHandler) ListByDepartment(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	departmentID, err := strconv.ParseInt(c.Param("departmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	actionPlans, err := h.actionPlanService.ListByDepartment(c.Request.Context(), partnerID, departmentID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionPlans)
}

func (h *ActionPlanHandler) ListByStatus(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status query parameter is required"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	actionPlans, err := h.actionPlanService.ListByStatus(c.Request.Context(), partnerID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actionPlans)
}
