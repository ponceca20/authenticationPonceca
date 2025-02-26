package users

import (
	"log"
	"os"
	"strings"
	"time"

	"practicev2/database"
	"practicev2/registry" // Agregado para registrar las rutas

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ------------------------
// Helpers y definiciones
// ------------------------

// Elimina la variable global jwtKey y agrega una función que obtiene la clave en tiempo de ejecución.
func getJWTKey() []byte {
	secret := os.Getenv("JWT_KEY")
	if secret == "" {
		log.Fatal("JWT_KEY no configurado en el archivo .env")
	}
	return []byte(secret)
}

// LoginRequest representa la carga útil del login.
type LoginRequest struct {
	Identifier string `json:"identifier"` // DNI o celular
	Password   string `json:"password"`
}

// ModuleAccess representa los módulos a los que tiene acceso el rol.
type ModuleAccess struct {
	ModuloID uint64 `json:"modulo_id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
}

// RoleAccess representa el rol dentro de una empresa y sus módulos.
type RoleAccess struct {
	RolID   uint64         `json:"rol_id"`
	Codigo  string         `json:"codigo"`
	Nombre  string         `json:"nombre"`
	Modules []ModuleAccess `json:"modules"`
}

// CompanyAccess representa la asociación de un usuario con una empresa.
type CompanyAccess struct {
	EmpresaID uint64     `json:"empresa_id"`
	Role      RoleAccess `json:"role"`
}

// TokenClaims define los claims personalizados.
type TokenClaims struct {
	UsuarioID uint64          `json:"usuario_id"`
	PersonaID uint64          `json:"persona_id"`
	Companies []CompanyAccess `json:"companies"`
	jwt.RegisteredClaims
}

// ------------------------
// Helpers de acceso a la BD
// ------------------------

// Helper: Obtener módulos para un rol dado.
func getModulesForRole(rolID uint) []ModuleAccess {
	var modules []ModuleAccess
	var rolModulos []RolModulo
	if err := database.DBconn.Where("rol_id = ? AND acceso = ?", rolID, true).Find(&rolModulos).Error; err != nil {
		return modules
	}
	for _, rm := range rolModulos {
		var modulo Modulo
		if err := database.DBconn.First(&modulo, rm.ModuloID).Error; err != nil {
			continue
		}
		modules = append(modules, ModuleAccess{
			ModuloID: uint64(modulo.ID),
			Codigo:   modulo.Codigo,
			Nombre:   modulo.Nombre,
		})
	}
	return modules
}

// Helper: Obtener accesos a empresas para un ID de persona.
func getCompanyAccesses(personaID uint) []CompanyAccess {
	var accesos []CompanyAccess
	var ueList []UsuarioEmpresa
	if err := database.DBconn.Where("usuario_id = ?", personaID).Find(&ueList).Error; err != nil {
		return accesos
	}
	for _, ue := range ueList {
		var rol Rol
		if err := database.DBconn.First(&rol, ue.RolID).Error; err != nil {
			continue
		}
		accesos = append(accesos, CompanyAccess{
			EmpresaID: ue.EmpresaID,
			Role: RoleAccess{
				RolID:   uint64(rol.ID),
				Codigo:  rol.Codigo,
				Nombre:  rol.Nombre,
				Modules: getModulesForRole(rol.ID),
			},
		})
	}
	return accesos
}

// ------------------------
// Repositorio
// Encargado de interactuar con la base de datos.
// ------------------------

// loginRepository busca la persona y el usuario activo correspondientes al identificador.
func loginRepository(identifier string) (Persona, Usuario, error) {
	var persona Persona
	if err := database.DBconn.Where("documento_numero = ? OR telefono = ?", identifier, identifier).First(&persona).Error; err != nil {
		return persona, Usuario{}, err
	}

	var usuario Usuario
	if err := database.DBconn.Where("persona_id = ? AND activo = ?", persona.ID, true).First(&usuario).Error; err != nil {
		return persona, usuario, err
	}
	return persona, usuario, nil
}

// ------------------------
// Servicio
// Encargado de la lógica de negocio del login.
// ------------------------

func loginService(req LoginRequest) (fiber.Map, error) {
	// Llamada al repositorio para obtener el usuario y la persona.
	persona, usuario, err := loginRepository(req.Identifier)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Credenciales inválidas")
	}

	// Validar la contraseña utilizando bcrypt.
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Credenciales inválidas")
	}

	// Obtener el tiempo actual y configuramos la expiración del token.
	now := time.Now()
	expirationTime := now.Add(24 * time.Hour)

	// Obtener accesos a las empresas.
	accesos := getCompanyAccesses(persona.ID)

	// Configurar los claims del token JWT.
	claims := TokenClaims{
		UsuarioID: uint64(usuario.ID),
		PersonaID: uint64(persona.ID),
		Companies: accesos,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "myapp",
		},
	}

	// Generar y firmar el token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error generando token")
	}

	// Registrar la sesión:
	session := Sesion{
		UsuarioID:       uint64(usuario.ID),
		Token:           tokenString,
		FechaExpiracion: expirationTime,
		Activa:          true,
	}
	if err = database.DBconn.Create(&session).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error creando sesión")
	}

	// Construir la respuesta.
	response := fiber.Map{
		"token":      tokenString,
		"expires_at": expirationTime.Format(time.RFC3339),
		"user": fiber.Map{
			"usuario_id": usuario.ID,
			"persona_id": persona.ID,
		},
		"companies": accesos,
	}
	return response, nil
}

// ------------------------
// Handler
// Encargado de manejar la petición HTTP.
// ------------------------

func loginHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json") // Forzar respuesta en JSON
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Solicitud inválida"})
	}
	if req.Identifier == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El identificador y la contraseña son requeridos"})
	}

	// Llamar al servicio para realizar el login.
	response, err := loginService(req)
	if err != nil {
		// Si el error es del servicio (ya contiene el status adecuado), se propaga.
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error interno"})
	}

	return c.JSON(response)
}

func logoutHandler(c *fiber.Ctx) error {
	tokenStr := c.Get("Authorization")
	if tokenStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token no provisto"})
	}

	tokenStr = strings.TrimSpace(strings.TrimPrefix(tokenStr, "Bearer"))

	var sesion Sesion
	if err := database.DBconn.Where("token = ? AND activa = ?", tokenStr, true).First(&sesion).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Sesión no encontrada o ya inactiva"})
	}

	sesion.Activa = false
	if err := database.DBconn.Save(&sesion).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al cerrar sesión"})
	}

	return c.JSON(fiber.Map{"message": "Sesión cerrada exitosamente"})
}

// ------------------------
// Registro de rutas
// ------------------------

// RegisterAdvancedLogin registra la ruta de login avanzada.
func RegisterAdvancedLogin(app *fiber.App) {
	app.Post("/login", loginHandler)
	app.Post("/logout", logoutHandler)
}

func init() {
	// Registrar el módulo de login avanzado en el registry.
	// Se asume que registry.RegisterModule es el mecanismo para cargar rutas.
	registry.RegisterModule(RegisterAdvancedLogin)
}
