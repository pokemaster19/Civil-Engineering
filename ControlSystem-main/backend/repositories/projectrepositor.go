package repositories

import (
	"ControlSystem/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	return &project, err
}

func (r *ProjectRepository) Update(project *models.Project, updates map[string]interface{}) error {
	return r.db.Model(project).Updates(updates).Error
}

func (r *ProjectRepository) AssignEngineer(project *models.Project, user *models.User) error {
	return r.db.Model(project).Association("Users").Append(user)
}

func (r *ProjectRepository) IsEngineerAssigned(projectID, engineerID uint) (bool, error) {
	var count int64
	err := r.db.Table("user_project").
		Where("project_id = ? AND user_id = ?", projectID, engineerID).
		Count(&count).Error
	return count > 0, err
}

func (r *ProjectRepository) GetProjectsByUser(userID uint, offset, limit int, search string) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{}).
		Joins("JOIN user_project ON user_project.project_id = projects.id").
		Where("user_project.user_id = ?", userID)

	if search != "" {
		query = query.Where("projects.name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Select("projects.id, projects.name").
		Offset(offset).Limit(limit).
		Scan(&projects).Error

	return projects, total, err
}

func (r *ProjectRepository) GetAllProjects(offset, limit int, search string) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.db.Model(&models.Project{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Select("id, name").
		Offset(offset).Limit(limit).
		Scan(&projects).Error

	return projects, total, err
}
