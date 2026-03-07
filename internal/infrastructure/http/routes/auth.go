package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authHandler *api.AuthHandler
}

func NewAuthRoutes(authHandler *api.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

func (ar *AuthRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	// Login não precisa de X-Tenant-ID - busca usuário globalmente por mobile
	{
		auth.POST("/login", ar.authHandler.Login) // Rota pública (sem JWT e sem Tenant)
	}
}
