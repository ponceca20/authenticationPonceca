package addModules

import (
	"practicev2/module/authentication"

	"github.com/gofiber/fiber/v2"
)

// RegisterAllModules registra las rutas de todos los módulos
func RegisterAllModules(app *fiber.App) {
	// Registrar módulo de autenticación
	authentication.RegisterRoutes(app)

	// TODO: Aquí se registrarán otros módulos cuando se implementen
	// billing.RegisterRoutes(app)
	// inventory.RegisterRoutes(app)
	// notifications.RegisterRoutes(app)
}
