package auth

import (
	"github.com/gofiber/fiber/v2"

	// ...existing imports if needed...
	"practicev2/registry"
)

// RegisterCentralRoutes centraliza la configuración de middleware y el registro de rutas.
func RegisterCentralRoutes(app *fiber.App) {
	// Middleware centralizado: ejemplo de logging

	// Llamadas a la configuración de rutas de cada handler
	RegisterRoutes2(app)
	RegisterRoutesUsuarioEmpresa(app)
	RegisterRoutesSesion(app)
	RegisterRoutesRolModulo(app)
	RegisterRoutesRol(app)
	RegisterRoutesPersona(app)
	RegisterRoutesModulo(app)
	RegisterRutasComplejas(app)

}

// Añadido init para registrar este módulo en el registry.
func init() {
	registry.RegisterModule(RegisterCentralRoutes)
}
