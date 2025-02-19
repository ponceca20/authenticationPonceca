// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/cmd/main.go
package main

import (
	//modulos importados
	_ "practicev2/modules"
	//fin de modulos importados
	"log"
	"practicev2/database"

	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDatabase()
	// Ejecuta las migraciones de forma secuencial.
	registry.RunMigrations(database.DBconn)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Bienvenido") })
	registry.RegisterAllModules(app)
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
