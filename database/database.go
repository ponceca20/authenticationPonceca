package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	// Construir el DSN para MySQL: usuario:contraseña@tcp(ruta:puerto)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Desactiva la pluralización
		},
	})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos MySQL: %v", err)
	}
	DBconn = db

	fmt.Println("Conexión a MySQL establecida correctamente y migración completada.")
}
