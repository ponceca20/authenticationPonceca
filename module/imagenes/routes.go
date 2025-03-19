package imagenes

import (
	"os"
	"practicev2/module/middleware"
	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
)

// Setup initializes the image module
func Setup() {
	registry.RegisterModule(RegisterRoutes)
}

// RegisterRoutes registers all routes for the image module
func RegisterRoutes(app *fiber.App) {
	// Create upload directory
	uploadDir := "./uploads/images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		// Agregamos verificación de error en la creación del directorio
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			panic("Unable to create upload directory: " + err.Error())
		}
	}
	base_api := "/api/images"
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	// Update la baseURL usando siempre APP_HOST y APP_PORT
	baseURLimg := host + ":" + port + base_api + "/serve"

	// Create repository, service and handler
	repo := NewImageRepository()
	service := NewImageService(repo, uploadDir, baseURLimg)
	handler := NewImageHandler(service)

	// Create API routes
	api := app.Group(base_api)

	// Public routes
	api.Get("/:id", handler.GetImage)
	api.Get("/serve/:filename", handler.ServeImage)
	api.Get("/user/:userId", handler.GetUserImages)

	// Protected routes
	api.Post("/upload", middleware.AuthMiddleware(), handler.UploadImage)
	api.Delete("/:id", middleware.AuthMiddleware(), handler.DeleteImage)

}

func init() {
	// Register this module with the application
	Setup()
}
