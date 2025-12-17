package routers

import (
	"ControlSystem/handlers"
	"ControlSystem/midlleware"
	"ControlSystem/models"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup, s *handlers.Server) {

	h := s.AdminHandler

	r.Use(midlleware.JWTMiddleware())
	r.POST("/register", midlleware.RoleMiddleware(), h.RegisterNewUser)
	r.PATCH("/edit-user/:userId", midlleware.RoleMiddleware(), h.EditUserInfo)
	r.GET("/get-users", midlleware.RoleMiddleware(models.RoleManager), h.GetUsers)
}
