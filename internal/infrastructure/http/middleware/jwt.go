package middleware

import (
	"net/http"
	"strings"

	"github.com/ericolvr/sec-back-v2/internal/core/utils"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware valida o token JWT e injeta os dados do usuário no contexto
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrai o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			c.Abort()
			return
		}

		// Remove o prefixo "Bearer " do token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Valida o token e extrai os claims
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Injeta os dados do usuário no contexto para uso nos handlers
		c.Set("user_id", claims.UserID)
		c.Set("partner_id", claims.TenantID)
		c.Set("mobile", claims.Mobile)
		c.Set("user_type", claims.Type)

		c.Next()
	}
}

// GetUserID extrai o user_id do contexto
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUserType extrai o user_type do contexto
func GetUserType(c *gin.Context) int {
	userType, exists := c.Get("user_type")
	if !exists {
		return 0
	}
	return userType.(int)
}

// GetMobile extrai o mobile do contexto
func GetMobile(c *gin.Context) string {
	mobile, exists := c.Get("mobile")
	if !exists {
		return ""
	}
	return mobile.(string)
}
