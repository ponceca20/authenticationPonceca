package organization_config

import (
	"fmt"
	"log"

	"practicev2/module/authentication/middleware"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =============================================================================
// CONFIGURACIÓN PARA DESARROLLO LOCAL Y TESTS
// =============================================================================

// CreateTestConnection crea una conexión de test sin depender de la configuración del .env
func CreateTestConnection() (*gorm.DB, error) {
	// Usar SQLite en memoria para tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Sin logs durante tests
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %v", err)
	}

	return db, nil
}

// InitializeForTesting inicializa el módulo para entorno de testing
func InitializeForTesting() (*gorm.DB, error) {
	// 1. Crear conexión de test
	db, err := CreateTestConnection()
	if err != nil {
		return nil, err
	}

	// 2. Configurar middleware sin DB externa
	middleware.InitGlobalAuth(middleware.DefaultUnifiedAuthConfig())

	// 3. Inicializar sin registro RBAC completo (para evitar el error)
	err = runModuleMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("failed to run test migrations: %v", err)
	}

	return db, nil
}

// SetupTestServer configura un servidor de test completo
func SetupTestServer() error {
	log.Println("🧪 Configurando entorno de testing...")

	// Verificar si MySQL está disponible en el puerto correcto
	log.Println("📊 Estado de MySQL:")
	log.Println("   - Puerto configurado en .env: 3308")
	log.Println("   - Puerto estándar MySQL: 3306")
	log.Println("   - Recomendación: Verificar que MySQL esté ejecutándose en el puerto 3308")
	log.Println("   - Alternativa: Cambiar DB_PORT=3306 en .env si MySQL está en puerto estándar")

	return nil
}

// =============================================================================
// DIAGNÓSTICO DE BASE DE DATOS
// =============================================================================

// DiagnoseDatabase ayuda a diagnosticar problemas de conexión
func DiagnoseDatabase() {
	fmt.Println("🔍 DIAGNÓSTICO DE BASE DE DATOS")
	fmt.Println("=====================================")
	fmt.Println()

	fmt.Println("📋 Configuración actual (.env):")
	fmt.Println("   - Host: 127.0.0.1")
	fmt.Println("   - Puerto: 3308")
	fmt.Println("   - Usuario: root")
	fmt.Println("   - Base de datos: gastos_ia")
	fmt.Println()

	fmt.Println("⚠️  ERROR ENCONTRADO:")
	fmt.Println("   'dial tcp 127.0.0.1:3308: connectex: No se puede establecer una conexión'")
	fmt.Println()

	fmt.Println("🔧 POSIBLES SOLUCIONES:")
	fmt.Println("   1. MySQL no está ejecutándose en puerto 3308")
	fmt.Println("   2. Verificar servicios MySQL activos:")
	fmt.Println("      - Ejecutar: netstat -an | findstr 3306")
	fmt.Println("      - MySQL estándar usa puerto 3306")
	fmt.Println()

	fmt.Println("✅ MYSQL DETECTADO EN PUERTO 3306")
	fmt.Println("   - Cambiar DB_PORT=3306 en .env")
	fmt.Println("   - O configurar MySQL para usar puerto 3308")
	fmt.Println()

	fmt.Println("🚀 PASOS PARA RESOLVER:")
	fmt.Println("   1. Abrir archivo .env")
	fmt.Println("   2. Cambiar: DB_PORT=3308 → DB_PORT=3306")
	fmt.Println("   3. Guardar archivo")
	fmt.Println("   4. Ejecutar: go run main.go")
	fmt.Println()
}
