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

type DepartmentHandler struct {
	departmentService *services.DepartmentService
}

func NewDepartmentHandler(service *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: service,
	}
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	var dto dto.DepartmentRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	department := &domain.Department{
		PartnerID: tenantID,
		CompanyID: dto.CompanyID,
		Name:      dto.Name,
		Active:    dto.Active,
		CreatedAt: time.Now(),
	}

	if err := h.departmentService.Create(c.Request.Context(), department); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar departamento",
			"details": err.Error(),
		})
		return
	}

	response := h.toDepartmentResponse(department)

	c.JSON(http.StatusCreated, response)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	companyIDStr := c.Query("company_id")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Se company_id for fornecido, filtra por company
	if companyIDStr != "" {
		companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "company_id deve ser um número válido",
			})
			return
		}

		departments, err := h.departmentService.ListByCompany(c.Request.Context(), tenantID, companyID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao listar departamentos",
				"details": err.Error(),
			})
			return
		}

		var responses []dto.DepartmentResponse
		for _, department := range departments {
			responses = append(responses, h.toDepartmentResponse(department))
		}

		c.JSON(http.StatusOK, responses)
		return
	}

	// Lista todos os departments do partner
	departments, err := h.departmentService.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar departamentos",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.DepartmentResponse
	for _, department := range departments {
		responses = append(responses, h.toDepartmentResponse(department))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	department, err := h.departmentService.GetByID(c.Request.Context(), tenantID, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Departamento não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toDepartmentResponse(department)
	c.JSON(http.StatusOK, response)
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	var updateDTO dto.DepartmentUpdateRequest
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	department := &domain.Department{
		ID:        int64(idInt),
		PartnerID: tenantID,
		Name:      updateDTO.Name,
		Active:    updateDTO.Active,
	}

	if err := h.departmentService.Update(c.Request.Context(), department); err != nil {
		if err.Error() == "department not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Departamento não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toDepartmentResponse(department)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	deletedDepartment, err := h.departmentService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "department not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Departamento não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Departamento deletado com sucesso",
		"data":    h.toDepartmentResponse(deletedDepartment),
	})
}

func (h *DepartmentHandler) ToggleActive(c *gin.Context) {
	partnerID := middleware.GetPartnerID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	var req dto.ToggleActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.departmentService.ToggleActive(c.Request.Context(), partnerID, id, req.Active); err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Departamento não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message := "Departamento desativado com sucesso"
	if req.Active {
		message = "Departamento reativado com sucesso"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *DepartmentHandler) ListDeleted(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	departments, err := h.departmentService.ListDeleted(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar departamentos",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.DepartmentResponse
	for _, department := range departments {
		responses = append(responses, h.toDepartmentResponse(department))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *DepartmentHandler) toDepartmentResponse(department *domain.Department) dto.DepartmentResponse {
	return dto.DepartmentResponse{
		ID:             int(department.ID),
		PartnerID:      int(department.PartnerID),
		CompanyID:      int(department.CompanyID),
		CompanyName:    department.CompanyName,
		Name:           department.Name,
		Active:         department.Active,
		TotalEmployees: department.TotalEmployees,
		CreatedAt:      department.CreatedAt,
	}
}
