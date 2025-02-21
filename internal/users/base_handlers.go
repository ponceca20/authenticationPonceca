package users

import (
	"practicev2/registry"

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

// RegisterRutasComplejas registers complex API routes.
func RegisterRutasComplejas(app *fiber.App) {
	api := app.Group("/api/complejas")
	api.Get("/seed", ExecuteSeed)
}

func init() {
	registry.RegisterModule(RegisterRutasComplejas)
}
