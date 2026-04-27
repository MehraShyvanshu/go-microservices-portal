package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleInterface, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		userRole := roleInterface.(string)

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		c.Abort()
	}
}
