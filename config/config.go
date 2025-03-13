// Package config centraliza la gestión de variables de entorno y configuración
package config

import (
	"log"
	"os"
)

// Variables globales para toda la aplicación
var (
	// JWT_KEY es la clave JWT global para todas las funciones
	JWT_KEY string
)

// Init carga las variables de entorno al inicio
func Init() {
	// Cargar JWT_KEY
	JWT_KEY = os.Getenv("JWT_KEY")
	if JWT_KEY == "" {
		log.Fatal("JWT_KEY no configurado en variables de entorno")
	}
	if len(JWT_KEY) < 15 {
		log.Fatal("JWT_KEY debe tener al menos 32 caracteres para seguridad")
	}

	// Configurar la variable de entorno "ENV" en el archivo .env:
	// Ejemplo para desarrollo: ENV=development
	// Ejemplo para producción: ENV=production

	log.Println("Configuración cargada")
}

// GetJWTKey devuelve la clave JWT como []byte para su uso en funciones de autenticación
func GetJWTKey() []byte {
	return []byte(JWT_KEY)
}

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
