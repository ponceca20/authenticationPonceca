package department

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

// DepartmentRepository defines the interface for database operations related to departments.
type DepartmentRepository interface {
	CreateDepartment(dept *models.Department) error
	FindDepartmentByID(orgID, deptID string) (*models.Department, error)
	ListDepartmentsByOrganization(orgID string) ([]models.Department, error)
	UpdateDepartment(dept *models.Department) error
	DeleteDepartment(dept *models.Department) error
	FindMembershipsByDepartment(orgID, deptID string) ([]models.OrganizationalMembership, error)
}

type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new instance of DepartmentRepository.
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

// CreateDepartment saves a new department to the database.
func (r *departmentRepository) CreateDepartment(dept *models.Department) error {
	return r.db.Create(dept).Error
}

// FindDepartmentByID finds a department by its ID within a specific organization.
func (r *departmentRepository) FindDepartmentByID(orgID, deptID string) (*models.Department, error) {
	var dept models.Department
	if err := r.db.Where("organization_id = ? AND id = ?", orgID, deptID).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

// ListDepartmentsByOrganization retrieves all departments for a given organization.
func (r *departmentRepository) ListDepartmentsByOrganization(orgID string) ([]models.Department, error) {
	var depts []models.Department
	if err := r.db.Where("organization_id = ?", orgID).Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

// UpdateDepartment saves changes to a department model.
func (r *departmentRepository) UpdateDepartment(dept *models.Department) error {
	return r.db.Save(dept).Error
}

// DeleteDepartment soft deletes a department.
func (r *departmentRepository) DeleteDepartment(dept *models.Department) error {
	return r.db.Delete(dept).Error
}

// FindMembershipsByDepartment finds all memberships associated with a specific department.
func (r *departmentRepository) FindMembershipsByDepartment(orgID, deptID string) ([]models.OrganizationalMembership, error) {
	var memberships []models.OrganizationalMembership
	err := r.db.
		Preload("Identity").
		Preload("Role").
		Where("organization_id = ? AND department = ?", orgID, deptID).
		Find(&memberships).Error
	return memberships, err
}
