package api

import (
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type InvitationHandler struct {
	invitationService *services.InvitationService
	assignmentService *services.AssessmentAssignmentService
}

func NewInvitationHandler(
	invitationService *services.InvitationService,
	assignmentService *services.AssessmentAssignmentService,
) *InvitationHandler {
	return &InvitationHandler{
		invitationService: invitationService,
		assignmentService: assignmentService,
	}
}

func (h *InvitationHandler) List(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	limit := int64(50)
	offset := int64(0)

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = parsedOffset
		}
	}

	invitations, err := h.invitationService.List(c.Request.Context(), partnerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar convites"})
		return
	}

	var response []dto.InvitationResponse
	for _, inv := range invitations {
		response = append(response, h.toInvitationResponse(inv))
	}

	c.JSON(http.StatusOK, response)
}

func (h *InvitationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	partnerID := c.GetInt64("partner_id")

	invitation, err := h.invitationService.GetByID(c.Request.Context(), partnerID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Convite não encontrado"})
		return
	}

	c.JSON(http.StatusOK, h.toInvitationResponse(invitation))
}

func (h *InvitationHandler) GetSummary(c *gin.Context) {
	partnerID := c.GetInt64("partner_id")

	assignments, err := h.assignmentService.List(c.Request.Context(), partnerID, 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar assignments"})
		return
	}

	var summaries []dto.InvitationSummaryResponse

	for _, assignment := range assignments {
		// Para cada department no assignment
		for i, deptID := range assignment.DepartmentIDs {
			invitations, err := h.invitationService.ListByTemplateAndDepartment(
				c.Request.Context(),
				partnerID,
				assignment.TemplateID,
				deptID,
			)
			if err != nil || len(invitations) == 0 {
				continue
			}

			pendingCount := 0
			sentCount := 0

			for _, inv := range invitations {
				if inv.Status == domain.InvitationStatusPending {
					pendingCount++
				} else if inv.Status == domain.InvitationStatusSent {
					sentCount++
				}
			}

			departmentName := ""
			if i < len(assignment.Data) {
				departmentName = assignment.Data[i].DepartmentName
			}

			summary := dto.InvitationSummaryResponse{
				TemplateID:       assignment.TemplateID,
				TemplateName:     assignment.TemplateName,
				DepartmentID:     deptID,
				DepartmentName:   departmentName,
				TotalInvitations: len(invitations),
				PendingCount:     pendingCount,
				SentCount:        sentCount,
			}

			summaries = append(summaries, summary)
		}
	}

	c.JSON(http.StatusOK, summaries)
}

func (h *InvitationHandler) MarkAsSent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	partnerID := c.GetInt64("partner_id")

	if err := h.invitationService.MarkAsSent(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convite marcado como enviado"})
}

func (h *InvitationHandler) MarkAsFailed(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	partnerID := c.GetInt64("partner_id")

	if err := h.invitationService.MarkAsFailed(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convite marcado como falha"})
}

func (h *InvitationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	partnerID := c.GetInt64("partner_id")

	if err := h.invitationService.Delete(c.Request.Context(), partnerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convite removido com sucesso"})
}

func (h *InvitationHandler) toInvitationResponse(invitation *domain.Invitation) dto.InvitationResponse {
	var sentAt string
	if invitation.SentAt != nil {
		sentAt = invitation.SentAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return dto.InvitationResponse{
		ID:              invitation.ID,
		PartnerID:       invitation.PartnerID,
		TemplateID:      invitation.TemplateID,
		TemplateName:    invitation.TemplateName,
		DepartmentID:    invitation.DepartmentID,
		ResponseID:      invitation.ResponseID,
		InvitationToken: invitation.EmployeeEmail,
		EmployeeEmail:   invitation.EmployeeEmail,
		Status:          invitation.Status,
		SentAt:          sentAt,
		CreatedAt:       invitation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       invitation.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
