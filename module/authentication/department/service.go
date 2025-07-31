package department

import (
	"errors"
	"fmt"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/user"
	"practicev2/module/authentication/utils"

	"github.com/google/uuid"
)

// DepartmentService defines the interface for department-related business logic.
type DepartmentService interface {
	CreateDepartment(orgID string, dto *DepartmentDTO) (*models.Department, error)
	GetDepartment(orgID, deptID string) (*models.Department, error)
	ListDepartments(orgID string) ([]models.Department, error)
	UpdateDepartment(orgID, deptID string, dto *DepartmentDTO) (*models.Department, error)
	DeleteDepartment(orgID, deptID string) error
	ListUsersInDepartment(orgID, deptID string) ([]user.UserDTO, error)
}

type departmentService struct {
	repo DepartmentRepository
}

// NewDepartmentService creates a new instance of DepartmentService.
func NewDepartmentService(repo DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

// CreateDepartment handles the logic for creating a new department.
func (s *departmentService) CreateDepartment(orgID string, dto *DepartmentDTO) (*models.Department, error) {
	if errs := utils.ValidateStruct(dto); errs != nil {
		return nil, fmt.Errorf("validation failed: %v", errs)
	}

	dept := &models.Department{
		ID:             uuid.New().String(),
		OrganizationID: orgID,
		Name:           dto.Name,
		Description:    dto.Description,
		ParentID:       dto.ParentID,
	}

	if err := s.repo.CreateDepartment(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// GetDepartment retrieves a single department.
func (s *departmentService) GetDepartment(orgID, deptID string) (*models.Department, error) {
	return s.repo.FindDepartmentByID(orgID, deptID)
}

// ListDepartments retrieves all departments for an organization.
func (s *departmentService) ListDepartments(orgID string) ([]models.Department, error) {
	return s.repo.ListDepartmentsByOrganization(orgID)
}

// UpdateDepartment updates a department's details.
func (s *departmentService) UpdateDepartment(orgID, deptID string, dto *DepartmentDTO) (*models.Department, error) {
	dept, err := s.repo.FindDepartmentByID(orgID, deptID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	if dto.Name != "" {
		dept.Name = dto.Name
	}
	if dto.Description != "" {
		dept.Description = dto.Description
	}
	dept.ParentID = dto.ParentID // Allow unsetting the parent

	if err := s.repo.UpdateDepartment(dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// DeleteDepartment deletes a department.
func (s *departmentService) DeleteDepartment(orgID, deptID string) error {
	dept, err := s.repo.FindDepartmentByID(orgID, deptID)
	if err != nil {
		return errors.New("department not found")
	}
	return s.repo.DeleteDepartment(dept)
}

// ListUsersInDepartment retrieves all users assigned to a specific department.
func (s *departmentService) ListUsersInDepartment(orgID, deptID string) ([]user.UserDTO, error) {
	// First, ensure the department exists in the organization
	_, err := s.repo.FindDepartmentByID(orgID, deptID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	memberships, err := s.repo.FindMembershipsByDepartment(orgID, deptID)
	if err != nil {
		return nil, err
	}

	var userDTOs []user.UserDTO
	for _, m := range memberships {
		if m.Identity.ID != "" {
			userDTOs = append(userDTOs, user.ToUserDTO(&m))
		}
	}

	return userDTOs, nil
}
