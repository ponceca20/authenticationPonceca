// Package config centraliza la gestión de variables de entorno y configuración
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// ===========================================
// 📋 ESTRUCTURAS DE CONFIGURACIÓN
// ===========================================

// Config contiene toda la configuración de la aplicación
type Config struct {
	Environment string
	App         AppConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Cookie      CookieConfig
	Security    SecurityConfig
	CORS        CORSConfig
	Email       EmailConfig
	Logging     LoggingConfig
}

// AppConfig configuración de la aplicación
type AppConfig struct {
	Host string
	Port string
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// RedisConfig configuración de Redis
type RedisConfig struct {
	URL string
}

// JWTConfig configuración JWT completa
type JWTConfig struct {
	Key           string
	AccessSecret  string
	RefreshSecret string
	Issuer        string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// CookieConfig configuración de cookies
type CookieConfig struct {
	Production bool
	Domain     string
	Secure     bool
}

// SecurityConfig configuración de seguridad
type SecurityConfig struct {
	RateLimitRequests      int
	RateLimitWindowMinutes int
}

// CORSConfig configuración CORS
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// EmailConfig configuración de email
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	Username    string
	Password    string
	FromAddress string
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level  string
	Format string
}

// ===========================================
// 🌍 VARIABLES GLOBALES
// ===========================================

var (
	GlobalConfig *Config
	// Variables de compatibilidad hacia atrás
	JWT_KEY            string
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	JWT_ISSUER         string
	JWT_ACCESS_TTL     time.Duration
	JWT_REFRESH_TTL    time.Duration
)

// ===========================================
// 🚀 FUNCIÓN DE INICIALIZACIÓN
// ===========================================

// loadEnvFile carga el archivo .env si existe
func loadEnvFile() {
	// Buscar el archivo .env en el directorio raíz del proyecto
	envPaths := []string{
		".env",                // directorio actual
		"../.env",             // un nivel arriba
		"../../.env",          // dos niveles arriba
		"../../../.env",       // tres niveles arriba
		"../../../../.env",    // cuatro niveles arriba
		"../../../../../.env", // cinco niveles arriba
	}

	for _, envPath := range envPaths {
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Printf("⚠️  Error cargando %s: %v", envPath, err)
			} else {
				log.Printf("✅ Archivo .env cargado desde: %s", envPath)
				return
			}
		}
	}

	// Si no se encuentra ningún archivo .env, continuar sin él
	log.Printf("ℹ️  No se encontró archivo .env, usando variables de entorno del sistema")
}

// Init carga todas las variables de entorno al inicio
func Init() {
	// Cargar archivo .env si existe
	loadEnvFile()

	config := &Config{}

	// Cargar configuración del entorno
	config.Environment = getEnvWithDefault("ENV", "development")

	// Configuración de la aplicación
	config.App = AppConfig{
		Host: getEnvWithDefault("APP_HOST", "127.0.0.1"),
		Port: getEnvWithDefault("APP_PORT", "3030"),
	}

	// Configuración de base de datos
	config.Database = DatabaseConfig{
		User:     getEnvWithDefault("DB_USER", "root"),
		Password: getEnvWithDefault("DB_PASSWORD", ""),
		Host:     getEnvWithDefault("DB_HOST", "127.0.0.1"),
		Port:     getEnvWithDefault("DB_PORT", "3306"),
		Name:     getEnvWithDefault("DB_NAME", "app_db"),
	}

	// Configuración de Redis
	config.Redis = RedisConfig{
		URL: getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0"),
	}

	// Configuración JWT
	config.JWT = loadJWTConfig()

	// Configuración de Cookies
	config.Cookie = loadCookieConfig()

	// Configuración de Seguridad
	config.Security = loadSecurityConfig()

	// Configuración CORS
	config.CORS = loadCORSConfig()

	// Configuración de Email
	config.Email = loadEmailConfig()

	// Configuración de Logging
	config.Logging = loadLoggingConfig()

	// Asignar a la variable global
	GlobalConfig = config

	// Mantener compatibilidad hacia atrás
	setBackwardCompatibilityVars(config)

	// Validar configuración
	validateConfig(config)

	// Log de confirmación
	logConfigurationSummary(config)
}

// ===========================================
// 🔧 FUNCIONES DE CARGA DE CONFIGURACIÓN
// ===========================================

