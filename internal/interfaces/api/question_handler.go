package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler(service *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: service,
	}
}

func (h *QuestionHandler) Create(c *gin.Context) {
	var req dto.QuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	tenantID := middleware.GetPartnerID(c)

	question := &domain.Question{
		PartnerID:       tenantID,
		QuestionnaireID: req.QuestionnaireID,
		Question:        req.Question,
		Type:            req.Type,
		Category:        req.Category,
		Options:         req.Options,
		ScoreValues:     req.ScoreValues,
		Weight:          req.Weight,
		Required:        req.Required,
		OrderNum:        req.OrderNum,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.questionService.Create(c.Request.Context(), question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar question",
			"details": err.Error(),
		})
		return
	}

	response := h.toQuestionResponse(question)

	c.JSON(http.StatusCreated, response)
}

func (h *QuestionHandler) ListAll(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 100
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	partnerID := middleware.GetPartnerID(c)

	// ListAll sem filtro de questionário - passar 0 como questionnaireID
	questions, err := h.questionService.List(c.Request.Context(), partnerID, 0, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar questions",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.QuestionResponse
	for _, question := range questions {
		responses = append(responses, h.toQuestionResponse(question))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *QuestionHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	questionnaireIDStr := c.Param("id")

	questionnaireID, err := strconv.ParseInt(questionnaireIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "questionnaire_id must be a valid number",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	tenantID := middleware.GetPartnerID(c)

	questions, err := h.questionService.List(c.Request.Context(), tenantID, questionnaireID, int64(limit), int64(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar questions",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.QuestionResponse
	for _, question := range questions {
		responses = append(responses, h.toQuestionResponse(question))
	}

	// Buscar dados do questionário
	questionnaireInfo := dto.QuestionnaireInfo{
		ID:          questionnaireID,
		Name:        "",
		Description: "",
		Active:      true,
	}

	// Se temos perguntas, pegar dados do questionário da primeira pergunta
	if len(questions) > 0 && questions[0].AssessmentTemplate != nil {
		questionnaireInfo.Name = questions[0].AssessmentTemplate.Name
		questionnaireInfo.Description = questions[0].AssessmentTemplate.Description
		questionnaireInfo.Active = questions[0].AssessmentTemplate.Active
	}

	response := dto.QuestionListResponse{
		Questionnaire:  questionnaireInfo,
		TotalQuestions: len(responses),
		Questions:      responses,
	}

	c.JSON(http.StatusOK, response)
}

func (h *QuestionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	tenantID := middleware.GetPartnerID(c)

	question, err := h.questionService.GetByID(c.Request.Context(), tenantID, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Question não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toQuestionResponse(question)
	c.JSON(http.StatusOK, response)
}

func (h *QuestionHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var updateDTO dto.QuestionRequest
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
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

	tenantID := middleware.GetPartnerID(c)

	question := &domain.Question{
		ID:              idInt,
		PartnerID:       tenantID,
		QuestionnaireID: updateDTO.QuestionnaireID,
		Question:        updateDTO.Question,
		Type:            updateDTO.Type,
		Category:        updateDTO.Category,
		Options:         updateDTO.Options,
		ScoreValues:     updateDTO.ScoreValues,
		Weight:          updateDTO.Weight,
		Required:        updateDTO.Required,
		OrderNum:        updateDTO.OrderNum,
	}

	if err := h.questionService.Update(c.Request.Context(), question); err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Question não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toQuestionResponse(question)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *QuestionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	tenantID := middleware.GetPartnerID(c)

	deletedQuestion, err := h.questionService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "question not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Question não encontrada",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question deletada com sucesso",
		"data":    h.toQuestionResponse(deletedQuestion),
	})
}

func (h *QuestionHandler) toQuestionResponse(question *domain.Question) dto.QuestionResponse {
	var questionnaireName string
	if question.AssessmentTemplate != nil {
		questionnaireName = question.AssessmentTemplate.Name
	}

	return dto.QuestionResponse{
		ID:                question.ID,
		PartnerID:         question.PartnerID,
		QuestionnaireID:   question.QuestionnaireID,
		QuestionnaireName: questionnaireName,
		Question:          question.Question,
		Type:              question.Type,
		Category:          question.Category,
		Options:           question.Options,
		ScoreValues:       question.ScoreValues,
		Weight:            question.Weight,
		Required:          question.Required,
		OrderNum:          question.OrderNum,
		CreatedAt:         question.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         question.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
