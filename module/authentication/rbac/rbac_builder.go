package rbac

import (
	"fmt"
	"time"

	"practicev2/module/authentication/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// RBAC BUILDER DINÁMICO - CONFIGURACIÓN FLUIDA Y DECLARATIVA
// =============================================================================

// RBACBuilder proporciona una interfaz fluida para configurar RBAC dinámico
type RBACBuilder struct {
	db         *gorm.DB
	config     *RBACConfig
	modules    map[string]*ModuleDefinition
	quickSetup bool
	isDynamic  bool
	autoSeed   bool
	seeder     *RBACSeeder
}

// NewRBACBuilder crea un nuevo builder para configuración RBAC dinámico
func NewRBACBuilder(db *gorm.DB, config *RBACConfig) *RBACBuilder {
	if config == nil {
		config = DefaultRBACConfig()
	}

	builder := &RBACBuilder{
		db:        db,
		config:    config,
		modules:   make(map[string]*ModuleDefinition),
		isDynamic: true, // Por defecto usar el sistema dinámico
		autoSeed:  true, // Por defecto inicializar permisos automáticamente
	}

	// Crear el seeder
	builder.seeder = NewRBACSeeder(db, config)

	return builder
}

// WithCaching configura el sistema de cache
func (b *RBACBuilder) WithCaching(enabled bool, ttl time.Duration) *RBACBuilder {
	b.config.CacheEnabled = enabled
	b.config.CacheTTL = ttl
	return b
}

// WithAuditLog habilita/deshabilita el log de auditoría
func (b *RBACBuilder) WithAuditLog(enabled bool) *RBACBuilder {
	b.config.EnableAuditLog = enabled
	return b
}

// WithDebugMode habilita/deshabilita el modo debug
func (b *RBACBuilder) WithDebugMode(enabled bool) *RBACBuilder {
	b.config.DebugMode = enabled
	return b
}

// WithAutoSeed configura si inicializar permisos automáticamente
func (b *RBACBuilder) WithAutoSeed(enabled bool) *RBACBuilder {
	b.autoSeed = enabled
	return b
}

// WithRoleHierarchy habilita/deshabilita jerarquía de roles
func (b *RBACBuilder) WithRoleHierarchy(enabled bool) *RBACBuilder {
	// Nota: Esta funcionalidad se puede implementar en futuras versiones
	return b
}

// WithOrganizationMode habilita/deshabilita el modo multi-tenant
func (b *RBACBuilder) WithOrganizationMode(enabled bool) *RBACBuilder {
	b.config.OrganizationMode = enabled
	return b
}

// AddModule agrega un módulo con configuración detallada
func (b *RBACBuilder) AddModule(name string, definition *ModuleDefinition) *RBACBuilder {
	b.modules[name] = definition
	return b
}

// AddQuickModule agrega un módulo con configuración simple
func (b *RBACBuilder) AddQuickModule(name string, resources map[string][]string) *RBACBuilder {
	definition := &ModuleDefinition{
		Name:      name,
		Version:   "1.0.0",
		Resources: make(map[string]*ResourceMetadata),
		Routes:    []*RouteDefinition{},
	}

	// Convertir recursos simples a ResourceMetadata
	for resourceName, actions := range resources {
		definition.Resources[resourceName] = &ResourceMetadata{
			Name:        resourceName,
			Description: fmt.Sprintf("%s management", resourceName),
			Actions:     actions,
			Scopes:      []string{"own", "department", "organization", "all"},
			Module:      name,
		}
	}

	b.modules[name] = definition
	b.quickSetup = true
	return b
}

// Build construye y retorna el motor RBAC dinámico configurado
func (b *RBACBuilder) Build() *DynamicRBACEngine {
	// 1. Inicializar sistema RBAC si está configurado
	if b.autoSeed {
		if err := b.seeder.InitializeRBACSystem(); err != nil {
			if b.config.DebugMode {
				fmt.Printf("[RBAC-Builder] Warning: Failed to initialize RBAC system: %v\n", err)
			}
		}
	}

	// 2. Crear el motor dinámico
	engine := NewDynamicRBACEngine(b.db, b.config)

	// 3. Registrar todos los módulos configurados
	for name, definition := range b.modules {
		if err := engine.RegisterModule(name, definition); err != nil {
			if b.config.DebugMode {
				fmt.Printf("[RBAC-Builder] Error registering module %s: %v\n", name, err)
			}
		}
	}

	if b.config.DebugMode {
		fmt.Printf("[RBAC-Builder] Dynamic engine built with %d modules\n", len(b.modules))
	}

	return engine
}

// =============================================================================
// MÉTODOS DE CONVENIENCIA PARA SETUP RÁPIDO
// =============================================================================

// QuickSetup crea un motor RBAC dinámico con configuración mínima
func QuickSetup(db *gorm.DB, moduleName string, resources map[string][]string) *DynamicRBACEngine {
	return NewRBACBuilder(db, DefaultRBACConfig()).
		AddQuickModule(moduleName, resources).
		Build()
}

// QuickSetupWithOrg crea un motor RBAC dinámico con soporte organizacional
func QuickSetupWithOrg(db *gorm.DB, moduleName string, resources map[string][]string) *DynamicRBACEngine {
	return NewRBACBuilder(db, DefaultRBACConfig()).
		WithOrganizationMode(true).
		AddQuickModule(moduleName, resources).
		Build()
}

// ProductionSetup crea un motor RBAC con configuración de producción
func ProductionSetup(db *gorm.DB, moduleName string, resources map[string][]string) *DynamicRBACEngine {
	config := &RBACConfig{
		CacheEnabled:     true,
		CacheTTL:         15 * time.Minute,
		DefaultDenyAll:   true,
		EnableAuditLog:   true,
		DebugMode:        false,
		OrganizationMode: true,
	}

	return NewRBACBuilder(db, config).
		AddQuickModule(moduleName, resources).
		Build()
}

// DevelopmentSetup crea un motor RBAC con configuración de desarrollo
func DevelopmentSetup(db *gorm.DB, moduleName string, resources map[string][]string) *DynamicRBACEngine {
	config := &RBACConfig{
		CacheEnabled:     false, // Sin cache en desarrollo para facilitar testing
		CacheTTL:         1 * time.Minute,
		DefaultDenyAll:   true,
		EnableAuditLog:   true,
		DebugMode:        true,
		OrganizationMode: true,
	}

	return NewRBACBuilder(db, config).
		AddQuickModule(moduleName, resources).
		Build()
}

// =============================================================================
// PRESETS PARA DIFERENTES TIPOS DE APLICACIÓN
// =============================================================================

// ECommerceSetup preset para aplicaciones de e-commerce
func ECommerceSetup(db *gorm.DB) *DynamicRBACEngine {
	resources := map[string][]string{
		"products":  {"create", "read", "update", "delete", "price", "inventory", "publish"},
		"orders":    {"create", "read", "update", "process", "cancel", "refund", "ship"},
		"customers": {"read", "update", "communicate", "discount", "suspend"},
		"inventory": {"read", "update", "adjust", "audit", "transfer"},
		"reports":   {"read", "create", "export", "share"},
		"settings":  {"read", "update", "configure"},
	}

	return QuickSetupWithOrg(db, "ecommerce", resources)
}

// EducationalSetup preset para instituciones educativas
func EducationalSetup(db *gorm.DB) *DynamicRBACEngine {
	resources := map[string][]string{
		"students":  {"create", "read", "update", "enroll", "grade", "communicate"},
		"teachers":  {"create", "read", "update", "assign", "evaluate", "schedule"},
		"courses":   {"create", "read", "update", "delete", "publish", "enroll"},
		"grades":    {"create", "read", "update", "approve", "publish", "report"},
		"schedules": {"create", "read", "update", "publish", "conflict"},
		"resources": {"create", "read", "update", "reserve", "maintain"},
		"reports":   {"read", "create", "export", "analyze"},
	}

	return QuickSetupWithOrg(db, "educational", resources)
}

// EnterpriseSetup preset para aplicaciones empresariales
func EnterpriseSetup(db *gorm.DB) *DynamicRBACEngine {
	resources := map[string][]string{
		"employees":    {"create", "read", "update", "deactivate", "hire", "evaluate"},
		"departments":  {"create", "read", "update", "delete", "manage", "budget"},
		"projects":     {"create", "read", "update", "delete", "assign", "track"},
		"finances":     {"read", "create", "approve", "report", "audit", "budget"},
		"documents":    {"create", "read", "update", "delete", "share", "approve"},
		"meetings":     {"create", "read", "update", "schedule", "cancel", "record"},
		"performance":  {"read", "evaluate", "report", "goal_setting"},
		"integrations": {"read", "configure", "manage", "monitor"},
	}

	return QuickSetupWithOrg(db, "enterprise", resources)
}

// =============================================================================
// MIDDLEWARE SHORTCUTS - ATAJOS PARA CASOS COMUNES
// =============================================================================

// ProtectResource protege un grupo de rutas Fiber con permisos de recurso
func (e *DynamicRBACEngine) ProtectResource(resource string) fiber.Handler {
	return e.Protect(resource, "read", "organization")
}

// ProtectCreate protege rutas de creación
func (e *DynamicRBACEngine) ProtectCreate(resource string) fiber.Handler {
	return e.Protect(resource, "create", "organization")
}

// ProtectUpdate protege rutas de actualización
func (e *DynamicRBACEngine) ProtectUpdate(resource string) fiber.Handler {
	return e.Protect(resource, "update", "organization")
}

// ProtectDelete protege rutas de eliminación
func (e *DynamicRBACEngine) ProtectDelete(resource string) fiber.Handler {
	return e.Protect(resource, "delete", "organization")
}

// ProtectAdmin protege rutas administrativas
func (e *DynamicRBACEngine) ProtectAdmin(resource string) fiber.Handler {
	return e.Protect(resource, "manage", "organization")
}

// ProtectOwn protege recursos propios del usuario
func (e *DynamicRBACEngine) ProtectOwn(resource string, action string) fiber.Handler {
	return e.Protect(resource, action, "own")
}

// =============================================================================
// CONFIGURACIONES PREDEFINIDAS PARA MÓDULOS ESPECÍFICOS
// =============================================================================

// GastosModuleSetup configura RBAC específicamente para el módulo de gastos
func GastosModuleSetup(db *gorm.DB) *DynamicRBACEngine {
	resources := map[string][]string{
		"expenses":       {"create", "read", "update", "delete", "approve", "reject", "export"},
		"categories":     {"create", "read", "update", "delete", "assign"},
		"budgets":        {"create", "read", "update", "delete", "approve", "track"},
		"reports":        {"read", "create", "export", "share", "analyze"},
		"receipts":       {"upload", "read", "update", "delete", "validate"},
		"approvals":      {"read", "approve", "reject", "delegate", "escalate"},
		"reimbursements": {"create", "read", "process", "pay", "track"},
	}

	return QuickSetupWithOrg(db, "gastos", resources)
}

// AuthModuleSetup configura RBAC para el módulo de autenticación
func AuthModuleSetup(db *gorm.DB) *DynamicRBACEngine {
	resources := map[string][]string{
		"users":         {"create", "read", "update", "delete", "invite", "suspend", "activate"},
		"roles":         {"create", "read", "update", "delete", "assign", "revoke"},
		"permissions":   {"read", "assign", "revoke", "audit"},
		"organizations": {"create", "read", "update", "delete", "manage", "configure"},
		"invitations":   {"create", "read", "resend", "cancel", "accept", "reject"},
		"audit":         {"read", "export", "analyze", "monitor"},
		"sessions":      {"read", "terminate", "monitor"},
	}

	return QuickSetupWithOrg(db, "authentication", resources)
}

// =============================================================================
// MÉTODOS DE GESTIÓN DE PERMISOS
// =============================================================================

// InitializeSystemPermissions inicializa todos los permisos del sistema
func (b *RBACBuilder) InitializeSystemPermissions() error {
	return b.seeder.InitializeRBACSystem()
}

// RefreshPermissions actualiza permisos del sistema
func (b *RBACBuilder) RefreshPermissions() error {
	return b.seeder.RefreshSystemPermissions()
}

// CreateOrganizationWithRoles crea una organización con roles por defecto
func (b *RBACBuilder) CreateOrganizationWithRoles(org interface{}) error {
	// Convertir a models.Organization si es necesario
	if orgModel, ok := org.(*models.Organization); ok {
		return b.seeder.CreateOrganizationWithDefaultRoles(orgModel)
	}
	return fmt.Errorf("invalid organization type")
}

// GetSeeder retorna el seeder para operaciones avanzadas
func (b *RBACBuilder) GetSeeder() *RBACSeeder {
	return b.seeder
}

// =============================================================================
// UTILIDADES DE CONFIGURACIÓN AVANZADA
// =============================================================================

// containsAll verifica si el slice contiene todos los elementos requeridos
func containsAll(slice []string, required []string) bool {
	found := make(map[string]bool)
	for _, item := range slice {
		found[item] = true
	}

	for _, req := range required {
		if !found[req] {
			return false
		}
	}

	return true
}

// =============================================================================
// REGISTRO DE MÓDULOS AUTOMATIZADO
// =============================================================================

// ModuleRegistry mantiene registro de todos los módulos y sus engines
type ModuleRegistry struct {
	modules map[string]*DynamicRBACEngine
}

var globalRegistry = &ModuleRegistry{
	modules: make(map[string]*DynamicRBACEngine),
}

// RegisterModule registra un módulo en el registry global
func RegisterModule(moduleName string, engine *DynamicRBACEngine) {
	globalRegistry.modules[moduleName] = engine
}

// GetModule obtiene el engine RBAC de un módulo
func GetModule(moduleName string) (*DynamicRBACEngine, bool) {
	engine, exists := globalRegistry.modules[moduleName]
	return engine, exists
}

// ListModules lista todos los módulos registrados
func ListModules() []string {
	modules := make([]string, 0, len(globalRegistry.modules))
	for moduleName := range globalRegistry.modules {
		modules = append(modules, moduleName)
	}
	return modules
}
