package services

import (
	"ControlSystem/models"
	"ControlSystem/repositories"
	"ControlSystem/storage"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ProjectService struct {
	repo        *repositories.ProjectRepository
	userRepo    *repositories.UserRepository
	defectRepo  *repositories.DefectRepository
	attachRepo  *repositories.AttachmentRepository
	minioClient *storage.MinioClient
}

func NewProjectService(
	repo *repositories.ProjectRepository,
	userRepo *repositories.UserRepository,
	defectRepo *repositories.DefectRepository,
	attachRepo *repositories.AttachmentRepository,
	minioClient *storage.MinioClient,
) *ProjectService {
	return &ProjectService{
		repo:        repo,
		userRepo:    userRepo,
		defectRepo:  defectRepo,
		attachRepo:  attachRepo,
		minioClient: minioClient,
	}
}

type CreateProjectInput struct {
	Name        string
	Description string
}

type UpdateProjectInput struct {
	Name        *string
	Description *string
	Status      *uint
}

type ProjectListResult struct {
	Projects   []ProjectWithPhoto
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type ProjectWithPhoto struct {
	ID       uint
	Name     string
	PhotoURL string
}

type ProjectDetails struct {
	Name        string
	Description string
}

type DefectExportData struct {
	ID          uint
	Title       string
	Description string
	Priority    uint
	Status      uint
	AssignedTo  uint
	CreatedBy   uint
	Deadline    *time.Time
}

type DefectStatusSummary struct {
	Open       int
	InProgress int
	Resolved   int
	Overdue    int
}

func (s *ProjectService) CreateProject(input CreateProjectInput) (models.Project, error) {
	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
		Status:      1,
	}

	if err := s.repo.Create(&project); err != nil {
		return models.Project{}, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(projectID uint, input UpdateProjectInput) (*models.Project, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}

	updates := make(map[string]interface{})
	if input.Name != nil && *input.Name != "" {
		updates["name"] = *input.Name
	}
	if input.Description != nil && *input.Description != "" {
		updates["description"] = *input.Description
	}
	if input.Status != nil && *input.Status > 0 {
		updates["status"] = *input.Status
	}

	if len(updates) == 0 {
		return project, nil
	}

	if err := s.repo.Update(project, updates); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) AssignEngineer(projectID, engineerID uint) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("project not found")
		}
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	engineer, err := s.userRepo.GetByID(engineerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("engineer with ID %d not found", engineerID)
		}
		return fmt.Errorf("failed to fetch engineer: %w", err)
	}

	if engineer.Role != models.RoleEngineer {
		return fmt.Errorf("user with ID %d is not an engineer", engineerID)
	}

	isAssigned, err := s.repo.IsEngineerAssigned(projectID, engineerID)
	if err != nil {
		return fmt.Errorf("failed to check engineer assignment: %w", err)
	}
	if isAssigned {
		return fmt.Errorf("engineer with ID %d is already assigned to project %d", engineerID, projectID)
	}

	if err := s.repo.AssignEngineer(project, engineer); err != nil {
		return fmt.Errorf("failed to assign engineer: %w", err)
	}

	return nil
}

func (s *ProjectService) GetProjects(userID, roleID uint, page, limit int, search string) (*ProjectListResult, error) {
	offset := (page - 1) * limit

	var projects []models.Project
	var total int64
	var err error

	if roleID == models.RoleEngineer {
		projects, total, err = s.repo.GetProjectsByUser(userID, offset, limit, search)
	} else {
		projects, total, err = s.repo.GetAllProjects(offset, limit, search)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	projectIDs := make([]uint, len(projects))
	for i, p := range projects {
		projectIDs[i] = p.ID
	}

	var projectsWithPhotos []ProjectWithPhoto
	if len(projectIDs) > 0 {
		attachments, err := s.attachRepo.GetByProjectIDs(projectIDs)
		if err == nil {
			attMap := make(map[uint]models.Attachment)
			for _, a := range attachments {
				if a.ProjectID != nil {
					attMap[*a.ProjectID] = a
				}
			}

			for _, p := range projects {
				pwp := ProjectWithPhoto{
					ID:   p.ID,
					Name: p.Name,
				}

				if att, ok := attMap[p.ID]; ok {
					if url, err := s.minioClient.GetFileURL(att.FilePath, att.FileName, 24*time.Hour); err == nil {
						pwp.PhotoURL = url
					}
				}

				projectsWithPhotos = append(projectsWithPhotos, pwp)
			}
		}
	} else {
		for _, p := range projects {
			projectsWithPhotos = append(projectsWithPhotos, ProjectWithPhoto{
				ID:   p.ID,
				Name: p.Name,
			})
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &ProjectListResult{
		Projects:   projectsWithPhotos,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *ProjectService) GetProjectDetails(projectID uint) (*ProjectDetails, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}

	return &ProjectDetails{
		Name:        project.Name,
		Description: project.Description,
	}, nil
}

func (s *ProjectService) GetDefectsForExport(projectID uint) ([]DefectExportData, *DefectStatusSummary, error) {
	_, err := s.repo.GetByID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("project not found")
		}
		return nil, nil, fmt.Errorf("failed to fetch project: %w", err)
	}

	defects, err := s.defectRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch defects: %w", err)
	}

	var exportData []DefectExportData
	summary := &DefectStatusSummary{}

	for _, d := range defects {
		exportData = append(exportData, DefectExportData{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Priority:    d.Priority,
			Status:      d.Status,
			AssignedTo:  d.AssignedTo,
			CreatedBy:   d.CreatedBy,
			Deadline:    d.Deadline,
		})

		switch d.Status {
		case 1:
			summary.Open++
		case 2:
			summary.InProgress++
		case 3:
			summary.Resolved++
		case 4:
			summary.Overdue++
		}
	}

	return exportData, summary, nil
}
