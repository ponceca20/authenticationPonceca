// Cambio el nombre del paquete para que coincida con la carpeta del módulo.
package auth

import (
	"log"
	"practicev2/config"
	"practicev2/database"
	"practicev2/registry"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ------------------------
// Helpers y definiciones
// ------------------------

// Constantes de configuración
const (
	// Changed cookieMaxAge to 30 days
	cookieMaxAge = 30 * 24 * time.Hour
	cookieName   = "auth_token"
)

// Obtiene la clave JWT desde el paquete config
func getJWTKey() []byte {
	return config.GetJWTKey()
}

// LoginRequest representa la carga útil del login.
type LoginRequest struct {
	Identifier string `json:"identifier"` // DNI o celular
	Password   string `json:"password"`
}

// Nuevo struct para representación mínima de empresa.
type MinimalCompany struct {
	EmpresaID uint64 `json:"empresa_id"`
}

// TokenClaims define los claims personalizados simplificados.
type TokenClaims struct {
	UsuarioID uint64           `json:"usuario_id"`
	PersonaID uint64           `json:"persona_id"`
	Companies []MinimalCompany `json:"companies"`
	jwt.RegisteredClaims
}

// ------------------------
// Helpers de acceso a la BD
// ------------------------

// Se modifica getCompanyAccesses para retornar solo MinimalCompany.
func getCompanyAccesses(usuarioID uint) []MinimalCompany {
	var companies []MinimalCompany
	var ueList []UsuarioEmpresa // ...existing struct definition...
	if err := database.DBconn.
		Select("empresa_id").
		Where("usuario_id = ? AND active_sesion = true", usuarioID). //
		Find(&ueList).Error; err != nil {
		return companies
	}
	for _, ue := range ueList {
		companies = append(companies, MinimalCompany{EmpresaID: ue.EmpresaID})
	}
	return companies
}

// loginRepository optimizado con una sola consulta
func loginRepository(identifier string) (Persona, Usuario, error) {
	var persona Persona
	var usuario Usuario

	// Primero, buscar la persona por documento o teléfono
	if err := database.DBconn.
		Where("documento_numero = ? OR telefono = ?", identifier, identifier).
		First(&persona).Error; err != nil {
		log.Printf("Error buscando persona: %v", err) // Registro para auditoría
		return persona, Usuario{}, err
	}

	// Luego, buscar el usuario activo asociado a la persona
	if err := database.DBconn.
		Where("persona_id = ? AND activo = ?", persona.ID, true).
		First(&usuario).Error; err != nil {
		log.Printf("Error buscando usuario activo para persona %d: %v", persona.ID, err)
		return persona, Usuario{}, err
	}

	return persona, usuario, nil
}

// Registrar un intento de login fallido
func registerLoginAttempt(identifier string, ip string) {
	// Guardar solo la información esencial para limitar intentos
	attempt := LoginAttempt{
		Identifier: identifier,
		IP:         ip,
		Success:    false,
	}
	database.DBconn.Create(&attempt)
}

// Verificar si hay demasiados intentos fallidos para este identificador
func checkLoginAttempts(identifier string, ip string) bool {
	var count int64
	timeWindow := time.Now().Add(-15 * time.Minute)

	// Contar intentos fallidos recientes
	database.DBconn.Model(&LoginAttempt{}).
		Where("identifier = ? AND ip = ? AND success = ? AND created_at > ?",
			identifier, ip, false, timeWindow).
		Count(&count)

	return count < 5 // Máximo 5 intentos en 15 minutos
}

// ------------------------
// Servicio
// Encargado de la lógica de negocio del login.
// ------------------------

func loginService(req LoginRequest, ip string) (fiber.Map, string, error) {
	// Verificar intentos de login
	if !checkLoginAttempts(req.Identifier, ip) {
		return nil, "", fiber.NewError(fiber.StatusTooManyRequests, "Demasiados intentos fallidos. Intente más tarde.")
	}

	// Llamada al repositorio para obtener el usuario y la persona.
	persona, usuario, err := loginRepository(req.Identifier)
	if err != nil {
		registerLoginAttempt(req.Identifier, ip)
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Credenciales inválidas")
	}

	// Validar la contraseña utilizando bcrypt.
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(req.Password)); err != nil {
		registerLoginAttempt(req.Identifier, ip)
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Credenciales inválidas")
	}

	// Obtener el tiempo actual y configuramos la expiración del token.
	now := time.Now()
	expirationTime := now.Add(cookieMaxAge)

	// Obtener accesos mínimos a empresas.
	companies := getCompanyAccesses(usuario.ID)

	// Configurar los claims simplificados.
	claims := TokenClaims{
		UsuarioID: uint64(usuario.ID),
		PersonaID: uint64(persona.ID),
		Companies: companies,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sistema-ponceca",
			Subject:   req.Identifier,
		},
	}

	// Generar y firmar el token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())

	if err != nil {
		return nil, "", fiber.NewError(fiber.StatusInternalServerError, "Error generando token")
	}

	// Registrar la sesión con más datos:
	session := Sesion{
		UsuarioID:       uint64(usuario.ID),
		Token:           tokenString,
		FechaExpiracion: expirationTime,
		Activa:          true,
		IP:              ip,
	}
	if err = database.DBconn.Create(&session).Error; err != nil {
		return nil, "", fiber.NewError(fiber.StatusInternalServerError, "Error creando sesión")
	}

	// Construir la respuesta sin incluir el token
	response := fiber.Map{
		"expires_at": expirationTime.Format(time.RFC3339),
		"user": fiber.Map{
			"usuario_id": usuario.ID,
			"persona_id": persona.ID,
		},
		"companies": companies,
	}

	// Devolver la respuesta y el token por separado
	return response, tokenString, nil
}