// loadJWTConfig carga la configuración JWT
func loadJWTConfig() JWTConfig {
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("❌ ERROR CRÍTICO: JWT_KEY es requerido")
	}

	accessSecret := getEnvWithDefault("JWT_ACCESS_SECRET", jwtKey)
	refreshSecret := getEnvWithDefault("JWT_REFRESH_SECRET", jwtKey+"_refresh")
	issuer := getEnvWithDefault("JWT_ISSUER", "ponceca-auth-system")

	// Parse TTL values
	accessTTL := parseDurationMinutes("JWT_ACCESS_TTL_MINUTES", 15)
	refreshTTL := parseDurationHours("JWT_REFRESH_TTL_HOURS", 168)

	return JWTConfig{
		Key:           jwtKey,
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		Issuer:        issuer,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
	}
}

// loadCookieConfig carga la configuración de cookies
func loadCookieConfig() CookieConfig {
	production := getEnvWithDefault("COOKIE_CONFIG", "false") == "true"
	domain := getEnvWithDefault("COOKIE_DOMAIN", "")
	secure := getEnvWithDefault("COOKIE_SECURE", "false") == "true"

	return CookieConfig{
		Production: production,
		Domain:     domain,
		Secure:     secure,
	}
}

// loadSecurityConfig carga la configuración de seguridad
func loadSecurityConfig() SecurityConfig {
	rateLimitRequests := parseIntWithDefault("RATE_LIMIT_REQUESTS", 100)
	rateLimitWindow := parseIntWithDefault("RATE_LIMIT_WINDOW_MINUTES", 15)

	return SecurityConfig{
		RateLimitRequests:      rateLimitRequests,
		RateLimitWindowMinutes: rateLimitWindow,
	}
}

// loadCORSConfig carga la configuración CORS
func loadCORSConfig() CORSConfig {
	origins := parseStringSlice("CORS_ALLOWED_ORIGINS", []string{"*"})
	methods := parseStringSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := parseStringSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Requested-With"})

	return CORSConfig{
		AllowedOrigins: origins,
		AllowedMethods: methods,
		AllowedHeaders: headers,
	}
}

// loadEmailConfig carga la configuración de email
func loadEmailConfig() EmailConfig {
	smtpPort := parseIntWithDefault("SMTP_PORT", 587)

	return EmailConfig{
		SMTPHost:    getEnvWithDefault("SMTP_HOST", ""),
		SMTPPort:    smtpPort,
		Username:    getEnvWithDefault("SMTP_USERNAME", ""),
		Password:    getEnvWithDefault("SMTP_PASSWORD", ""),
		FromAddress: getEnvWithDefault("EMAIL_FROM", ""),
	}
}

// loadLoggingConfig carga la configuración de logging
func loadLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:  getEnvWithDefault("LOG_LEVEL", "info"),
		Format: getEnvWithDefault("LOG_FORMAT", "text"),
	}
}

// ===========================================
// 🔧 FUNCIONES AUXILIARES
// ===========================================

// getEnvWithDefault obtiene una variable de entorno con valor por defecto
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseIntWithDefault parsea un entero con valor por defecto
func parseIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		log.Printf("⚠️  %s inválido (%s), usando valor por defecto: %d", key, value, defaultValue)
	}
	return defaultValue
}

// parseDurationMinutes parsea duración en minutos
func parseDurationMinutes(key string, defaultMinutes int) time.Duration {
	minutes := parseIntWithDefault(key, defaultMinutes)
	return time.Duration(minutes) * time.Minute
}

// parseDurationHours parsea duración en horas
func parseDurationHours(key string, defaultHours int) time.Duration {
	hours := parseIntWithDefault(key, defaultHours)
	return time.Duration(hours) * time.Hour
}

// parseStringSlice parsea un string separado por comas en slice
func parseStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(strings.ReplaceAll(value, " ", ""), ",")
	}
	return defaultValue
}

// setBackwardCompatibilityVars mantiene compatibilidad hacia atrás
func setBackwardCompatibilityVars(config *Config) {
	JWT_KEY = config.JWT.Key
	JWT_ACCESS_SECRET = config.JWT.AccessSecret
	JWT_REFRESH_SECRET = config.JWT.RefreshSecret
	JWT_ISSUER = config.JWT.Issuer
	JWT_ACCESS_TTL = config.JWT.AccessTTL
	JWT_REFRESH_TTL = config.JWT.RefreshTTL
}

