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

type EmployeeHandler struct {
	employeeService *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: service,
	}
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	var dto dto.EmployeeRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	employee := &domain.Employee{
		PartnerID:     tenantID,
		Name:         dto.Name,
		DepartmentID: dto.DepartmentID,
		Email:        dto.Email,
		Active:       true,
		CreatedAt:    time.Now(),
	}

	if err := h.employeeService.Create(c.Request.Context(), employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar employee",
			"details": err.Error(),
		})
		return
	}

	response := h.toEmployeeResponse(employee)

	c.JSON(http.StatusCreated, response)
}

func (h *EmployeeHandler) List(c *gin.Context) {
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

	employees, err := h.employeeService.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar employees",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.EmployeeResponse
	for _, employee := range employees {
		responses = append(responses, h.toEmployeeResponse(employee))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	employee, err := h.employeeService.GetByID(c.Request.Context(), tenantID, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Employee não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toEmployeeResponse(employee)
	c.JSON(http.StatusOK, response)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	var updateDTO dto.EmployeeRequest
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

	employee := &domain.Employee{
		ID:           idInt,
		PartnerID:     tenantID,
		DepartmentID: updateDTO.DepartmentID,
		Name:         updateDTO.Name,
		Email:        updateDTO.Email,
		Active:       updateDTO.Active,
	}

	if err := h.employeeService.Update(c.Request.Context(), employee); err != nil {
		if err.Error() == "employee not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Employee não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toEmployeeResponse(employee)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	deletedEmployee, err := h.employeeService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "employee not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Employee não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deletada com sucesso",
		"data":    h.toEmployeeResponse(deletedEmployee),
	})
}

func (h *EmployeeHandler) toEmployeeResponse(employee *domain.Employee) dto.EmployeeResponse {
	response := dto.EmployeeResponse{
		ID:           int(employee.ID),
		DepartmentID: employee.DepartmentID,
		Name:         employee.Name,
		Email:        employee.Email,
		Active:       employee.Active,
		CreatedAt:    employee.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if employee.Department != nil {
		response.DepartmentName = employee.Department.Name
	}

	return response
}