// ------------------------
// Handler
// Encargado de manejar la petición HTTP.
// ------------------------

func loginHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json") // Forzar respuesta en JSON

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error al parsear body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Solicitud inválida"})
	}

	if req.Identifier == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El identificador y la contraseña son requeridos"})
	}

	// Obtener IP del cliente para seguridad
	ip := c.IP()
	if forwardedIP := c.Get("X-Forwarded-For"); forwardedIP != "" {
		ip = strings.Split(forwardedIP, ",")[0]
	}

	// Llamar al servicio para realizar el login
	response, token, err := loginService(req, ip)
	if err != nil {
		log.Printf("Error en servicio de login para '%s': %v", req.Identifier, err)
		// Si el error es del servicio (ya contiene el status adecuado), se propaga.
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error interno"})
	}

	// Configurar una cookie segura para guardar el token
	expiresAt, _ := time.Parse(time.RFC3339, response["expires_at"].(string))

	// Establecer dominio vacío para usar cookie solo en el mismo sitio
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		Domain:   "", // Siempre mismo sitio
		Expires:  expiresAt,
		MaxAge:   int(cookieMaxAge.Seconds()),
		Secure:   config.IsProductionCookie(), // usar solo secure en producción
		HTTPOnly: true,
		SameSite: "Strict",
	})

	// Add the token to the response for client-side storage options
	// (useful as a fallback if cookies don't work)
	response["token"] = token
	// Agregamos "isAuthenticated": true al response
	response["isAuthenticated"] = true
	return c.JSON(response)
}

func logoutHandler(c *fiber.Ctx) error {
	// Intentar obtener el token tanto del header como de la cookie
	tokenStr := c.Get("Authorization")
	if tokenStr == "" { // Leer la cookie "auth_token"
		tokenStr = c.Cookies(cookieName)
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token no provisto"})
		}
	} else {
		tokenStr = strings.TrimSpace(strings.Replace(tokenStr, "Bearer ", "", 1))
	}

	var sesion Sesion
	if err := database.DBconn.Where("token = ? AND activa = ?", tokenStr, true).First(&sesion).Error; err != nil {
		log.Printf("Error obteniendo sesión para token %s: %v", tokenStr, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Sesión no encontrada o ya inactiva"})
	}

	// Invalidar la sesión
	sesion.Activa = false
	if err := database.DBconn.Save(&sesion).Error; err != nil {
		log.Printf("Error al invalidar sesión para token %s: %v", tokenStr, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al cerrar sesión"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Domain:   "", // Siempre mismo sitio
		Expires:  time.Now().Add(-24 * time.Hour),
		MaxAge:   -1,
		Secure:   config.IsProductionCookie(), // usar solo secure en producción
		HTTPOnly: true,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Sesión cerrada exitosamente"})
}

// Nuevo handler para verificar el estado del token (cookie)
func statusHandler(c *fiber.Ctx) error {
	// Intentar obtener el token desde la cookie
	tokenStr := c.Cookies(cookieName)
	if tokenStr == "" {
		// Siempre devolver status 200 con isAuthenticated: false
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"isAuthenticated": false,
		})
	}

	// Validar el token
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return getJWTKey(), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Token inválido: %v", err)
		// Siempre devolver status 200 con isAuthenticated: false
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"isAuthenticated": false,
		})
	}

	// Extraer claims del token
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		log.Printf("No se pudieron extraer claims del token")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"isAuthenticated": false,
		})
	}

	// Responder con status 200, isAuthenticated: true y la información del usuario y empresas
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isAuthenticated": true,
		"timestamp":       time.Now().Format(time.RFC3339),
		"user": fiber.Map{
			"usuario_id": claims.UsuarioID,
			"persona_id": claims.PersonaID,
		},
		"companies": claims.Companies,
	})
}

// ------------------------
// Registro de rutas
// ------------------------

// RegisterAdvancedLogin registra la ruta de login avanzada.
func RegisterAdvancedLogin(app *fiber.App) {
	api := app.Group("/api")
	status := app.Group("/auth") // Creamos un nuevo grupo de rutas para status

	api.Post("/login", loginHandler)
	api.Post("/logout", logoutHandler)
	status.Get("/status", statusHandler) // Comentamos o eliminamos esta línea

	//api.Get("/status", statusHandler) // Agregamos la nueva ruta independiente
}

func init() {

	registry.RegisterModule(RegisterAdvancedLogin)
}
