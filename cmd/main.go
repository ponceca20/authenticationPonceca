// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/cmd/main.go
package main

import (
	//modulos importados

	"os"
	_ "practicev2/addModules"
	"practicev2/config"
	"practicev2/module/middleware"

	"github.com/joho/godotenv"

	//fin de modulos importados
	"log"
	"practicev2/database"
	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando archivo .env")
	}
	config.Init()
	database.ConnectDatabase()
	// Ejecuta las migraciones de forma secuencial.
	registry.RunMigrations(database.DBconn)
	app := fiber.New()

	// Actualización de configuración CORS mejorada para mejor compatibilidad con HTTPS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost:3000,http://localhost:3000,https://192.168.1.36:3000,http://192.168.1.36:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowCredentials: true,
		MaxAge:           3600, // Aumentar el tiempo de caché de preflight para evitar problemas
	}))
	app.Use(middleware.ConfigMiddleware())

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Bienvenido") })
	registry.RegisterAllModules(app)

	//Usar la dirección del servidor desde el paquete config
	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	log.Fatal(app.Listen(addr))
}
