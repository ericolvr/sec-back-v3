package routes

import (
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/http/middleware"
	"github.com/ericolvr/sec-back-v2/internal/interfaces/api"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userHandler *api.UserHandler
}

func NewUserRoutes(userHandler *api.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

func (ur *UserRoutes) SetupRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	users.Use(middleware.TenantMiddleware())
	// users.Use(middleware.JWTMiddleware())
	{
		users.POST("", ur.userHandler.Create)
		users.GET("", ur.userHandler.List)
		users.GET("/:id", ur.userHandler.GetByID)
		users.PATCH("/:id", ur.userHandler.Update)
		users.DELETE("/:id", ur.userHandler.Delete)
		users.POST("/:id/reset-password", ur.userHandler.ResetPassword)
	}
}
