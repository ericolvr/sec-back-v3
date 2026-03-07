package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/core/services"
	"github.com/ericolvr/sec-back-v2/internal/core/utils"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	var dto dto.UserRequest

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Generate 6-digit numeric password
	plainPassword, err := utils.GenerateNumericPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate password",
		})
		return
	}

	// Log senha gerada
	fmt.Printf("🔑 Senha gerada para usuário %s (mobile: %s): %s\n", dto.Name, dto.Mobile, plainPassword)

	// Hash password
	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := &domain.User{
		PartnerID:  tenantID,
		Name:      dto.Name,
		Mobile:    dto.Mobile,
		Password:  hashedPassword,
		Type:      dto.Type,
		Active:    dto.Active,
		CreatedAt: time.Now(),
	}

	if err := h.userService.CreateWithPassword(c.Request.Context(), user, plainPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao criar usuário",
			"details": err.Error(),
		})
		return
	}

	response := h.toUserResponse(user)

	c.JSON(http.StatusCreated, gin.H{
		"user":    response,
		"message": "Usuário criado com sucesso. SMS enviado com a senha.",
	})
}

func (h *UserHandler) List(c *gin.Context) {
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

	users, err := h.userService.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao listar usuários",
			"details": err.Error(),
		})
		return
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, h.toUserResponse(user))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), tenantID, idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Usuário não encontrado",
			"details": err.Error(),
		})
		return
	}

	response := h.toUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Update(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	id := c.Param("id")

	var updateDTO dto.UserRequest
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

	user := &domain.User{
		ID:       idInt,
		PartnerID: tenantID,
		Name:     updateDTO.Name,
		Mobile:   updateDTO.Mobile,
		Type:     updateDTO.Type,
		Active:   updateDTO.Active,
	}

	if err := h.userService.Update(c.Request.Context(), user); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuário não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	responseDTO := h.toUserResponse(user)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *UserHandler) Delete(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	deletedUser, err := h.userService.Delete(c.Request.Context(), tenantID, id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuário não encontrado",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuário deletado com sucesso",
		"data":    h.toUserResponse(deletedUser),
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	tenantID := middleware.GetPartnerID(c)
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID deve ser um número válido",
		})
		return
	}

	// Get user
	user, err := h.userService.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuário não encontrado",
		})
		return
	}

	// Generate new 6-digit password
	plainPassword, err := utils.GenerateNumericPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate password",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(plainPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Update user password
	user.Password = hashedPassword
	if err := h.userService.UpdatePasswordWithSMS(c.Request.Context(), user, plainPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar senha",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Senha resetada com sucesso. SMS enviado para o usuário.",
		"mobile":  user.Mobile,
	})
}

func (h *UserHandler) toUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        int(user.ID),
		PartnerID:  int(user.PartnerID),
		Name:      user.Name,
		Mobile:    user.Mobile,
		Type:      user.Type,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}
}
