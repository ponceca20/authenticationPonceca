// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/cmd/main.go
package main

import (
	//modulos importados
	"context"
	"os"
	"os/signal"
	"practicev2/addModules"
	"practicev2/config"
	_ "practicev2/module/authentication" // Import for side effects (init)
	"practicev2/module/middleware"
	"syscall"
	"time"

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

	// Inicializar configuración
	config.Init()

	// Conectar a la base de datos
	database.ConnectDatabase()

	// Ejecutar migraciones de forma segura
	if err := registry.RunMigrations(database.DBconn); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	// Configurar Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:             10 * 1024 * 1024, // 10MB en bytes
		DisableStartupMessage: false,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           60 * time.Second,
		ErrorHandler:          customErrorHandler,
	})

	// Configuración CORS mejorada
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost:3000,http://localhost:3000,https://192.168.1.36:3000,http://192.168.1.36:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// Middleware personalizado
	app.Use(middleware.ConfigMiddleware())

	// Ruta principal
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bienvenido al sistema de autenticación Ponceca")
	})

	// Registrar todos los módulos
	addModules.RegisterAllModules(app)

	// Configurar canal para señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Obtener dirección del servidor
	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")

	// Iniciar servidor en una goroutine
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Printf("Error iniciando servidor: %v", err)
		}
	}()

	log.Printf("🚀 Servidor iniciado en %s", addr)
	log.Println("💡 Presiona Ctrl+C para detener el servidor de forma segura")

	// Esperar señal de interrupción
	<-quit
	log.Println("🛑 Recibida señal de interrupción, iniciando cierre seguro...")

	// Crear contexto con timeout para el cierre
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cerrar servidor de forma segura
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("❌ Error durante el cierre del servidor: %v", err)
	} else {
		log.Println("✅ Servidor cerrado correctamente")
	}

	// Cerrar conexión a la base de datos
	if database.DBconn != nil {
		sqlDB, err := database.DBconn.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("❌ Error cerrando conexión a la base de datos: %v", err)
			} else {
				log.Println("✅ Conexión a la base de datos cerrada correctamente")
			}
		}
	}

	log.Println("👋 ¡Hasta luego!")
}

// customErrorHandler maneja los errores de forma personalizada
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Printf("❌ Error en %s %s: %v", c.Method(), c.Path(), err)

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}
