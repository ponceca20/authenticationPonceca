package users

import (
	"fmt"
	"time"

	"practicev2/database"
)

func Seeders() bool {
	// Inicializar la conexión a MySQL usando el database.go
	database.ConnectDatabase()
	db := database.DBconn

	// Auto-migrate de modelos
	if err := db.AutoMigrate(
		&Persona{},
		&Usuario{},
		&UsuarioEmpresa{},
		&Rol{},
		&Sesion{},
		&Modulo{},
		&RolModulo{},
	); err != nil {
		fmt.Println("Error during auto migrate:", err)
		return false
	}

	// Insertar datos de ejemplo

	// Crear Persona
	persona := Persona{
		DocumentoNumero: "12345678",
		Foto:            "foto.jpg",
		Nombre:          "Juan",
		Apellidos:       "Perez",
		Email:           "juan@example.com",
		Telefono:        "123456789",
		Direccion:       "Calle Falsa 123",
		Ciudad:          "Ciudad",
		Pais:            "Pais",
		FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if result := db.Create(&persona); result.Error != nil {
		fmt.Println("Error creating Persona:", result.Error)
		return false
	}

	// Crear Usuario
	usuario := Usuario{
		PasswordHash:   "hashedpassword",
		Activo:         true,
		CreadoPor:      1,
		ActualizadoPor: 1,
	}
	if result := db.Create(&usuario); result.Error != nil {
		fmt.Println("Error creating Usuario:", result.Error)
		return false
	}

	// Crear Rol
	rol := Rol{
		Codigo: "ADMIN",
		Nombre: "Administrador",
		Activo: true,
	}
	if result := db.Create(&rol); result.Error != nil {
		fmt.Println("Error creating Rol:", result.Error)
		return false
	}

	// Crear UsuarioEmpresa
	usuarioEmpresa := UsuarioEmpresa{
		EmpresaID:       1,
		RolID:           uint64(rol.ID),
		FechaAsignacion: time.Now(),
	}
	if result := db.Create(&usuarioEmpresa); result.Error != nil {
		fmt.Println("Error creating UsuarioEmpresa:", result.Error)
		return false
	}

	// Crear Sesion
	sesion := Sesion{
		Token:           "token_example",
		RefreshToken:    "refresh_example",
		FechaExpiracion: time.Now().Add(24 * time.Hour),
		Activa:          true,
	}
	if result := db.Create(&sesion); result.Error != nil {
		fmt.Println("Error creating Sesion:", result.Error)
		return false
	}

	// Crear Modulo
	modulo := Modulo{
		Nombre: "Modulo 1",
		Ruta:   "/modulo1",
		Activo: true,
		Orden:  1,
	}
	if result := db.Create(&modulo); result.Error != nil {
		fmt.Println("Error creating Modulo:", result.Error)
		return false
	}

	// Crear RolModulo
	rolModulo := RolModulo{
		ModuloID: uint64(modulo.ID),
		Acceso:   true,
	}
	if result := db.Create(&rolModulo); result.Error != nil {
		fmt.Println("Error creating RolModulo:", result.Error)
		return false
	}

	fmt.Println("Seed data inserted successfully!")
	return true
}
