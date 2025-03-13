package middleware

import (
	"practicev2/config"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/patrickmn/go-cache"
)

// Variables para configuración
const (
	cookieName      = "auth_token"
	tokenExpiry     = 12 * time.Hour
	cleanupInterval = 10 * time.Minute // Intervalo para limpieza de tokens expirados
	tokenCacheTTL   = 5 * time.Minute  // Tiempo de vida para tokens en caché
	tokenCacheClean = 10 * time.Minute // Intervalo de limpieza para caché
)

// Estructura simple para tokens revocados
type RevokedToken struct {
	ExpiredAt time.Time
}

// Mapa sincronizado para tokens revocados y caché de tokens válidos
var (
	revokedTokens = make(map[string]RevokedToken)
	tokenMutex    = &sync.RWMutex{}
	cleanupDone   = make(chan bool)
	cleanupActive bool
	// Caché para tokens válidos: clave=tokenString, valor=claims
	tokenCache = cache.New(tokenCacheTTL, tokenCacheClean)
)

// StartTokenCleanup inicia la limpieza periódica de tokens expirados
func StartTokenCleanup() {
	// Evitar iniciar múltiples goroutines de limpieza
	tokenMutex.Lock()
	if cleanupActive {
		tokenMutex.Unlock()
		return
	}
	cleanupActive = true
	tokenMutex.Unlock()

	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cleanupExpiredTokens()
			case <-cleanupDone:
				return
			}
		}
	}()
}

// StopTokenCleanup detiene la limpieza de tokens
func StopTokenCleanup() {
	tokenMutex.Lock()
	if cleanupActive {
		cleanupActive = false
		cleanupDone <- true
	}
	tokenMutex.Unlock()
}

// cleanupExpiredTokens elimina tokens expirados
func cleanupExpiredTokens() {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	now := time.Now()
	for token, info := range revokedTokens {
		if now.After(info.ExpiredAt) {
			delete(revokedTokens, token)
		}
	}
}

// RevokeToken añade un token a la lista de revocados
func RevokeToken(tokenString string, expiredAt time.Time) {
	// Eliminar del caché si existe
	tokenCache.Delete(tokenString)

	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	revokedTokens[tokenString] = RevokedToken{ExpiredAt: expiredAt}
}

// isTokenRevoked verifica si un token está revocado
func isTokenRevoked(tokenString string) bool {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()

	_, exists := revokedTokens[tokenString]
	return exists
}

// AuthMiddleware verifica si el usuario está autenticado
func AuthMiddleware() fiber.Handler {
	// Asegurar que la limpieza de tokens esté activa
	StartTokenCleanup()
	// Preparar la clave JWT una sola vez para reutilizarla
	jwtKey := config.GetJWTKey()

	return func(c *fiber.Ctx) error {
		// Intentar obtener el token primero desde la cookie HTTP-only
		tokenString := c.Cookies(cookieName)

		// Si no hay cookie, intentar desde el header de autorización como fallback
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Verificar que exista un token
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Acceso no autorizado: Token no proporcionado",
				"code":  "AUTH_NO_TOKEN",
			})
		}

		// Verificar si el token está en la lista negra (rápido)
		if isTokenRevoked(tokenString) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token revocado o en lista negra",
				"code":  "AUTH_TOKEN_REVOKED",
			})
		}

		// Verificar si el token está en caché
		if cachedClaims, found := tokenCache.Get(tokenString); found {
			// Token encontrado en caché, usar claims almacenados
			c.Locals("user", cachedClaims)
			return c.Next()
		}

		// Validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de firma JWT inválido")
			}
			return jwtKey, nil
		})

		// Verificar errores durante la validación
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido o expirado",
				"code":  "AUTH_INVALID_TOKEN",
			})
		}

		// Extraer y validar claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido: estructura de claims incorrecta",
				"code":  "AUTH_INVALID_CLAIMS",
			})
		}

		// Verificar expiración
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if expTime.Before(time.Now()) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token expirado",
					"code":  "AUTH_TOKEN_EXPIRED",
				})
			}

			// Calcular TTL para caché basado en expiración del token
			remaining := time.Until(expTime)
			cacheTTL := tokenCacheTTL
			if remaining < tokenCacheTTL {
				cacheTTL = remaining
			}

			// Guardar en caché para futuras validaciones
			tokenCache.Set(tokenString, claims, cacheTTL)
		}

		// Verificar y extraer usuario_id con manejo seguro de tipos
		_, ok = claims["usuario_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido: falta o formato incorrecto de usuario_id",
				"code":  "AUTH_INVALID_USER_ID",
			})
		}

		// Almacenar claims en el contexto

		return c.Next()
	}
}

// init inicia la limpieza de tokens al cargar el paquete
func init() {
	StartTokenCleanup()
}