// validateConfig valida la configuración cargada
func validateConfig(config *Config) {
	errors := []string{}

	// Validar JWT Key
	if len(config.JWT.Key) < 32 {
		errors = append(errors, "JWT_KEY debe tener al menos 32 caracteres")
	}

	// Validar configuración de base de datos
	if config.Database.User == "" {
		errors = append(errors, "DB_USER es requerido")
	}
	if config.Database.Name == "" {
		errors = append(errors, "DB_NAME es requerido")
	}

	// Si hay errores críticos, salir
	if len(errors) > 0 {
		for _, err := range errors {
			log.Printf("❌ ERROR DE CONFIGURACIÓN: %s", err)
		}
		log.Fatal("❌ Configuración inválida, cerrando aplicación")
	}
}

// logConfigurationSummary muestra un resumen de la configuración
func logConfigurationSummary(config *Config) {
	log.Println("✅ Configuración cargada exitosamente")
	log.Printf("🌍 Entorno: %s", config.Environment)
	log.Printf("🚀 Aplicación: %s:%s", config.App.Host, config.App.Port)
	log.Printf("🗄️  Base de datos: %s@%s:%s/%s", config.Database.User, config.Database.Host, config.Database.Port, config.Database.Name)
	log.Printf("🔐 JWT Issuer: %s", config.JWT.Issuer)
	log.Printf("📊 JWT Access TTL: %v", config.JWT.AccessTTL)
	log.Printf("📊 JWT Refresh TTL: %v", config.JWT.RefreshTTL)
	log.Printf("🍪 Cookies Production: %v", config.Cookie.Production)
	log.Printf("📝 Log Level: %s", config.Logging.Level)
}

// ===========================================
// 🔧 FUNCIONES DE ACCESO (COMPATIBILIDAD HACIA ATRÁS)
// ===========================================

// GetJWTKey devuelve la clave JWT como []byte para su uso en funciones de autenticación
func GetJWTKey() []byte {
	return []byte(JWT_KEY)
}

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

// IsProductionCookie retorna true si estamos en entorno de producción
func IsProductionCookie() bool {
	if GlobalConfig != nil {
		return GlobalConfig.Cookie.Production
	}
	// Fallback para compatibilidad
	return getEnvWithDefault("COOKIE_CONFIG", "false") == "true"
}

// GetCookieDomain devuelve el dominio para las cookies
func GetCookieDomain() string {
	if GlobalConfig != nil {
		return GlobalConfig.Cookie.Domain
	}
	// Fallback para compatibilidad
	return getEnvWithDefault("COOKIE_DOMAIN", "")
}

// GetCookieSecure devuelve si las cookies deben ser secure
func GetCookieSecure() bool {
	if GlobalConfig != nil {
		return GlobalConfig.Cookie.Secure
	}
	// Fallback para compatibilidad
	return getEnvWithDefault("COOKIE_SECURE", "false") == "true"
}

// ===========================================
// 🔧 FUNCIONES DE ACCESO A CONFIGURACIÓN NUEVA
// ===========================================

// GetConfig devuelve la configuración completa
func GetConfig() *Config {
	return GlobalConfig
}

// GetEnvironment devuelve el entorno actual
func GetEnvironment() string {
	if GlobalConfig != nil {
		return GlobalConfig.Environment
	}
	return getEnvWithDefault("ENV", "development")
}

// IsProduction retorna true si estamos en producción
func IsProduction() bool {
	return GetEnvironment() == "production"
}

// IsDevelopment retorna true si estamos en desarrollo
func IsDevelopment() bool {
	return GetEnvironment() == "development"
}

// GetDatabaseURL construye la URL de conexión a la base de datos
func GetDatabaseURL() string {
	if GlobalConfig != nil {
		cfg := GlobalConfig.Database
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	}
	// Fallback para compatibilidad
	user := getEnvWithDefault("DB_USER", "root")
	password := getEnvWithDefault("DB_PASSWORD", "")
	host := getEnvWithDefault("DB_HOST", "127.0.0.1")
	port := getEnvWithDefault("DB_PORT", "3306")
	name := getEnvWithDefault("DB_NAME", "app_db")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
}

// GetRedisURL devuelve la URL de Redis
func GetRedisURL() string {
	if GlobalConfig != nil {
		return GlobalConfig.Redis.URL
	}
	return getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0")
}

// GetAppAddress devuelve la dirección completa de la aplicación
func GetAppAddress() string {
	if GlobalConfig != nil {
		return fmt.Sprintf("%s:%s", GlobalConfig.App.Host, GlobalConfig.App.Port)
	}
	host := getEnvWithDefault("APP_HOST", "127.0.0.1")
	port := getEnvWithDefault("APP_PORT", "3030")
	return fmt.Sprintf("%s:%s", host, port)
}
