package handlers

import (
	"ControlSystem/services"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	type CreateProjectRequest struct {
		Name        string `json:"name" binding:"required,min=3"`
		Description string `json:"description"`
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := services.CreateProjectInput{
		Name:        req.Name,
		Description: req.Description,
	}

	project, err := h.service.CreateProject(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

func (h *ProjectHandler) EditProjectInfo(c *gin.Context) {
	type EditProjectRequest struct {
		Name        string `json:"name" binding:"omitempty"`
		Description string `json:"description" binding:"omitempty"`
		Status      uint   `json:"status" binding:"omitempty"`
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req EditProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := services.UpdateProjectInput{}
	if req.Name != "" {
		input.Name = &req.Name
	}
	if req.Description != "" {
		input.Description = &req.Description
	}
	if req.Status > 0 {
		input.Status = &req.Status
	}

	project, err := h.service.UpdateProject(uint(projectID), input)
	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project updated", "project": project})
}

func (h *ProjectHandler) AssignEngineer(c *gin.Context) {
	type AssignEngineerRequest struct {
		EngineerId uint `json:"engineer_id" binding:"required"`
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req AssignEngineerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.AssignEngineer(uint(projectID), req.EngineerId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" ||
			err.Error() == fmt.Sprintf("engineer with ID %d not found", req.EngineerId) {
			statusCode = http.StatusNotFound
		} else if err.Error() == fmt.Sprintf("user with ID %d is not an engineer", req.EngineerId) ||
			err.Error() == fmt.Sprintf("engineer with ID %d is already assigned to project %d", req.EngineerId, projectID) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "engineer assigned successfully",
		"project_id":  projectID,
		"engineer_id": req.EngineerId,
	})
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	type ProjectResponse struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		PhotoURL string `json:"photoUrl,omitempty"`
	}

	roleID, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}

	search := c.DefaultQuery("search", "")

	result, err := h.service.GetProjects(userID.(uint), roleID.(uint), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ProjectResponse, len(result.Projects))
	for i, p := range result.Projects {
		response[i] = ProjectResponse{
			ID:       p.ID,
			Name:     p.Name,
			PhotoURL: p.PhotoURL,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": response,
		"pagination": gin.H{
			"page":       result.Page,
			"limit":      result.Limit,
			"total":      result.Total,
			"totalPages": result.TotalPages,
		},
	})
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid projectId"})
		return
	}

	details, err := h.service.GetProjectDetails(uint(projectID))
	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project_name":        details.Name,
		"project_description": details.Description,
	})
}

func (h *ProjectHandler) ExportDefectsCSV(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	defects, summary, err := h.service.GetDefectsForExport(uint(projectID))
	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error exporting defects: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export defects"})
		}
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=project_%d_defects_report.csv", projectID))

	if _, err := c.Writer.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		log.Printf("Error writing UTF-8 BOM: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write CSV"})
		return
	}

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{
		"ID", "Title", "Description", "Priority",
		"Open", "In Progress", "Resolved", "Overdue",
		"Assigned To", "Created By", "Deadline",
	}
	if err := writer.Write(headers); err != nil {
		log.Printf("Error writing CSV header: %v", err)
		return
	}

	priorityMap := map[uint]string{1: "Low", 2: "Medium", 3: "High"}

	for _, defect := range defects {
		deadlineStr := ""
		if defect.Deadline != nil {
			deadlineStr = defect.Deadline.Format("2006-01-02")
		}

		row := []string{
			strconv.Itoa(int(defect.ID)),
			defect.Title,
			defect.Description,
			priorityMap[defect.Priority],
			"", "", "", "",
			strconv.Itoa(int(defect.AssignedTo)),
			strconv.Itoa(int(defect.CreatedBy)),
			deadlineStr,
		}
		if err := writer.Write(row); err != nil {
			log.Printf("Error writing CSV row: %v", err)
			return
		}
	}

	summaryRow := []string{
		"", "Summary", "", "",
		strconv.Itoa(summary.Open),
		strconv.Itoa(summary.InProgress),
		strconv.Itoa(summary.Resolved),
		strconv.Itoa(summary.Overdue),
		"", "", "",
	}
	if err := writer.Write(summaryRow); err != nil {
		log.Printf("Error writing CSV summary: %v", err)
	}
}
