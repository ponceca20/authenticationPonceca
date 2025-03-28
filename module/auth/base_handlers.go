package auth

import (
	"github.com/gofiber/fiber/v2"
)

// ExecuteSeed performs the seeding process and returns a JSON response.
func ExecuteSeed(c *fiber.Ctx) error {
	if success := Seeders(); !success {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "seeder execution failed",
		})
	}
	return c.JSON(fiber.Map{
		"status": "seeder executed successfully",
	})
}

// RegisterRutasComplejas registra las rutas de prueba para las cookies.
func RegisterRutasComplejas(app *fiber.App) {
	group := app.Group("/api/complejas")
	group.Get("/seed", ExecuteSeed)

}
