package routers

import (
	"ControlSystem/handlers"
	"ControlSystem/midlleware"
	"ControlSystem/models"

	"github.com/gin-gonic/gin"
)

func RegisterAttachmentsRoutes(r *gin.RouterGroup, s *handlers.Server) {

	h := s.AttachmentHandler

	r.Use(midlleware.JWTMiddleware())
	r.POST("/", midlleware.RoleMiddleware(models.RoleEngineer, models.RoleManager), h.UploadAttachment)
}
