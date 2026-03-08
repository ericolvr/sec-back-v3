package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	answerService     *services.AnswerService
	submissionService *services.EmployeeSubmissionService
}

func NewAnswerHandler(answerService *services.AnswerService, submissionService *services.EmployeeSubmissionService) *AnswerHandler {
	return &AnswerHandler{
		answerService:     answerService,
		submissionService: submissionService,
	}
}

func (h *AnswerHandler) Create(c *gin.Context) {
	var req dto.AnswerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	tenantID := c.GetInt64("tenant_id")

	answer := &domain.Answer{
		PartnerID:    tenantID,
		SubmissionID: req.ResponseID,
		QuestionID:   req.QuestionID,
		Value:        req.Value,
	}

	if err := h.answerService.Create(c.Request.Context(), answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar answer",
			"details": err.Error(),
		})
		return
	}

	response := h.toAnswerResponse(answer)

	c.JSON(http.StatusCreated, response)
}

func (h *AnswerHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	responseIDStr := c.Query("response_id")

	if responseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "response_id is required",
		})
		return
	}

	responseID, err := strconv.ParseInt(responseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "response_id must be a valid number",
		})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	tenantID := c.GetInt64("tenant_id")

	answers, err := h.answerService.List(c.Request.Context(), tenantID, responseID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar answers",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.AnswerResponse
	for _, answer := range answers {
		responses = append(responses, h.toAnswerResponse(answer))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *AnswerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	tenantID := c.GetInt64("tenant_id")

	answer, err := h.answerService.GetByID(c.Request.Context(), tenantID, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Answer não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toAnswerResponse(answer)
	c.JSON(http.StatusOK, response)
}

func (h *AnswerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	tenantID := c.GetInt64("tenant_id")

	answer := &domain.Answer{
		ID:           idInt,
		PartnerID:    tenantID,
		SubmissionID: req.ResponseID,
		QuestionID:   req.QuestionID,
		Value:        req.Value,
	}

	if err := h.answerService.Update(c.Request.Context(), answer); err != nil {
		if err.Error() == "answer not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Answer não encontrado",
			})
			return
		}
		if err.Error() == "cannot update answer: response already completed" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Não é possível atualizar resposta: template já foi finalizado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toAnswerResponse(answer)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *AnswerHandler) SubmitBatch(c *gin.Context) {
	var req struct {
		Token     string `json:"token" binding:"required"`
		Responses []struct {
			QuestionID int64       `json:"question_id" binding:"required"`
			Answer     interface{} `json:"answer" binding:"required"`
		} `json:"responses" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Buscar Response pelo token (sem tenant_id pois é público)
	response, err := h.submissionService.GetByToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token inválido ou expirado",
		})
		return
	}

	// Verificar se já foi completado
	if response.Status == domain.SubmissionStatusCompleted {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Este template já foi respondido",
		})
		return
	}

	// Salvar todas as respostas
	savedCount := 0
	var errors []map[string]interface{}

	for _, answerReq := range req.Responses {
		// Converter answer para string
		var answerValue string
		switch v := answerReq.Answer.(type) {
		case string:
			answerValue = v
		case float64:
			answerValue = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			answerValue = strconv.Itoa(v)
		default:
			answerValue = fmt.Sprintf("%v", v)
		}

		answer := &domain.Answer{
			PartnerID:    response.PartnerID,
			SubmissionID: response.ID,
			QuestionID:   answerReq.QuestionID,
			Value:        answerValue,
		}

		if err := h.answerService.Create(c.Request.Context(), answer); err != nil {
			errors = append(errors, map[string]interface{}{
				"question_id": answerReq.QuestionID,
				"error":       err.Error(),
			})
		} else {
			savedCount++
		}
	}

	// Atualizar status do Response para completed
	response.Status = domain.SubmissionStatusCompleted
	// TODO: Implementar submissionService.Update quando necessário
	// if err := h.submissionService.Update(c.Request.Context(), response); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error":   "Erro ao finalizar template",
	// 		"details": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":     "Questionário enviado com sucesso",
		"total":       len(req.Responses),
		"saved":       savedCount,
		"failed":      len(errors),
		"errors":      errors,
		"response_id": response.ID,
		"status":      response.Status,
	})
}

func (h *AnswerHandler) toAnswerResponse(answer *domain.Answer) dto.AnswerResponse {
	return dto.AnswerResponse{
		ID:         answer.ID,
		PartnerID:  answer.PartnerID,
		ResponseID: answer.SubmissionID,
		QuestionID: answer.QuestionID,
		Value:      answer.Value,
		Score:      answer.Score,
		CreatedAt:  answer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  answer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
