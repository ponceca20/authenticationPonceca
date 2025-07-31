package middleware

import (
	"net/http"
	"practicev2/config"
	"practicev2/database"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAppWithMiddleware(db *gorm.DB) *fiber.App {
	app := fiber.New()
	database.DBconn = db // Override global DB connection for the test
	app.Use(SmartAuthMiddleware())
	return app
}

// setupTestDatabase initializes an in-memory SQLite database for testing purposes.
func setupTestDatabase(t *testing.T) *gorm.DB {
	// Initialize config with default values for tests
	config.Init()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Auto-migrate the models we need for testing
	err = db.AutoMigrate(
		&models.Identity{},
		&models.User{},
		&models.UserProfile{},
		&models.RefreshToken{},
	)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestSmartAuthMiddleware(t *testing.T) {
	// Setup Test DB and a test user
	db := setupTestDatabase(t)
	authRepo := auth.NewAuthRepository(db)
	jwtService := utils.NewJWTService()
	authService := auth.NewAuthService(authRepo, jwtService)

	registerDTO := &auth.RegisterDTO{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "secure-password-456",
	}
	identity, err := authService.Register(registerDTO)
	assert.NoError(t, err)

	// Generate a valid token for the test user
	accessToken, _, err := jwtService.GenerateTokenPair(identity, nil, nil)
	assert.NoError(t, err)

	app := setupAppWithMiddleware(db)

	// Test handler to inspect the context
	var capturedCtx *AuthContext
	app.Get("/protected", func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("authContext").(*AuthContext)
		if ok {
			capturedCtx = ctx
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// --- Test Cases ---

	t.Run("Valid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NotNil(t, capturedCtx)
		assert.False(t, capturedCtx.IsGuest)
		assert.Equal(t, identity.ID, capturedCtx.Identity.ID)
	})

	t.Run("No Token (Guest)", func(t *testing.T) {
		capturedCtx = nil // Reset captured context
		req, _ := http.NewRequest("GET", "/protected", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.NotNil(t, capturedCtx)
		assert.True(t, capturedCtx.IsGuest)
		assert.Nil(t, capturedCtx.Identity)
	})

	t.Run("Malformed Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer_invalid"+accessToken) // Malformed header
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken+"invalid") // Invalid signature part
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
