package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	// DBconn es la variable global que almacena la conexión a la base de datos.
	DBconn *gorm.DB
)

// ConnectDatabase se conecta a MySQL utilizando GORM y configura la variable global DBconn.
// Los parámetros de conexión se cargan desde variables de entorno.
func ConnectDatabase() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		log.Fatal("Variables de entorno de base de datos incompletas")
	}

	// Construir el DSN para MySQL con configuraciones adicionales para estabilidad
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s",
		user, password, host, port, name)

	// Configurar logger personalizado para reducir logs innecesarios
	var gormLogger logger.Interface

	// Solo en desarrollo mostrar queries con errores y warnings
	if os.Getenv("ENV") == "development" {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Solo queries lentas
				LogLevel:                  logger.Warn,            // Solo warnings y errores
				IgnoreRecordNotFoundError: true,                   // Ignorar errores de "not found"
				Colorful:                  true,                   // Colores en desarrollo
			},
		)
	} else {
		// En producción, solo errores críticos
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             500 * time.Millisecond, // Solo queries muy lentas
				LogLevel:                  logger.Error,           // Solo errores
				IgnoreRecordNotFoundError: true,                   // Ignorar errores de "not found"
				Colorful:                  false,                  // Sin colores en producción
			},
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Desactiva la pluralización
		},
		Logger: gormLogger,
		// Configuraciones para mejor rendimiento y estabilidad
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos MySQL: %v", err)
	}

	// Configurar pool de conexiones para mejor rendimiento
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error obteniendo la instancia de base de datos SQL: %v", err)
	}

	// Configuraciones del pool de conexiones
	sqlDB.SetMaxIdleConns(10)                  // Máximo número de conexiones inactivas
	sqlDB.SetMaxOpenConns(100)                 // Máximo número de conexiones abiertas
	sqlDB.SetConnMaxLifetime(time.Hour)        // Tiempo máximo de vida de una conexión
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Tiempo máximo que una conexión puede estar inactiva

	// Verificar la conexión
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error verificando la conexión a la base de datos: %v", err)
	}

	DBconn = db
	fmt.Println("Conexión a MySQL establecida correctamente y migración completada.")
}
