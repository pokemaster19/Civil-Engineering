package routers

import (
	"ControlSystem/handlers"
	"ControlSystem/midlleware"
	"ControlSystem/models"

	"github.com/gin-gonic/gin"
)

func RegisterDefectsRoutes(r *gin.RouterGroup, s *handlers.Server) {

	h := s.DefectHandler
	r.Use(midlleware.JWTMiddleware())
	r.POST("/", midlleware.RoleMiddleware(models.RoleEngineer), h.CreateDefect)
	r.POST("/:defectId/comments", h.LeaveComment)
	r.GET("/:defectId/comments", h.GetComments)
	r.GET("/", h.GetDefects)
	r.GET("/:defectId", h.GetDefectById)
	r.PATCH("/:defectId", midlleware.RoleMiddleware(models.RoleEngineer), h.UpdateDefect)
}
