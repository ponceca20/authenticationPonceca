// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/cmd/main.go
package main

import (
	//modulos importados
	_ "practicev2/addModules"
	//fin de modulos importados
	"log"
	"os"
	"practicev2/database"

	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Ajusta la ruta al archivo .env ya que el working directory es /cmd
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No se pudo cargar el archivo .env, revisa la ruta o ejecuta desde el directorio raíz")
	}
	database.ConnectDatabase()
	// Ejecuta las migraciones de forma secuencial.
	registry.RunMigrations(database.DBconn)
	app := fiber.New()

	// Registrar CORS con configuración actualizada para solicitudes con credenciales
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Bienvenido") })
	registry.RegisterAllModules(app)

	// Obtiene la IP y el Puerto desde el .env
	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	log.Fatal(app.Listen(addr))
}
