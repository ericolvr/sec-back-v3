package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "partner_id"

func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantIDStr := c.GetHeader("X-Partner-ID")
		
		if tenantIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "X-Partner-ID header is required",
			})
			c.Abort()
			return
		}

		tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
		if err != nil || tenantID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid X-Partner-ID header",
			})
			c.Abort()
			return
		}

		c.Set(TenantIDKey, tenantID)
		c.Next()
	}
}

func GetPartnerID(c *gin.Context) int64 {
	tenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return 0
	}
	return tenantID.(int64)
}
