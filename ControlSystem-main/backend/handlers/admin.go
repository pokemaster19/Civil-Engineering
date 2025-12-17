package handlers

import (
	"ControlSystem/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) RegisterNewUser(c *gin.Context) {
	type RegisterInput struct {
		FirstName  string `json:"first_name" binding:"required,min=1"`
		MiddleName string `json:"middle_name" binding:"required,min=1"`
		LastName   string `json:"last_name" binding:"required,min=1"`
		OrigEmail  string `json:"email" binding:"required,email"`
		Role       uint   `json:"role" binding:"required"`
	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	roleID, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	serviceInput := services.RegisterUserInput{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		OrigEmail:  input.OrigEmail,
		Role:       input.Role,
	}

	corporateEmail, err := h.service.RegisterUser(serviceInput, roleID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "admin access required" || err.Error() == "insufficient permissions to assign this role" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "invalid input" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"corporate_email": corporateEmail})
}

func (h *AdminHandler) EditUserInfo(c *gin.Context) {
	type EditUserInput struct {
		FirstName  string `json:"first_name" binding:"omitempty,min=1"`
		MiddleName string `json:"middle_name" binding:"omitempty,min=1"`
		LastName   string `json:"last_name" binding:"omitempty,min=1"`
		Role       uint   `json:"role" binding:"omitempty"`
		IsEnabled  bool   `json:"is_enabled" binding:"omitempty"`
	}

	roleID, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var input EditUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	serviceInput := services.EditUserInput{}
	if input.FirstName != "" {
		serviceInput.FirstName = &input.FirstName
	}
	if input.MiddleName != "" {
		serviceInput.MiddleName = &input.MiddleName
	}
	if input.LastName != "" {
		serviceInput.LastName = &input.LastName
	}
	if input.Role != 0 {
		serviceInput.Role = &input.Role
	}
	if input.IsEnabled {
		serviceInput.IsEnabled = &input.IsEnabled
	}

	user, err := h.service.EditUser(uint(userID), roleID.(uint), serviceInput)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "admin access required" || err.Error() == "insufficient permissions to edit this user" || err.Error() == "insufficient permissions to assign this role" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
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

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	emailFilter := c.Query("email")
	roleFilter := c.Query("role")
	statusFilter := c.Query("is_enabled")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit value"})
		return
	}

	input := services.UserListInput{
		Page:         page,
		Limit:        limit,
		EmailFilter:  emailFilter,
		RoleFilter:   roleFilter,
		StatusFilter: statusFilter,
	}

	result, err := h.service.GetUsers(userID.(uint), roleID.(uint), input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "admin access required" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "invalid page number" || err.Error() == "invalid limit value" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": result.Users,
		"pagination": gin.H{
			"page":       result.Page,
			"limit":      result.Limit,
			"total":      result.Total,
			"totalPages": result.TotalPages,
		},
	})
}
