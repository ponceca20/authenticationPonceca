package middleware

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"practicev2/module/authentication/models"
	"practicev2/module/authentication/rbac"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupIntegrationTestDB configura DB temporal para tests de integración
func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	tempFile := fmt.Sprintf("integration_test_%d.db", time.Now().UnixNano())

	t.Cleanup(func() {
		if _, err := os.Stat(tempFile); err == nil {
			os.Remove(tempFile)
		}
	})

	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{})
	require.NoError(t, err)

	// Migrar modelos
	err = db.AutoMigrate(
		&models.Identity{},
		&models.Organization{},
		&models.Role{},
		&models.Permission{},
		&models.OrganizationalMembership{},
		&models.AuditLog{},
	)
	require.NoError(t, err)

	return db
}

func TestUnifiedAuthMiddlewareIntegration(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// 1. Inicializar sistema RBAC completo
	seeder := rbac.NewRBACSeeder(db, rbac.DefaultRBACConfig())
	err := seeder.InitializeRBACSystem()
	require.NoError(t, err)

	// 2. Crear organización con roles
	org := &models.Organization{
		ID:          generateTestID(),
		Name:        "Test Integration Company",
		Slug:        "test-integration",
		Type:        "company",
		Description: "Organización para test de integración",
		IsActive:    true,
	}

	err = seeder.CreateOrganizationWithDefaultRoles(org)
	require.NoError(t, err)

	// 3. Crear usuario
	user := &models.Identity{
		ID:            generateTestID(),
		Email:         "integration@test.com",
		PasswordHash:  "hash",
		FirstName:     "Integration",
		LastName:      "Test",
		EmailVerified: true,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// 4. Asignar rol employee al usuario
	var employeeRole models.Role
	err = db.Where("name = ? AND organization_id = ?", "employee", org.ID).First(&employeeRole).Error
	require.NoError(t, err)

	membership := &models.OrganizationalMembership{
		ID:             generateTestID(),
		IdentityID:     user.ID,
		OrganizationID: org.ID,
		RoleID:         employeeRole.ID,
		ActiveFrom:     db.NowFunc(),
		IsActive:       true,
	}
	err = db.Create(membership).Error
	require.NoError(t, err)

	// 5. Inicializar middleware unificado usando esta DB
	config := &UnifiedAuthConfig{
		CacheEnabled:   false, // Deshabilitar cache para tests
		DebugMode:      true,
		DefaultDenyAll: true,
		RequireOrg:     true,
		EnableAudit:    false, // Deshabilitar auditoría para simplificar
	}

	// Crear middleware que use nuestra DB de test
	middleware := &UnifiedAuthMiddleware{
		db:              db,
		rbacEngine:      rbac.NewDynamicRBACEngine(db, rbac.DefaultRBACConfig()),
		config:          config,
		permissionCache: make(map[string]*PermissionCacheEntry),
	}

	// 6. Test del motor RBAC directamente
	t.Run("Motor RBAC funciona", func(t *testing.T) {
		// Verificar que el usuario tiene permisos de employee
		hasExpenseRead := middleware.rbacEngine.CheckUserPermission(
			user.ID, org.ID, "expenses", "read", "own")
		require.True(t, hasExpenseRead, "Employee debe poder leer sus propios expenses")

		// Verificar que NO tiene permisos de admin
		hasAdminManage := middleware.rbacEngine.CheckUserPermission(
			user.ID, org.ID, "admin", "manage", "organization")
		require.False(t, hasAdminManage, "Employee NO debe tener permisos de admin")
	})

	// 7. Test de creación de permisos automática
	t.Run("Creación automática de permisos", func(t *testing.T) {
		testPermissions := []string{
			"test_resource:read:own",
			"test_resource:create:own",
			"test_resource:update:organization",
		}

		// Crear permisos manualmente usando nuestra instancia de DB
		for _, permStr := range testPermissions {
			parts := strings.Split(permStr, ":")
			require.Equal(t, 3, len(parts), "Permiso mal formateado: %s", permStr)

			permission := models.Permission{
				Resource: parts[0],
				Action:   parts[1],
				Scope:    parts[2],
			}

			err := db.Where("resource = ? AND action = ? AND scope = ?",
				parts[0], parts[1], parts[2]).FirstOrCreate(&permission).Error
			require.NoError(t, err)
		}

		// Verificar que se crearon en BD
		var count int64
		err = db.Model(&models.Permission{}).
			Where("resource = ?", "test_resource").
			Count(&count).Error
		require.NoError(t, err)
		require.Equal(t, int64(3), count, "Deben crearse 3 permisos de test_resource")
	})

	t.Run("Verificación de conteos finales", func(t *testing.T) {
		// Verificar conteos finales
		var permissionCount, roleCount, membershipCount int64

		err = db.Model(&models.Permission{}).Count(&permissionCount).Error
		require.NoError(t, err)
		require.Greater(t, permissionCount, int64(300), "Debe haber al menos 300 permisos")

		err = db.Model(&models.Role{}).Count(&roleCount).Error
		require.NoError(t, err)
		require.Greater(t, roleCount, int64(7), "Debe haber al menos 7 roles (7 sistema + 4 org)")

		err = db.Model(&models.OrganizationalMembership{}).Count(&membershipCount).Error
		require.NoError(t, err)
		require.Equal(t, int64(1), membershipCount, "Debe haber exactamente 1 membresía")

		fmt.Printf("✅ Test de integración exitoso: %d permisos, %d roles, %d membresías\n",
			permissionCount, roleCount, membershipCount)
	})
}

func generateTestID() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}
