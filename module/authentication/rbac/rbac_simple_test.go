package rbac

import (
	"fmt"
	"os"
	"practicev2/module/authentication/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDBDebug creates a fresh in-memory database for testing
func setupTestDBDebug(t *testing.T, dbName string) *gorm.DB {
	fmt.Printf("🔧 Configurando DB de test: %s\n", dbName)

	// Usar archivo temporal único para cada test
	tempFile := fmt.Sprintf("test_%s_%d.db", dbName, time.Now().UnixNano())

	// Cleanup al final del test
	t.Cleanup(func() {
		if _, err := os.Stat(tempFile); err == nil {
			os.Remove(tempFile)
		}
	})

	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{})
	require.NoError(t, err)

	// Migrate models
	err = db.AutoMigrate(
		&models.Identity{},
		&models.Organization{},
		&models.Role{},
		&models.Permission{},
		&models.OrganizationalMembership{},
	)
	require.NoError(t, err)

	return db
}

func TestRBACSystemDebug(t *testing.T) {
	db := setupTestDBDebug(t, "debug_test_db")

	// 1. Inicializar sistema RBAC
	seeder := NewRBACSeeder(db, DefaultRBACConfig())
	err := seeder.InitializeRBACSystem()
	require.NoError(t, err)

	// Verificar roles de sistema directamente
	var roles []models.Role
	err = db.Where("is_system_role = ?", true).Find(&roles).Error
	require.NoError(t, err)
	require.Equal(t, 7, len(roles), "Deben existir 7 roles de sistema")

	// Buscar específicamente org_admin
	var orgAdminRole models.Role
	err = db.Where("name = ? AND is_system_role = ?", "org_admin", true).First(&orgAdminRole).Error
	require.NoError(t, err, "org_admin debe existir")

	// Intentar usar el método del seeder
	sysRole, err := seeder.GetSystemRoleByName("org_admin")
	require.NoError(t, err, "GetSystemRoleByName debe funcionar")
	require.Equal(t, orgAdminRole.ID, sysRole.ID, "Los IDs deben coincidir")
}

func TestRBACCompleteWorkflow(t *testing.T) {
	db := setupTestDBDebug(t, "complete_flow_test_db")

	// 1. Inicializar sistema RBAC
	seeder := NewRBACSeeder(db, DefaultRBACConfig())
	err := seeder.InitializeRBACSystem()
	require.NoError(t, err)

	// Verificar permisos y roles de sistema
	var permissionCount, systemRoleCount int64
	err = db.Model(&models.Permission{}).Count(&permissionCount).Error
	require.NoError(t, err)
	err = db.Model(&models.Role{}).Where("is_system_role = ?", true).Count(&systemRoleCount).Error
	require.NoError(t, err)

	require.Greater(t, permissionCount, int64(0), "Deben existir permisos")
	require.Equal(t, int64(7), systemRoleCount, "Deben existir 7 roles de sistema")

	// 2. Crear organización con roles
	org := &models.Organization{
		ID:          generateID(),
		Name:        "Test Company",
		Slug:        "test-company",
		Type:        "company",
		Description: "Empresa de prueba",
		IsActive:    true,
	}

	err = seeder.CreateOrganizationWithDefaultRoles(org)
	require.NoError(t, err)

	// Verificar roles de organización
	var orgRoleCount int64
	err = db.Model(&models.Role{}).Where("organization_id = ?", org.ID).Count(&orgRoleCount).Error
	require.NoError(t, err)
	require.Greater(t, orgRoleCount, int64(0), "La organización debe tener roles")

	// 3. Crear usuario y asignar a organización
	user := &models.Identity{
		ID:            generateID(),
		Email:         "test@company.com",
		PasswordHash:  "hash",
		FirstName:     "Test",
		LastName:      "User",
		EmailVerified: true,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Obtener rol employee
	var employeeRole models.Role
	err = db.Where("name = ? AND organization_id = ?", "employee", org.ID).First(&employeeRole).Error
	require.NoError(t, err)

	// Crear membresía organizacional
	membership := &models.OrganizationalMembership{
		ID:             generateID(),
		IdentityID:     user.ID,
		OrganizationID: org.ID,
		RoleID:         employeeRole.ID,
		ActiveFrom:     db.NowFunc(),
		IsActive:       true,
	}
	err = db.Create(membership).Error
	require.NoError(t, err)

	// 4. Verificar el flujo completo
	var membershipCount int64
	err = db.Model(&models.OrganizationalMembership{}).
		Where("organization_id = ? AND is_active = ?", org.ID, true).
		Count(&membershipCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), membershipCount, "Debe haber 1 membresía activa")

	// Verificar que se pueden contar todos los elementos correctamente
	require.Greater(t, permissionCount, int64(300), "Debe haber al menos 300 permisos")
	require.Equal(t, int64(7), systemRoleCount, "Debe haber exactamente 7 roles de sistema")
	require.Equal(t, int64(4), orgRoleCount, "Debe haber exactamente 4 roles organizacionales")
	require.Equal(t, int64(1), membershipCount, "Debe haber exactamente 1 membresía")
}

func checkTablesDebug(t *testing.T, db *gorm.DB) {
	var tables []string
	err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name").Scan(&tables).Error
	if err != nil {
		fmt.Printf("❌ Error listando tablas: %v\n", err)
		return
	}

	fmt.Printf("� Tablas existentes (%d):\n", len(tables))
	for _, table := range tables {
		fmt.Printf("   - %s\n", table)
	}
}
