package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	// DBconn es la variable global que almacena la conexión a la base de datos.
	DBconn *gorm.DB
)

// ConnectDatabase se conecta a MySQL utilizando GORM y configura la variable global DBconn.
// Se recomienda cargar el DSN desde variables de entorno en producción.
func ConnectDatabase() {
	// DSN para MySQL: usuario:contraseña@tcp(ruta)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:toor@tcp(127.0.0.1:3306)/practicav1?charset=utf8mb4&parseTime=True&loc=Local"
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
