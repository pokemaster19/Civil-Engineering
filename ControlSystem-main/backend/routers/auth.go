package routers

import (
	"ControlSystem/handlers"
	"ControlSystem/midlleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, s *handlers.Server) {
	r.POST("/login", s.Login)
	r.POST("/refresh", s.RefreshTokenHandler)
	r.POST("logout", midlleware.JWTMiddleware(), s.Logout)
}
