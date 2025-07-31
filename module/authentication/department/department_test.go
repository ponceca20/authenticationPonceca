package department

import (
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DepartmentTestSuite struct {
	suite.Suite
	db          *gorm.DB
	deptService DepartmentService
	testOrg     *models.Organization
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *DepartmentTestSuite) setupTestDatabase() *gorm.DB {
	// Initialize config with default values for tests
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the models we need for testing
	err = db.AutoMigrate(
		&models.Identity{},
		&models.User{},
		&models.UserProfile{},
		&models.Organization{},
		&models.OrganizationalMembership{},
		&models.Department{},
		&models.Role{},
		&models.RefreshToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *DepartmentTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	deptRepo := NewDepartmentRepository(suite.db)
	suite.deptService = NewDepartmentService(deptRepo)

	// We need an organization to exist for these tests
	orgRepo := organization.NewOrganizationRepository(suite.db)
	authRepo := auth.NewAuthRepository(suite.db)
	orgService := organization.NewOrganizationService(orgRepo, authRepo)

	dto := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Dept",
			LastName:  "Head",
			Email:     "dept.head@example.com",
			Password:  "password123",
		},
		Name: "Department Test Corp",
		Type: "company",
	}
	org, err := orgService.CreateOrganization(dto)
	assert.NoError(suite.T(), err)
	suite.testOrg = org
}

func TestDepartmentTestSuite(t *testing.T) {
	suite.Run(t, new(DepartmentTestSuite))
}

func (suite *DepartmentTestSuite) TestDepartmentCRUD() {
	// 1. Create a department
	dto := &DepartmentDTO{Name: "Engineering"}
	createdDept, err := suite.deptService.CreateDepartment(suite.testOrg.ID, dto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Engineering", createdDept.Name)

	// 2. Get the department
	fetchedDept, err := suite.deptService.GetDepartment(suite.testOrg.ID, createdDept.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdDept.ID, fetchedDept.ID)

	// 3. List departments
	depts, err := suite.deptService.ListDepartments(suite.testOrg.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), depts, 1)

	// 4. Update the department
	updateDTO := &DepartmentDTO{Name: "Software Engineering"}
	updatedDept, err := suite.deptService.UpdateDepartment(suite.testOrg.ID, createdDept.ID, updateDTO)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Software Engineering", updatedDept.Name)

	// 5. Delete the department
	err = suite.deptService.DeleteDepartment(suite.testOrg.ID, createdDept.ID)
	assert.NoError(suite.T(), err)

	// 6. Verify deletion
	_, err = suite.deptService.GetDepartment(suite.testOrg.ID, createdDept.ID)
	assert.Error(suite.T(), err)
}
