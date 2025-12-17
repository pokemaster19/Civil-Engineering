package handlers

import (
	"ControlSystem/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DefectHandler struct {
	service *services.DefectService
}

func NewDefectHandler(service *services.DefectService) *DefectHandler {
	return &DefectHandler{service: service}
}

func (h *DefectHandler) CreateDefect(c *gin.Context) {
	type CreateDefectInput struct {
		Title       string `json:"title" binding:"required,min=3"`
		Description string `json:"description" binding:"required"`
		ProjectID   uint   `json:"project_id" binding:"required"`
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "user not identified"})
		return
	}

	var input CreateDefectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defect, err := h.service.CreateDefect(services.CreateDefectInput{
		Title:       input.Title,
		Description: input.Description,
		ProjectID:   input.ProjectID,
	}, userId.(uint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "project not found"})
		} else {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(201, gin.H{
		"message": "defect created successfully",
		"defect":  defect,
	})
}

func (h *DefectHandler) GetDefects(c *gin.Context) {
	type GetDefectsInput struct {
		ProjectID uint `form:"projectId" binding:"required"`
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "8")
	search := c.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}

	var input GetDefectsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetDefects(services.GetDefectsInput{
		ProjectID: input.ProjectID,
		Page:      page,
		Limit:     limit,
		Search:    search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"defects": result.Defects,
		"pagination": gin.H{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"totalPages":  result.TotalPages,
			"hasNextPage": result.Page < result.TotalPages,
			"hasPrevPage": result.Page > 1,
		},
		"statusCounts": result.Counts,
	})
}

func (h *DefectHandler) GetDefectById(c *gin.Context) {
	defectId, exists := c.Params.Get("defectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "defectId is required"})
		return
	}

	defectIDUint, err := strconv.ParseUint(defectId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid defectId: " + err.Error()})
		return
	}

	details, err := h.service.GetDefectByID(uint(defectIDUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Defect not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve defect"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"defect": details})
}

func (h *DefectHandler) UpdateDefect(c *gin.Context) {
	defectId, exists := c.Params.Get("defectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "defectId is required"})
		return
	}

	defectIDUint, err := strconv.ParseUint(defectId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid defectId: " + err.Error()})
		return
	}

	var input services.UpdateDefectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	defect, err := h.service.UpdateDefect(uint(defectIDUint), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Defect not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, defect)
}

func (h *DefectHandler) LeaveComment(c *gin.Context) {
	type CommentInput struct {
		Content string `json:"content" binding:"required,max=255"`
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	roleId, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: role not found"})
		return
	}

	defectId, exists := c.Params.Get("defectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "defectId is required"})
		return
	}

	var input CommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	defectIDUint, err := strconv.ParseUint(defectId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid defectId: " + err.Error()})
		return
	}

	comment, err := h.service.AddComment(services.AddCommentInput{
		Content:  input.Content,
		DefectID: uint(defectIDUint),
	}, userId.(uint), roleId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Comment added successfully",
		"id":         comment.ID,
		"content":    comment.Content,
		"authorName": comment.Creator.FirstName,
	})
}

func (h *DefectHandler) GetComments(c *gin.Context) {
	defectId, exists := c.Params.Get("defectId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "defectId is required"})
		return
	}

	defectIDUint, err := strconv.ParseUint(defectId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid defectId: " + err.Error()})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}

	result, err := h.service.GetComments(services.GetCommentsInput{
		DefectID: uint(defectIDUint),
		Page:     page,
		Limit:    limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": result.Comments,
		"pagination": gin.H{
			"page":        result.Page,
			"limit":       result.Limit,
			"total":       result.Total,
			"totalPages":  result.TotalPages,
			"hasNextPage": result.Page < result.TotalPages,
			"hasPrevPage": result.Page > 1,
		},
	})
}
