// Package config centraliza la gestión de variables de entorno y configuración
package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Variables globales para toda la aplicación
var (
	// JWT_KEY es la clave JWT global para todas las funciones
	JWT_KEY string

	// JWT Configuration para autenticación avanzada
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	JWT_ISSUER         string
	JWT_ACCESS_TTL     time.Duration
	JWT_REFRESH_TTL    time.Duration
)

// Init carga las variables de entorno al inicio
func Init() {
	// Cargar JWT_KEY (compatibilidad con código existente)
	JWT_KEY = os.Getenv("JWT_KEY")
	if JWT_KEY == "" {
		JWT_KEY = "default-super-secret-key-for-testing"
		log.Println("⚠️  JWT_KEY no configurado, usando valor por defecto para testing")
	}
	if len(JWT_KEY) < 32 {
		// In a real app, this should probably be a fatal error.
		// For this context, we'll just log a warning.
		log.Println("⚠️  ADVERTENCIA: JWT_KEY debe tener al menos 32 caracteres para seguridad")
	}

	// ===========================================
	// 🔐 CONFIGURACIÓN JWT AVANZADA
	// ===========================================

	// JWT Access Secret
	JWT_ACCESS_SECRET = os.Getenv("JWT_ACCESS_SECRET")
	if JWT_ACCESS_SECRET == "" {
		JWT_ACCESS_SECRET = JWT_KEY // Usar JWT_KEY como fallback
		log.Println("⚠️  JWT_ACCESS_SECRET no configurado, usando JWT_KEY como fallback")
	}

	// JWT Refresh Secret
	JWT_REFRESH_SECRET = os.Getenv("JWT_REFRESH_SECRET")
	if JWT_REFRESH_SECRET == "" {
		JWT_REFRESH_SECRET = JWT_KEY + "_refresh" // Fallback diferenciado
		log.Println("⚠️  JWT_REFRESH_SECRET no configurado, usando fallback")
	}

	// JWT Issuer
	JWT_ISSUER = os.Getenv("JWT_ISSUER")
	if JWT_ISSUER == "" {
		JWT_ISSUER = "practicev2-auth"
		log.Println("⚠️  JWT_ISSUER no configurado, usando valor por defecto")
	}

	// JWT Access TTL (Time To Live)
	accessTTLStr := os.Getenv("JWT_ACCESS_TTL_MINUTES")
	if accessTTLStr != "" {
		if minutes, err := strconv.Atoi(accessTTLStr); err == nil {
			JWT_ACCESS_TTL = time.Duration(minutes) * time.Minute
		} else {
			log.Printf("⚠️  JWT_ACCESS_TTL_MINUTES inválido (%s), usando valor por defecto", accessTTLStr)
			JWT_ACCESS_TTL = 15 * time.Minute
		}
	} else {
		JWT_ACCESS_TTL = 15 * time.Minute // 15 minutos por defecto
	}

	// JWT Refresh TTL
	refreshTTLStr := os.Getenv("JWT_REFRESH_TTL_HOURS")
	if refreshTTLStr != "" {
		if hours, err := strconv.Atoi(refreshTTLStr); err == nil {
			JWT_REFRESH_TTL = time.Duration(hours) * time.Hour
		} else {
			log.Printf("⚠️  JWT_REFRESH_TTL_HOURS inválido (%s), usando valor por defecto", refreshTTLStr)
			JWT_REFRESH_TTL = 168 * time.Hour
		}
	} else {
		JWT_REFRESH_TTL = 168 * time.Hour // 7 días por defecto
	}

	// Configurar la variable de entorno "ENV" en el archivo .env:
	// Ejemplo para desarrollo: ENV=development
	// Ejemplo para producción: ENV=production

	log.Println("✅ Configuración cargada exitosamente")
	log.Printf("📊 JWT Access TTL: %v", JWT_ACCESS_TTL)
	log.Printf("📊 JWT Refresh TTL: %v", JWT_REFRESH_TTL)
	log.Printf("📊 JWT Issuer: %s", JWT_ISSUER)
}

// GetJWTKey devuelve la clave JWT como []byte para su uso en funciones de autenticación
func GetJWTKey() []byte {
	return []byte(JWT_KEY)
}

// ===========================================
// 🔐 FUNCIONES DE CONFIGURACIÓN JWT AVANZADA
// ===========================================

// GetJWTAccessSecret devuelve el secreto para tokens de acceso
func GetJWTAccessSecret() string {
	return JWT_ACCESS_SECRET
}

// GetJWTRefreshSecret devuelve el secreto para tokens de refresh
func GetJWTRefreshSecret() string {
	return JWT_REFRESH_SECRET
}

// GetJWTIssuer devuelve el emisor de tokens JWT
func GetJWTIssuer() string {
	return JWT_ISSUER
}

// GetJWTAccessTTL devuelve el tiempo de vida de los tokens de acceso
func GetJWTAccessTTL() time.Duration {
	return JWT_ACCESS_TTL
}

// GetJWTRefreshTTL devuelve el tiempo de vida de los tokens de refresh
func GetJWTRefreshTTL() time.Duration {
	return JWT_REFRESH_TTL
}

// ===========================================
// 🍪 FUNCIONES DE CONFIGURACIÓN DE COOKIES
// ===========================================

// Aplicación:
// - Para desarrollo: No es necesario configurar nada, por defecto se considera entorno de desarrollo.
// - Para producción: Establecer la variable de entorno COOKIE_CONFIG="true"
//
// Retorna:
//   - true: si estamos en entorno de producción (COOKIE_CONFIG="true")
//   - false: si estamos en entorno de desarrollo (cualquier otro valor o sin configurar)
func IsProductionCookie() bool {
	cookieConfig := os.Getenv("COOKIE_CONFIG")

	return cookieConfig == "true"
}

// GetCookieDomain returns the domain to use for cookies
func GetCookieDomain() string {
	// In development, return empty string for localhost compatibility
	if !IsProductionCookie() {
		return ""
	}

	// In production, return your domain
	// You could also read this from an environment variable
	return "" // Empty string means the cookie applies to the current domain only
}
