package services

import (
	"ControlSystem/models"
	"ControlSystem/repositories"
	"ControlSystem/storage"
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type AttachmentService struct {
	repo        *repositories.AttachmentRepository
	projectRepo *repositories.ProjectRepository
	defectRepo  *repositories.DefectRepository
	userRepo    *repositories.UserRepository
	minioClient *storage.MinioClient
}

func NewAttachmentService(
	repo *repositories.AttachmentRepository,
	projectRepo *repositories.ProjectRepository,
	defectRepo *repositories.DefectRepository,
	userRepo *repositories.UserRepository,
	minioClient *storage.MinioClient,
) *AttachmentService {
	return &AttachmentService{
		repo:        repo,
		projectRepo: projectRepo,
		defectRepo:  defectRepo,
		userRepo:    userRepo,
		minioClient: minioClient,
	}
}

type UploadAttachmentInput struct {
	FileName    string
	FileContent []byte
	ContentType string
	ProjectID   *uint
	DefectID    *uint
	UploadedBy  uint
	RoleID      uint
}

func (s *AttachmentService) UploadAttachment(input UploadAttachmentInput) (*models.Attachment, error) {
	if input.DefectID == nil && input.ProjectID == nil {
		return nil, errors.New("either defectId or projectId must be provided")
	}

	if input.DefectID != nil {
		if input.RoleID != models.RoleEngineer && input.RoleID <= models.RoleAdmin {
			return nil, errors.New("forbidden: only engineers can attach files to defects")
		}
	}
	if input.ProjectID != nil {
		if input.RoleID != models.RoleManager && input.RoleID <= models.RoleAdmin {
			return nil, errors.New("forbidden: only managers can attach files to projects")
		}
	}

	if input.ProjectID != nil {
		_, err := s.projectRepo.GetByID(*input.ProjectID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("project not found")
			}
			return nil, fmt.Errorf("failed to fetch project: %w", err)
		}
	}
	if input.DefectID != nil {
		_, err := s.defectRepo.GetByID(*input.DefectID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("defect not found")
			}
			return nil, fmt.Errorf("failed to fetch defect: %w", err)
		}
	}

	ext := strings.ToLower(filepath.Ext(input.FileName))
	var bucketName string
	switch ext {
	case ".png", ".jpg", ".jpeg":
		bucketName = "images"
	case ".pdf", ".docx":
		bucketName = "files"
	default:
		return nil, errors.New("unsupported file type")
	}

	fileName := fmt.Sprintf("file-%d%s", time.Now().UnixNano(), ext)

	_, err := s.minioClient.Client.PutObject(context.Background(), bucketName, fileName, bytes.NewReader(input.FileContent), int64(len(input.FileContent)), minio.PutObjectOptions{
		ContentType: input.ContentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	attachment := models.Attachment{
		DefectID:   input.DefectID,
		ProjectID:  input.ProjectID,
		FileName:   fileName,
		FilePath:   bucketName,
		FileType:   ext,
		UploadedBy: input.UploadedBy,
	}
	if err := s.repo.Create(&attachment); err != nil {
		return nil, fmt.Errorf("failed to save attachment to database: %w", err)
	}

	return &attachment, nil
}
