package midlleware

import (
	"ControlSystem/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := utils.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong token"})
			c.Abort()
			return
		}

		token, err := utils.GetToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token parsing error"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Dead token"})
			c.Abort()
			return
		}

		userId, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user id"})
			c.Abort()
			return
		}

		roleId, ok := claims["role"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user role"})
			c.Abort()
			return
		}

		c.Set("user_id", uint(userId))
		c.Set("role", uint(roleId))
		c.Next()
	}

}
