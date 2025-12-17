package handlers

import (
	"ControlSystem/services"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service *services.AttachmentService
}

func NewAttachmentHandler(service *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

func (h *AttachmentHandler) UploadAttachment(c *gin.Context) {
	roleID, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("FormFile error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	fileContent, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	projectIDStr := c.PostForm("projectId")
	defectIDStr := c.PostForm("defectId")

	var projectID *uint
	if projectIDStr != "" {
		pID, _ := strconv.ParseUint(projectIDStr, 10, 32)
		tmp := uint(pID)
		projectID = &tmp
	}

	var defectID *uint
	if defectIDStr != "" {
		dID, _ := strconv.ParseUint(defectIDStr, 10, 32)
		tmp := uint(dID)
		defectID = &tmp
	}

	input := services.UploadAttachmentInput{
		FileName:    file.Filename,
		FileContent: fileContent,
		ContentType: file.Header.Get("Content-Type"),
		ProjectID:   projectID,
		DefectID:    defectID,
		UploadedBy:  userID.(uint),
		RoleID:      roleID.(uint),
	}

	attachment, err := h.service.UploadAttachment(input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "project not found" || err.Error() == "defect not found" {
			statusCode = http.StatusNotFound
		} else if strings.HasPrefix(err.Error(), "forbidden:") || err.Error() == "unsupported file type" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "either defectId or projectId must be provided" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "file uploaded successfully",
		"attachment": attachment.ID,
	})
}
