package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/gin-gonic/gin"
)

type AssessmentAssignmentHandler struct {
	assignmentService *services.AssessmentAssignmentService
}

func NewAssessmentAssignmentHandler(assignmentService *services.AssessmentAssignmentService) *AssessmentAssignmentHandler {
	return &AssessmentAssignmentHandler{
		assignmentService: assignmentService,
	}
}

func (h *AssessmentAssignmentHandler) Create(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	var req struct {
		TemplateID    int64   `json:"template_id" binding:"required"`
		DepartmentIDs []int64 `json:"department_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var successAssignments []*domain.AssessmentAssignment
	var errors []map[string]interface{}

	assignment := &domain.AssessmentAssignment{
		PartnerID:     partnerID,
		TemplateID:    req.TemplateID,
		DepartmentIDs: req.DepartmentIDs,
		Active:        true,
	}

	submissionsCount, err := h.assignmentService.Create(c.Request.Context(), assignment)
	if err != nil {
		// Log detalhado para debug
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"details": map[string]interface{}{
				"partner_id":     partnerID,
				"template_id":    req.TemplateID,
				"department_ids": req.DepartmentIDs,
			},
		})
		return
	}

	successAssignments = append(successAssignments, assignment)

	// Se todos falharam, retorna erro
	if len(successAssignments) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create any assignments",
			"errors": errors,
		})
		return
	}

	// Se houve erros parciais, retorna sucesso com erros
	if len(errors) > 0 {
		c.JSON(http.StatusCreated, gin.H{
			"success":             successAssignments,
			"submissions_created": submissionsCount,
			"errors":              errors,
		})
		return
	}

	// Sucesso total
	c.JSON(http.StatusCreated, gin.H{
		"success":             successAssignments,
		"submissions_created": submissionsCount,
	})
}

func (h *AssessmentAssignmentHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	assignments, err := h.assignmentService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

func (h *AssessmentAssignmentHandler) GetByID(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	assignment, err := h.assignmentService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func (h *AssessmentAssignmentHandler) Close(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	assignment, err := h.assignmentService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	assignment.Active = false
	if err := h.assignmentService.Deactivate(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment closed successfully"})
}

func (h *AssessmentAssignmentHandler) Update(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		TemplateID    *int64  `json:"template_id,omitempty"`
		DepartmentIDs []int64 `json:"department_ids,omitempty"`
		Active        *bool   `json:"active,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.assignmentService.Update(c.Request.Context(), partnerID, id, req.TemplateID, req.DepartmentIDs, req.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	assignment, err := h.assignmentService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated assignment"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func (h *AssessmentAssignmentHandler) Delete(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.assignmentService.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}
