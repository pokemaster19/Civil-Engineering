package routers

import (
	"ControlSystem/handlers"
	"ControlSystem/midlleware"
	"ControlSystem/models"

	"github.com/gin-gonic/gin"
)

func RegisterProjectsRoutes(r *gin.RouterGroup, s *handlers.Server) {

	h := s.ProjectHandler

	r.Use(midlleware.JWTMiddleware())
	r.POST("/", midlleware.RoleMiddleware(models.RoleManager), h.CreateProject)
	r.PATCH("/:projectId", midlleware.RoleMiddleware(models.RoleManager), h.EditProjectInfo)
	r.GET("/", h.GetProjects)
	r.POST("/:projectId/assign", midlleware.RoleMiddleware(models.RoleManager), h.AssignEngineer)
	r.GET("/:projectId", h.GetProject)
	r.GET("/:projectId/export", midlleware.RoleMiddleware(models.RoleManager, models.RoleObserver), h.ExportDefectsCSV)
}
