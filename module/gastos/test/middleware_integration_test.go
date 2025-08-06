package gastos_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/rbac"
	"practicev2/module/authentication/utils"
	"practicev2/module/gastos"
	gastos_model "practicev2/module/gastos/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// cleanupOldTestFiles elimina archivos de test antiguos para evitar acumulación
func cleanupOldTestFiles() {
	files, err := os.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		fileName := file.Name()
		// Limpiar archivos de test de gastos y otros tests antiguos
		if (strings.HasPrefix(fileName, "gastos_middleware_test_") ||
			strings.HasPrefix(fileName, "test_gastos_")) &&
			strings.HasSuffix(fileName, ".db") {
			os.Remove(fileName)
		}
	}
} // setupTestDB configura una base de datos temporal para testing
func setupTestDB(t *testing.T) *gorm.DB {
	// ⚠️ INICIALIZAR CONFIGURACIÓN PARA TESTING
	os.Setenv("JWT_KEY", "test-jwt-key-32-characters-minimum-required")
	os.Setenv("JWT_ACCESS_TTL_MINUTES", "60") // 1 hora para testing
	config.Init()

	// Limpiar archivos de test antiguos al inicio (por si quedaron de ejecuciones anteriores)
	cleanupOldTestFiles()

	// ✅ Usar archivo temporal con nombre único y limpieza garantizada
	tempFile := fmt.Sprintf("test_gastos_%d.db", time.Now().UnixNano())

	// Garantizar limpieza al final del test
	t.Cleanup(func() {
		os.Remove(tempFile)
	})

	db, err := gorm.Open(sqlite.Open(tempFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Migrar todos los modelos necesarios
	err = db.AutoMigrate(
		&models.Identity{},
		&models.Organization{},
		&models.Role{},
		&models.Permission{},
		&models.OrganizationalMembership{},
		&models.AuditLog{},
		// Agregar modelos de gastos
		&gastos_model.Expense{},
	)
	require.NoError(t, err)

	return db
}

// setupTestUser crea un usuario de prueba con permisos
func setupTestUser(t *testing.T, db *gorm.DB) (*models.Identity, *models.Organization, *models.OrganizationalMembership) {
	// 1. Inicializar sistema RBAC
	seeder := rbac.NewRBACSeeder(db, rbac.DefaultRBACConfig())
	err := seeder.InitializeRBACSystem()
	require.NoError(t, err)

	// 2. Crear organización
	org := &models.Organization{
		ID:          generateTestID(),
		Name:        "Test Gastos Company",
		Slug:        "test-gastos",
		Type:        "company",
		Description: "Organización para test de gastos",
		IsActive:    true,
	}

	err = seeder.CreateOrganizationWithDefaultRoles(org)
	require.NoError(t, err)

	// 3. Crear usuario
	user := &models.Identity{
		ID:            generateTestID(),
		Email:         "gastos@test.com",
		PasswordHash:  "hash",
		FirstName:     "Gastos",
		LastName:      "Test",
		EmailVerified: true,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// 4. Asignar rol employee
	var employeeRole models.Role
	err = db.Where("name = ? AND organization_id = ?", "employee", org.ID).First(&employeeRole).Error
	require.NoError(t, err)

	// 5. Asignar permisos específicos de gastos al rol employee
	expensePermissions := []models.Permission{
		{Resource: "expenses", Action: "read", Scope: "own"},
		{Resource: "expenses", Action: "read", Scope: "organization"},
		{Resource: "expenses", Action: "create", Scope: "own"},
		{Resource: "expenses", Action: "update", Scope: "own"},
		{Resource: "expenses", Action: "delete", Scope: "own"},
	}

	for _, perm := range expensePermissions {
		// Buscar o crear el permiso
		var permission models.Permission
		err := db.Where("resource = ? AND action = ? AND scope = ?", perm.Resource, perm.Action, perm.Scope).First(&permission).Error
		if err != nil {
			// Si no existe, crear el permiso
			permission = perm
			db.Create(&permission)
		}

		// Asignar permiso al rol usando la asociación many2many
		err = db.Model(&employeeRole).Association("Permissions").Append(&permission)
		require.NoError(t, err)
	}

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

	// 6. Recargar membership con role para el JWT
	err = db.Preload("Role").Where("id = ?", membership.ID).First(membership).Error
	require.NoError(t, err)

	return user, org, membership
}

// generateFreshToken genera un token JWT fresco para testing
func generateFreshToken(t *testing.T, user *models.Identity, membership *models.OrganizationalMembership) string {
	jwtService := utils.NewJWTService()

	// Preparar memberships slice para el JWT
	memberships := []models.OrganizationalMembership{*membership}

	accessToken, _, err := jwtService.GenerateTokenPair(user, memberships, nil)
	require.NoError(t, err)

	return accessToken
}

func generateTestID() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}

func TestGastosMiddlewareIntegration(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	user, org, membership := setupTestUser(t, db)

	// ⚠️ HACK TEMPORAL: Inicializar database.DBconn para evitar nil pointer
	// En una implementación real, el middleware debería ser más flexible
	originalDBconn := database.DBconn
	database.DBconn = db
	defer func() {
		database.DBconn = originalDBconn
	}()

	// Crear app Fiber
	app := fiber.New()

	// Configurar rutas de gastos
	gastos.SetupExpenseRoutes(app, db)

	t.Run("Debug - Verificar token JWT", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		t.Logf("Token generado: %s", token)
		t.Logf("Usuario ID: %s", user.ID)
		t.Logf("Organización ID: %s", org.ID)

		// Verificar que el token se puede validar
		jwtService := utils.NewJWTService()
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			t.Logf("Error validando token: %v", err)
		} else {
			t.Logf("Claims válidos: IdentityID=%s, Email=%s", claims.IdentityID, claims.Email)
		}
	})

	t.Run("Acceso sin autenticación debería fallar", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/org/test-gastos/expenses", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Acceso con token válido debería funcionar", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		req := httptest.NewRequest("GET", "/api/org/test-gastos/expenses", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST crear gasto con permisos", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		expense := map[string]interface{}{
			"title":       "Test Expense",
			"description": "Gasto de prueba",
			"amount":      100.50,
		}

		body, _ := json.Marshal(expense)
		req := httptest.NewRequest("POST", "/api/org/test-gastos/expenses", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Acceso a dashboard de usuario", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		req := httptest.NewRequest("GET", "/api/org/test-gastos/user/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Acceso admin sin permisos debería fallar", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		req := httptest.NewRequest("GET", "/api/org/test-gastos/admin/expenses/all", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Verificar contexto de usuario en responses", func(t *testing.T) {
		token := generateFreshToken(t, user, membership)
		req := httptest.NewRequest("GET", "/api/org/test-gastos/user/expenses", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verificar que la respuesta contiene datos del usuario
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.Equal(t, "success", response["status"])
	})

	fmt.Printf("✅ Test de integración del middleware de gastos exitoso\n")
	fmt.Printf("   Usuario: %s (%s)\n", user.Email, user.ID)
	fmt.Printf("   Organización: %s (%s)\n", org.Name, org.ID)
}
