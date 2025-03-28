package auth

import (
	"fmt"
	"time"

	"practicev2/database"

	"golang.org/x/crypto/bcrypt"
)

func Seeders() bool {
	// Inicializar la conexión a MySQL usando el database.go
	database.ConnectDatabase()
	db := database.DBconn

	// Auto-migrate de modelos
	if err := db.AutoMigrate(
		&TipoDocumentoSunat{},
		&Persona{},
		&Usuario{},
		&UsuarioEmpresa{},
		&Rol{},
		&Sesion{},
		&Modulo{},
		&RolModulo{},
		&LoginAttempt{},
	); err != nil {
		fmt.Println("Error during auto migrate:", err)
		return false
	}

	// Insertar datos de ejemplo para TipoDocumentoSunat
	dni := TipoDocumentoSunat{
		Codigo: "01",
		Nombre: "Documento Nacional de Identidad",
	}
	if result := db.FirstOrCreate(&dni, TipoDocumentoSunat{Codigo: "01"}); result.Error != nil {
		fmt.Println("Error creating TipoDocumentoSunat DNI:", result.Error)
		return false
	}

	ruc := TipoDocumentoSunat{
		Codigo: "06",
		Nombre: "Registro Único de Contribuyentes",
	}
	if result := db.FirstOrCreate(&ruc, TipoDocumentoSunat{Codigo: "06"}); result.Error != nil {
		fmt.Println("Error creating TipoDocumentoSunat RUC:", result.Error)
		return false
	}

	// Crear Persona
	persona1 := Persona{
		TipoDocumentoID:    uint64(dni.ID),
		DocumentoNumero:    "12345678",
		Foto:               "foto1.jpg",
		Nombre:             "Juan",
		Apellidos:          "Perez",
		Email:              "juan@example.com",
		Telefono:           "123456789",
		TelefonoSecundario: "987654321",
		Direccion:          "Calle Falsa 123",
		FechaNacimiento:    time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	if result := db.FirstOrCreate(&persona1, Persona{DocumentoNumero: "12345678"}); result.Error != nil {
		fmt.Println("Error creating Persona 1:", result.Error)
		return false
	}

	persona2 := Persona{
		TipoDocumentoID: uint64(dni.ID),
		DocumentoNumero: "41822932",
		Foto:            "foto2.jpg",
		Nombre:          "Fredy",
		Apellidos:       "Ponceca",
		Email:           "fredy@example.com",
		Telefono:        "987654321",
		Direccion:       "Avenida Siempre Viva 742",
		FechaNacimiento: time.Date(1985, time.June, 15, 0, 0, 0, 0, time.UTC),
	}
	if result := db.FirstOrCreate(&persona2, Persona{DocumentoNumero: "41822932"}); result.Error != nil {
		fmt.Println("Error creating Persona 2:", result.Error)
		return false
	}

	// Generar password hashes
	password1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	password2, _ := bcrypt.GenerateFromPassword([]byte("secret456"), bcrypt.DefaultCost)

	// Crear Usuarios
	usuario1 := Usuario{
		PersonaID:      uint64(persona1.ID),
		PasswordHash:   string(password1),
		Activo:         true,
		CreadoPor:      1,
		ActualizadoPor: 1,
	}
	if result := db.FirstOrCreate(&usuario1, Usuario{PersonaID: uint64(persona1.ID)}); result.Error != nil {
		fmt.Println("Error creating Usuario 1:", result.Error)
		return false
	}

	usuario2 := Usuario{
		PersonaID:      uint64(persona2.ID),
		PasswordHash:   string(password2),
		Activo:         true,
		CreadoPor:      1,
		ActualizadoPor: 1,
	}
	if result := db.FirstOrCreate(&usuario2, Usuario{PersonaID: uint64(persona2.ID)}); result.Error != nil {
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
	if result := db.FirstOrCreate(&rolAdmin, Rol{Codigo: "ADMIN", EmpresaID: 1}); result.Error != nil {
		fmt.Println("Error creating Rol Admin:", result.Error)
		return false
	}

	rolUsuario := Rol{
		EmpresaID: 1,
		Codigo:    "USER",
		Nombre:    "Usuario Regular",
		Activo:    true,
	}
	if result := db.FirstOrCreate(&rolUsuario, Rol{Codigo: "USER", EmpresaID: 1}); result.Error != nil {
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
	if result := db.FirstOrCreate(&usuarioEmpresa1, UsuarioEmpresa{
		UsuarioID: uint64(usuario1.ID),
		EmpresaID: 1,
	}); result.Error != nil {
		fmt.Println("Error creating UsuarioEmpresa 1:", result.Error)
		return false
	}

	usuarioEmpresa2 := UsuarioEmpresa{
		UsuarioID:       uint64(usuario2.ID),
		EmpresaID:       1,
		RolID:           uint64(rolUsuario.ID),
		FechaAsignacion: time.Now(),
	}
	if result := db.FirstOrCreate(&usuarioEmpresa2, UsuarioEmpresa{
		UsuarioID: uint64(usuario2.ID),
		EmpresaID: 1,
	}); result.Error != nil {
		fmt.Println("Error creating UsuarioEmpresa 2:", result.Error)
		return false
	}

	// Crear Módulos
	modulos := []Modulo{
		{Codigo: "DASHBOARD", Nombre: "Dashboard", Ruta: "/dashboard", Activo: true, Orden: 1},
		{Codigo: "USERS", Nombre: "Usuarios", Ruta: "/users", Activo: true, Orden: 2},
		{Codigo: "REPORTS", Nombre: "Reportes", Ruta: "/reports", Activo: true, Orden: 3},
		{Codigo: "CONFIG", Nombre: "Configuración", Ruta: "/config", Activo: true, Orden: 4},
	}

	for i, m := range modulos {
		if result := db.FirstOrCreate(&modulos[i], Modulo{Codigo: m.Codigo}); result.Error != nil {
			fmt.Printf("Error creating Modulo %s: %v\n", m.Codigo, result.Error)
			return false
		}
	}

	// Crear RolModulos (permisos)
	// Para rol Admin (tiene acceso a todos los módulos)
	for _, modulo := range modulos {
		rolModulo := RolModulo{
			RolID:    uint64(rolAdmin.ID),
			ModuloID: uint64(modulo.ID),
			Acceso:   true,
		}
		if result := db.FirstOrCreate(&rolModulo, RolModulo{
			RolID:    uint64(rolAdmin.ID),
			ModuloID: uint64(modulo.ID),
		}); result.Error != nil {
			fmt.Printf("Error creating RolModulo for Admin-%s: %v\n", modulo.Codigo, result.Error)
			return false
		}
	}

	// Para rol Usuario (solo acceso a Dashboard y Reportes)
	userModulos := []string{"DASHBOARD", "REPORTS"}
	for _, codigo := range userModulos {
		var modulo Modulo
		if result := db.Where("codigo = ?", codigo).First(&modulo); result.Error != nil {
			fmt.Printf("Error finding Modulo %s: %v\n", codigo, result.Error)
			return false
		}

		rolModulo := RolModulo{
			RolID:    uint64(rolUsuario.ID),
			ModuloID: uint64(modulo.ID),
			Acceso:   true,
		}
		if result := db.FirstOrCreate(&rolModulo, RolModulo{
			RolID:    uint64(rolUsuario.ID),
			ModuloID: uint64(modulo.ID),
		}); result.Error != nil {
			fmt.Printf("Error creating RolModulo for User-%s: %v\n", modulo.Codigo, result.Error)
			return false
		}
	}

	// Crear ejemplo de sesiones
	ahora := time.Now()
	sesiones := []Sesion{
		{
			UsuarioID:       uint64(usuario1.ID),
			Token:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3VhcmlvX2lkIjoxLCJwZXJzb25hX2lkIjoxfQ.ejemplo1",
			RefreshToken:    "refresh_token_1",
			FechaExpiracion: ahora.Add(24 * time.Hour),
			Activa:          true,
			IP:              "192.168.1.100",
		},
		{
			UsuarioID:       uint64(usuario2.ID),
			Token:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3VhcmlvX2lkIjoyLCJwZXJzb25hX2lkIjoyfQ.ejemplo2",
			RefreshToken:    "refresh_token_2",
			FechaExpiracion: ahora.Add(24 * time.Hour),
			Activa:          false, // sesión inactiva
			IP:              "192.168.1.101",
		},
	}

	for _, sesion := range sesiones {
		// No usamos FirstOrCreate aquí porque las sesiones son únicas y pueden repetirse
		if result := db.Create(&sesion); result.Error != nil {
			fmt.Printf("Error creating Sesion for Usuario %d: %v\n", sesion.UsuarioID, result.Error)
			return false
		}
	}

	// Crear ejemplos de intentos de login
	attempts := []LoginAttempt{
		{
			Identifier: "12345678", // usando DNI
			IP:         "192.168.1.100",
			Success:    true,
		},
		{
			Identifier: "invalid@example.com",
			IP:         "192.168.1.200",
			Success:    false,
		},
		{
			Identifier: "41822932", // DNI de Fredy
			IP:         "192.168.1.150",
			Success:    true,
		},
	}

	for _, attempt := range attempts {
		if result := db.Create(&attempt); result.Error != nil {
			fmt.Printf("Error creating LoginAttempt for %s: %v\n", attempt.Identifier, result.Error)
			return false
		}
	}

	fmt.Println("Seed data inserted successfully!")
	return true
}
