package audit

import (
	"practicev2/config"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/organization"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuditTestSuite struct {
	suite.Suite
	db           *gorm.DB
	auditService AuditService
	auditRepo    AuditRepository
	testOrg      *models.Organization
	testUser     *models.Identity
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func (suite *AuditTestSuite) setupTestDatabase() *gorm.DB {
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
		&models.AuditLog{},
		&models.Role{},
		&models.RefreshToken{},
	)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func (suite *AuditTestSuite) SetupSuite() {
	suite.db = suite.setupTestDatabase()
	suite.auditRepo = NewAuditRepository(suite.db)
	suite.auditService = NewAuditService(suite.auditRepo)

	// Create test organization and user
	authRepo := auth.NewAuthRepository(suite.db)
	orgRepo := organization.NewOrganizationRepository(suite.db)
	orgService := organization.NewOrganizationService(orgRepo, authRepo)

	orgDTO := &organization.OrganizationRegistrationDTO{
		Identity: auth.RegisterDTO{
			FirstName: "Audit",
			LastName:  "User",
			Email:     "audit.user@example.com",
			Password:  "audit-password",
		},
		Name: "Audit Test Corp",
		Type: "company",
	}
	org, err := orgService.CreateOrganization(orgDTO)
	assert.NoError(suite.T(), err)
	suite.testOrg = org

	user, err := authRepo.FindIdentityByEmail(orgDTO.Identity.Email)
	assert.NoError(suite.T(), err)
	suite.testUser = user

	// Create some test audit logs
	suite.createTestAuditLogs()
}

func (suite *AuditTestSuite) createTestAuditLogs() {
	logs := []models.AuditLog{
		{
			ID:             "audit-1",
			OrganizationID: &suite.testOrg.ID,
			IdentityID:     suite.testUser.ID,
			Action:         "login",
			Resource:       "authentication",
			ResourceID:     "auth-resource-1",
			Status:         "success",
			IPAddress:      "192.168.1.100",
			Details:        "User logged in successfully",
			Timestamp:      time.Now().Add(-2 * time.Hour),
		},
		{
			ID:             "audit-2",
			OrganizationID: &suite.testOrg.ID,
			IdentityID:     suite.testUser.ID,
			Action:         "create",
			Resource:       "organization",
			ResourceID:     suite.testOrg.ID,
			Status:         "success",
			IPAddress:      "192.168.1.100",
			Details:        "Organization created",
			Timestamp:      time.Now().Add(-1 * time.Hour),
		},
		{
			ID:             "audit-3",
			OrganizationID: &suite.testOrg.ID,
			IdentityID:     suite.testUser.ID,
			Action:         "update",
			Resource:       "profile",
			ResourceID:     "profile-1",
			Status:         "failed",
			IPAddress:      "192.168.1.101",
			Details:        "Profile update failed - validation error",
			Timestamp:      time.Now().Add(-30 * time.Minute),
		},
	}

	for _, log := range logs {
		log.Identity = *suite.testUser // Set the identity relation
		suite.db.Create(&log)
	}
}

func TestAuditTestSuite(t *testing.T) {
	suite.Run(t, new(AuditTestSuite))
}

func (suite *AuditTestSuite) TestListAuditLogs() {
	// 1. Test basic listing with pagination
	query := &AuditQueryDTO{
		Page:     1,
		PageSize: 10,
	}
	response, err := suite.auditService.ListAuditLogs(suite.testOrg.ID, query)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), int64(3), response.Total)
	assert.Len(suite.T(), response.Data, 3)
	assert.Equal(suite.T(), 1, response.Page)
	assert.Equal(suite.T(), 10, response.PageSize)
	assert.Equal(suite.T(), 1, response.TotalPages)

	// 2. Test filtering by action
	queryByAction := &AuditQueryDTO{
		Action:   "login",
		Page:     1,
		PageSize: 10,
	}
	response, err = suite.auditService.ListAuditLogs(suite.testOrg.ID, queryByAction)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), response.Total)
	assert.Len(suite.T(), response.Data, 1)
	assert.Equal(suite.T(), "login", response.Data[0].Action)

	// 3. Test filtering by resource
	queryByResource := &AuditQueryDTO{
		Resource: "organization",
		Page:     1,
		PageSize: 10,
	}
	response, err = suite.auditService.ListAuditLogs(suite.testOrg.ID, queryByResource)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), response.Total)
	assert.Len(suite.T(), response.Data, 1)
	assert.Equal(suite.T(), "organization", response.Data[0].Resource)

	// 4. Test pagination
	queryPage2 := &AuditQueryDTO{
		Page:     2,
		PageSize: 2,
	}
	response, err = suite.auditService.ListAuditLogs(suite.testOrg.ID, queryPage2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), response.Total)
	assert.Equal(suite.T(), 2, response.Page)
	assert.Equal(suite.T(), 2, response.PageSize)
	assert.Equal(suite.T(), 2, response.TotalPages) // ceil(3/2) = 2
}
