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
		&LoginAttempt{}, // Nuevo modelo agregado
	); err != nil {
		fmt.Println("Error during auto migrate:", err)
		return false
	}

	// Insertar datos de ejemplo

	// Crear Persona
	persona1 := Persona{
		DocumentoTipo:   "DNI",
		DocumentoNumero: "12345678",
		Foto:            "foto1.jpg",
		Nombre:          "Juan",
		Apellidos:       "Perez",
		Email:           "juan@example.com",
		Telefono:        "123456789",
		Direccion:       "Calle Falsa 123",
		Ciudad:          "Lima",
		Pais:            "Peru",
		FechaNacimiento: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if result := db.Create(&persona1); result.Error != nil {
		fmt.Println("Error creating Persona 1:", result.Error)
		return false
	}

	persona2 := Persona{
		DocumentoTipo:   "DNI",
		DocumentoNumero: "87654321",
		Foto:            "foto2.jpg",
		Nombre:          "Maria",
		Apellidos:       "Lopez",
		Email:           "maria@example.com",
		Telefono:        "987654321",
		Direccion:       "Avenida Siempre Viva 742",
		Ciudad:          "Bogota",
		Pais:            "Colombia",
		FechaNacimiento: time.Date(1985, time.June, 15, 0, 0, 0, 0, time.UTC),
	}
	if result := db.Create(&persona2); result.Error != nil {
		fmt.Println("Error creating Persona 2:", result.Error)
		return false
	}

	// Crear Usuarios
	usuario1 := Usuario{
		PersonaID:      uint64(persona1.ID),
		PasswordHash:   "$2a$10$xJJJ4O.UhJ1nRlbmtJ1y0ujh0C5q/YGGpBnLCFpZ9icWZ5bQ7fUTa", // hash para "password123"
		Activo:         true,
		CreadoPor:      1,
		ActualizadoPor: 1,
	}
	if result := db.Create(&usuario1); result.Error != nil {
		fmt.Println("Error creating Usuario 1:", result.Error)
		return false
	}

	usuario2 := Usuario{
		PersonaID:      uint64(persona2.ID),
		PasswordHash:   "$2a$10$k4HT2KkfCPv1/zhDfg1NoeEkV86HMjHvM3ndCGGAiRwBoIppyizEy", // hash para "secret456"
		Activo:         true,
		CreadoPor:      1,
		ActualizadoPor: 1,
	}
	if result := db.Create(&usuario2); result.Error != nil {
		fmt.Println("Error creating Usuario 2:", result.Error)
		return false
	}

	// Crear Roles
	rolAdmin := Rol{
		EmpresaID: 1,
		Codigo:    "ADMIN",
		Nombre:    "Administrador",
		Activo:    true,
	}
	if result := db.Create(&rolAdmin); result.Error != nil {
		fmt.Println("Error creating Rol Admin:", result.Error)
		return false
	}

	rolUsuario := Rol{
		EmpresaID: 1,
		Codigo:    "USER",
		Nombre:    "Usuario Regular",
		Activo:    true,
	}
	if result := db.Create(&rolUsuario); result.Error != nil {
		fmt.Println("Error creating Rol Usuario:", result.Error)
		return false
	}

	// Crear UsuarioEmpresa relaciones
	usuarioEmpresa1 := UsuarioEmpresa{
		UsuarioID:       uint64(usuario1.ID),
		EmpresaID:       1,
		RolID:           uint64(rolAdmin.ID),
		FechaAsignacion: time.Now(),
	}
	if result := db.Create(&usuarioEmpresa1); result.Error != nil {
		fmt.Println("Error creating UsuarioEmpresa 1:", result.Error)
		return false
	}

	usuarioEmpresa2 := UsuarioEmpresa{
		UsuarioID:       uint64(usuario2.ID),
		EmpresaID:       1,
		RolID:           uint64(rolUsuario.ID),
		FechaAsignacion: time.Now(),
	}
	if result := db.Create(&usuarioEmpresa2); result.Error != nil {
		fmt.Println("Error creating UsuarioEmpresa 2:", result.Error)
		return false
	}

	// Crear Sesiones
	sesion1 := Sesion{
		UsuarioID:       uint64(usuario1.ID),
		Token:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
		RefreshToken:    "refresh_token_1_example",
		FechaExpiracion: time.Now().Add(24 * time.Hour),
		Activa:          true,
		IP:              "192.168.1.100",
	}
	if result := db.Create(&sesion1); result.Error != nil {
		fmt.Println("Error creating Sesion 1:", result.Error)
		return false
	}

	sesion2 := Sesion{
		UsuarioID:       uint64(usuario2.ID),
		Token:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI5ODc2NTQzMjEwIiwibmFtZSI6Ik1hcmlhIExvcGV6IiwiaWF0IjoxNTE2MjM5MDIyfQ",
		RefreshToken:    "refresh_token_2_example",
		FechaExpiracion: time.Now().Add(24 * time.Hour),
		Activa:          true,
		IP:              "192.168.1.101",
	}
	if result := db.Create(&sesion2); result.Error != nil {
		fmt.Println("Error creating Sesion 2:", result.Error)
		return false
	}

	// Crear Módulos
	modulo1 := Modulo{
		Codigo: "DASHBOARD",
		Nombre: "Dashboard",
		Ruta:   "/dashboard",
		Activo: true,
		Orden:  1,
	}
	if result := db.Create(&modulo1); result.Error != nil {
		fmt.Println("Error creating Modulo 1:", result.Error)
		return false
	}

	modulo2 := Modulo{
		Codigo: "USERS",
		Nombre: "Usuarios",
		Ruta:   "/users",
		Activo: true,
		Orden:  2,
	}
	if result := db.Create(&modulo2); result.Error != nil {
		fmt.Println("Error creating Modulo 2:", result.Error)
		return false
	}

	modulo3 := Modulo{
		Codigo: "REPORTS",
		Nombre: "Reportes",
		Ruta:   "/reports",
		Activo: true,
		Orden:  3,
	}
	if result := db.Create(&modulo3); result.Error != nil {
		fmt.Println("Error creating Modulo 3:", result.Error)
		return false
	}

	// Crear RolModulos (permisos)
	// Admin tiene acceso a todos los módulos
	rolModulo1 := RolModulo{
		RolID:    uint64(rolAdmin.ID),
		ModuloID: uint64(modulo1.ID),
		Acceso:   true,
	}
	if result := db.Create(&rolModulo1); result.Error != nil {
		fmt.Println("Error creating RolModulo 1:", result.Error)
		return false
	}

	rolModulo2 := RolModulo{
		RolID:    uint64(rolAdmin.ID),
		ModuloID: uint64(modulo2.ID),
		Acceso:   true,
	}
	if result := db.Create(&rolModulo2); result.Error != nil {
		fmt.Println("Error creating RolModulo 2:", result.Error)
		return false
	}

	rolModulo3 := RolModulo{
		RolID:    uint64(rolAdmin.ID),
		ModuloID: uint64(modulo3.ID),
		Acceso:   true,
	}
	if result := db.Create(&rolModulo3); result.Error != nil {
		fmt.Println("Error creating RolModulo 3:", result.Error)
		return false
	}

	// Usuario regular solo tiene acceso al dashboard
	rolModulo4 := RolModulo{
		RolID:    uint64(rolUsuario.ID),
		ModuloID: uint64(modulo1.ID),
		Acceso:   true,
	}
	if result := db.Create(&rolModulo4); result.Error != nil {
		fmt.Println("Error creating RolModulo 4:", result.Error)
		return false
	}

	// Crear LoginAttempts
	loginAttempt1 := LoginAttempt{
		Identifier: "juan@example.com",
		IP:         "192.168.1.100",
		Success:    true,
	}
	if result := db.Create(&loginAttempt1); result.Error != nil {
		fmt.Println("Error creating LoginAttempt 1:", result.Error)
		return false
	}

	loginAttempt2 := LoginAttempt{
		Identifier: "unknown@example.com",
		IP:         "192.168.1.200",
		Success:    false,
	}
	if result := db.Create(&loginAttempt2); result.Error != nil {
		fmt.Println("Error creating LoginAttempt 2:", result.Error)
		return false
	}

	fmt.Println("Seed data inserted successfully!")
	return true
}
