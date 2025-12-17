package midlleware

import (
	"ControlSystem/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		roleUint := role.(uint)
		if roleUint == models.RoleAdmin || roleUint == models.RoleSuperAdmin {
			c.Next()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleUint == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
