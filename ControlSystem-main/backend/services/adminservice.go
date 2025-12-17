package services

import (
	"ControlSystem/models"
	"ControlSystem/repositories"
	"ControlSystem/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/resend/resend-go/v2"
	"gorm.io/gorm"
)

type AdminService struct {
	db           *gorm.DB
	resendClient *resend.Client
	userRepo     *repositories.UserRepository
	adminRepo    *repositories.AdminRepository
}

func NewAdminService(db *gorm.DB, resendClient *resend.Client, userRepo *repositories.UserRepository, adminRepo *repositories.AdminRepository) *AdminService {
	return &AdminService{
		db:           db,
		resendClient: resendClient,
		userRepo:     userRepo,
		adminRepo:    adminRepo,
	}
}

type RegisterUserInput struct {
	FirstName  string
	MiddleName string
	LastName   string
	OrigEmail  string
	Role       uint
}

type EditUserInput struct {
	FirstName  *string
	MiddleName *string
	LastName   *string
	Role       *uint
	IsEnabled  *bool
}

type UserListInput struct {
	Page         int
	Limit        int
	EmailFilter  string
	RoleFilter   string
	StatusFilter string
}

type UserListResult struct {
	Users      []UserResponse
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type UserResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Role       uint   `json:"role"`
	IsEnabled  bool   `json:"is_enabled"`
}

func (s *AdminService) RegisterUser(input RegisterUserInput, currentUserRole uint) (string, error) {

	corporateEmail, err := utils.GenerateCorporateEmail(input.FirstName, input.LastName, s.db)
	if err != nil {
		return "", errors.New("failed to generate unique corporate email")
	}

	password, err := utils.GeneratePassword(12)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	if input.Role >= currentUserRole {
		return "", errors.New("insufficient permissions to assign this role")
	}

	newUser := models.User{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Email:      corporateEmail,
		Password:   password,
		Role:       input.Role,
	}

	newUser.HashPassword()

	if err := s.adminRepo.CreateUser(&newUser); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	err = utils.SendRegistrationEmail(input.OrigEmail, corporateEmail, password, s.resendClient)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return corporateEmail, nil
}

func (s *AdminService) EditUser(userID, currentUserRole uint, input EditUserInput) (*models.User, error) {

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Role >= currentUserRole || user.ID != userID {
		return nil, errors.New("insufficient permissions to edit this user")
	}

	if input.Role != nil && *input.Role >= currentUserRole {
		return nil, errors.New("insufficient permissions to assign this role")
	}

	updates := make(map[string]interface{})
	if input.FirstName != nil && *input.FirstName != "" {
		updates["first_name"] = *input.FirstName
	}
	if input.MiddleName != nil && *input.MiddleName != "" {
		updates["middle_name"] = *input.MiddleName
	}
	if input.LastName != nil && *input.LastName != "" {
		updates["last_name"] = *input.LastName
	}
	if input.Role != nil {
		updates["role"] = *input.Role
	}
	if input.IsEnabled != nil {
		updates["is_enabled"] = *input.IsEnabled
	}

	if len(updates) == 0 {
		return user, nil
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *AdminService) GetUsers(currentUserID, currentUserRole uint, input UserListInput) (*UserListResult, error) {
	if input.Page < 1 {
		return nil, errors.New("invalid page number")
	}
	if input.Limit < 1 {
		return nil, errors.New("invalid limit value")
	}

	offset := (input.Page - 1) * input.Limit

	query := s.db.Model(&models.User{}).Where("id <> ?", currentUserID)

	if input.EmailFilter != "" {
		query = query.Where("LOWER(email) LIKE ?", "%"+strings.ToLower(input.EmailFilter)+"%")
	}

	if input.RoleFilter != "" {
		if roleValue, err := strconv.Atoi(input.RoleFilter); err == nil {
			query = query.Where("role = ?", roleValue)
		}
	}

	if input.StatusFilter != "" {
		if status, err := strconv.ParseBool(input.StatusFilter); err == nil {
			query = query.Where("is_enabled = ?", status)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	var users []models.User
	if err := query.Offset(offset).Limit(input.Limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	response := make([]UserResponse, 0, len(users))
	for _, u := range users {
		response = append(response, UserResponse{
			ID:         u.ID,
			FirstName:  u.FirstName,
			MiddleName: u.MiddleName,
			LastName:   u.LastName,
			Email:      u.Email,
			Role:       u.Role,
			IsEnabled:  u.IsEnabled,
		})
	}

	totalPages := int((total + int64(input.Limit) - 1) / int64(input.Limit))

	return &UserListResult{
		Users:      response,
		Total:      total,
		Page:       input.Page,
		Limit:      input.Limit,
		TotalPages: totalPages,
	}, nil
}
