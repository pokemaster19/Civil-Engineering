package services

import (
	"ControlSystem/models"
	"ControlSystem/repositories"
	"ControlSystem/storage"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DefectService struct {
	db          *gorm.DB
	projectRepo *repositories.ProjectRepository
	userRepo    *repositories.UserRepository
	defectRepo  *repositories.DefectRepository
	attachRepo  *repositories.AttachmentRepository
	minio       *storage.MinioClient
}

func NewDefectService(
	db *gorm.DB,
	projectRepo *repositories.ProjectRepository,
	userRepo *repositories.UserRepository,
	defectRepo *repositories.DefectRepository,
	attachRepo *repositories.AttachmentRepository,
	minio *storage.MinioClient,
) *DefectService {
	return &DefectService{
		db:          db,
		projectRepo: projectRepo,
		userRepo:    userRepo,
		defectRepo:  defectRepo,
		attachRepo:  attachRepo,
		minio:       minio,
	}
}

type CreateDefectInput struct {
	Title       string
	Description string
	ProjectID   uint
}

type GetDefectsInput struct {
	ProjectID uint
	Page      int
	Limit     int
	Search    string
}

type DefectResponse struct {
	ID        uint                 `json:"id"`
	Title     string               `json:"title"`
	Priority  uint                 `json:"priority"`
	Status    uint                 `json:"status"`
	CreatedBy uint                 `json:"createdBy"`
	Creator   *CreatorUserResponse `json:"creator,omitempty"`
	PhotoUrl  string               `json:"photoUrl,omitempty"`
}

type DefectDetails struct {
	ID          uint                 `json:"id"`
	ProjectID   uint                 `json:"projectId"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Priority    uint                 `json:"priority"`
	Status      uint                 `json:"status"`
	AssignedTo  uint                 `json:"assignedTo"`
	Assignee    *CreatorUserResponse `json:"assignee,omitempty"`
	CreatedBy   uint                 `json:"createdBy"`
	Creator     *CreatorUserResponse `json:"creator,omitempty"`
	Deadline    *time.Time           `json:"deadline"`
	PhotosUrl   []string             `json:"photosUrl,omitempty"`
	FilesUrl    []string             `json:"filesUrl,omitempty"`
}

type UpdateDefectInput struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    uint   `json:"priority,omitempty"`
	Status      uint   `json:"status,omitempty"`
}

type CreatorUserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type GetDefectsResult struct {
	Defects    []DefectResponse
	Counts     map[string]int64
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type AddCommentInput struct {
	Content  string
	DefectID uint
}

type CommentResponse struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	AuthorName string    `json:"authorName"`
	CreatedAt  time.Time `json:"createdAt"`
}

type GetCommentsInput struct {
	DefectID uint
	Page     int
	Limit    int
}

type GetCommentsResult struct {
	Comments   []CommentResponse
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

func (s *DefectService) CreateDefect(input CreateDefectInput, userID uint) (*models.Defect, error) {
	var project models.Project
	if err := s.db.First(&project, input.ProjectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var user models.User
	if err := s.db.Preload("Projects", "id = ?", input.ProjectID).First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if len(user.Projects) == 0 && user.Role < 4 {
		return nil, errors.New("user not assigned to this project")
	}

	defect := models.Defect{
		Title:       input.Title,
		Description: input.Description,
		ProjectID:   input.ProjectID,
		CreatedBy:   userID,
		AssignedTo:  userID,
		Status:      1,
	}

	if err := s.db.Create(&defect).Error; err != nil {
		return nil, fmt.Errorf("failed to create defect: %w", err)
	}

	return &defect, nil
}

func (s *DefectService) GetDefects(input GetDefectsInput) (*GetDefectsResult, error) {
	if input.Page < 1 {
		return nil, errors.New("invalid page number")
	}
	if input.Limit < 1 {
		return nil, errors.New("invalid limit value")
	}

	type StatusCount struct {
		Status uint
		Count  int64
	}
	var statusCounts []StatusCount
	if err := s.db.Model(&models.Defect{}).
		Where("project_id = ?", input.ProjectID).
		Group("status").
		Select("status, COUNT(*) as count").
		Find(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	counts := map[string]int64{
		"open":        0,
		"in_progress": 0,
		"resolved":    0,
		"overdue":     0,
	}
	for _, sc := range statusCounts {
		switch sc.Status {
		case 1:
			counts["open"] = sc.Count
		case 2:
			counts["in_progress"] = sc.Count
		case 3:
			counts["resolved"] = sc.Count
		case 4:
			counts["overdue"] = sc.Count
		}
	}

	offset := (input.Page - 1) * input.Limit
	var defects []models.Defect
	var total int64

	likeSearch := "%" + input.Search + "%"
	query := s.db.Preload("Creator").Model(&models.Defect{}).
		Where("project_id = ?", input.ProjectID)

	if input.Search != "" {
		query = query.Where("title LIKE ?", likeSearch)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(input.Limit).
		Find(&defects).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defectIDs := make([]uint, len(defects))
	for i, d := range defects {
		defectIDs[i] = d.ID
	}

	attachments, err := s.attachRepo.GetByDefectIDs(defectIDs)
	attMap := make(map[uint]models.Attachment)
	if err == nil {
		for _, a := range attachments {
			if a.DefectID != nil {
				if _, exists := attMap[*a.DefectID]; !exists {
					attMap[*a.DefectID] = a
				}
			}
		}
	}

	defectResponses := make([]DefectResponse, len(defects))
	for i, defect := range defects {
		var creator *CreatorUserResponse
		if defect.Creator.ID != 0 {
			creator = &CreatorUserResponse{
				ID:        defect.Creator.ID,
				FirstName: defect.Creator.FirstName,
				LastName:  defect.Creator.LastName,
				Email:     defect.Creator.Email,
			}
		}

		defectResponses[i] = DefectResponse{
			ID:        defect.ID,
			Title:     defect.Title,
			Priority:  defect.Priority,
			Status:    defect.Status,
			CreatedBy: defect.CreatedBy,
			Creator:   creator,
		}

		if att, ok := attMap[defect.ID]; ok {
			if url, err := s.minio.GetFileURL(att.FilePath, att.FileName, 24*time.Hour); err == nil {
				defectResponses[i].PhotoUrl = url
			}
		}
	}

	totalPages := int((total + int64(input.Limit) - 1) / int64(input.Limit))

	return &GetDefectsResult{
		Defects:    defectResponses,
		Counts:     counts,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *DefectService) GetDefectByID(id uint) (*DefectDetails, error) {
	var defect models.Defect
	if err := s.db.Preload("Creator").Preload("Assignee").First(&defect, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("defect not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var creator *CreatorUserResponse
	if defect.Creator.ID != 0 {
		creator = &CreatorUserResponse{
			ID:        defect.Creator.ID,
			FirstName: defect.Creator.FirstName,
			LastName:  defect.Creator.LastName,
			Email:     defect.Creator.Email,
		}
	}

	var assignee *CreatorUserResponse
	if defect.Assignee.ID != 0 {
		assignee = &CreatorUserResponse{
			ID:        defect.Assignee.ID,
			FirstName: defect.Assignee.FirstName,
			LastName:  defect.Assignee.LastName,
			Email:     defect.Assignee.Email,
		}
	}

	attachments, err := s.attachRepo.GetByDefectIDs([]uint{id})
	var photosUrl, filesUrl []string
	if err == nil {
		for _, a := range attachments {
			url, err := s.minio.GetFileURL(a.FilePath, a.FileName, 24*time.Hour)
			if err != nil {
				continue
			}
			if a.FilePath == "images" {
				photosUrl = append(photosUrl, url)
			} else {
				filesUrl = append(filesUrl, url)
			}
		}
	}

	return &DefectDetails{
		ID:          defect.ID,
		ProjectID:   defect.ProjectID,
		Title:       defect.Title,
		Description: defect.Description,
		Priority:    defect.Priority,
		Status:      defect.Status,
		AssignedTo:  defect.AssignedTo,
		Assignee:    assignee,
		CreatedBy:   defect.CreatedBy,
		Creator:     creator,
		Deadline:    defect.Deadline,
		PhotosUrl:   photosUrl,
		FilesUrl:    filesUrl,
	}, nil
}

func (s *DefectService) AddComment(input AddCommentInput, userID, roleID uint) (*models.Comment, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if roleID < 2 {
		var count int64
		if err := s.db.Model(&models.Defect{}).
			Where("id = ? AND assigned_to = ?", input.DefectID, userID).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		if count == 0 {
			return nil, errors.New("you do not have access to comment on this defect")
		}
	}

	var comment models.Comment
	comment.Content = input.Content
	comment.DefectID = input.DefectID
	comment.CreatedBy = userID

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to save comment: %w", err)
	}

	comment.Creator = user
	return &comment, nil
}

func (s *DefectService) UpdateDefect(id uint, input UpdateDefectInput) (*models.Defect, error) {
	var defect models.Defect
	if err := s.db.First(&defect, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("defect not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Priority != 0 {
		updates["priority"] = input.Priority
	}
	if input.Status != 0 {
		updates["status"] = input.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&defect).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update defect: %w", err)
		}
	}

	return &defect, nil
}

func (s *DefectService) GetComments(input GetCommentsInput) (*GetCommentsResult, error) {
	if input.Page < 1 {
		return nil, errors.New("invalid page number")
	}
	if input.Limit < 1 {
		return nil, errors.New("invalid limit value")
	}

	var total int64
	if err := s.db.Model(&models.Comment{}).
		Where("defect_id = ?", input.DefectID).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count comments: %w", err)
	}

	var comments []models.Comment
	if err := s.db.Preload("Creator").
		Where("defect_id = ?", input.DefectID).
		Order("created_at DESC").
		Offset((input.Page - 1) * input.Limit).
		Limit(input.Limit).
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve comments: %w", err)
	}

	commentResponses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		authorName := ""
		if comment.Creator.ID != 0 {
			authorName = strings.TrimSpace(comment.Creator.FirstName + " " + comment.Creator.LastName)
		}
		commentResponses[i] = CommentResponse{
			ID:         comment.ID,
			Content:    comment.Content,
			AuthorName: authorName,
			CreatedAt:  comment.CreatedAt,
		}
	}

	totalPages := int((total + int64(input.Limit) - 1) / int64(input.Limit))

	return &GetCommentsResult{
		Comments:   commentResponses,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
		TotalPages: totalPages,
	}, nil
}
