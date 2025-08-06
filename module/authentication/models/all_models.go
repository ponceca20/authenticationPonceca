package models

import "time"

// =============================================================================
// ARCHIVO CONSOLIDADO DE TODOS LOS MODELOS DE AUTENTICACIÓN
// =============================================================================
// Este archivo contiene todos los modelos del sistema de autenticación
// en un solo lugar para facilitar el acceso y análisis por IA.
//
// Modelos incluidos:
// - Identity (User)
// - Organization
// - OrganizationalMembership
// - Role y Permission
// - CustomerProfile
// - CustomerPreferences
// - UserProfile
// - ShippingAddress
// - Department
// - AuditLog
// - RefreshToken
// - PasswordResetToken
// - Invitation
// - GuestSession
// =============================================================================

// =============================================================================
// MODELO PRINCIPAL: IDENTITY (USER)
// =============================================================================

// Identity represents the central, unified identity of a person in the system.
// It holds the core credentials and personal information.
type Identity struct {
	ID                  string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Email               string     `json:"email" gorm:"uniqueIndex;size:191;not null"`
	PasswordHash        string     `json:"-" gorm:"column:password_hash;size:255;not null"`
	EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty"`
	FirstName           string     `json:"first_name" gorm:"size:100;not null"`
	LastName            string     `json:"last_name" gorm:"size:100;not null"`
	Avatar              string     `json:"avatar,omitempty" gorm:"size:500"`
	Phone               string     `json:"phone,omitempty" gorm:"size:20"`
	DateOfBirth         *time.Time `json:"date_of_birth,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempts int        `json:"-" gorm:"default:0;type:smallint"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// User represents the primary user entity in the system.
// This is an alias to Identity to maintain backward compatibility and provide
// a more intuitive interface for general user operations.
type User = Identity

// TableName ensures User uses the same table as Identity
func (User) TableName() string {
	return "identity"
}

// =============================================================================
// ORGANIZACIONES Y MEMBRESÍAS
// =============================================================================

// Organization represents a multi-tenant entity like a company or a school.
type Organization struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string     `json:"name" gorm:"not null;size:255"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;size:191;not null"`
	Type        string     `json:"type" gorm:"not null;size:100"` // e.g., "company", "educational_institution"
	Description string     `json:"description" gorm:"size:1000"`
	Avatar      string     `json:"avatar,omitempty" gorm:"size:500"`
	Website     string     `json:"website,omitempty" gorm:"size:255"`
	Phone       string     `json:"phone,omitempty" gorm:"size:20"`
	Address     string     `json:"address,omitempty" gorm:"size:500"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for Organization model
func (Organization) TableName() string {
	return "organization"
}

// OrganizationalMembership links an Identity to an Organization with a specific Role and context.
type OrganizationalMembership struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID     string       `json:"identity_id" gorm:"index;type:varchar(36);not null"`
	Identity       Identity     `json:"identity" gorm:"foreignKey:IdentityID"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36);not null"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	RoleID         string       `json:"role_id" gorm:"index;type:varchar(36);not null"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`

	// Context-specific data
	Department string `json:"department,omitempty" gorm:"size:255"`        // For employees
	Grade      string `json:"grade,omitempty" gorm:"size:50"`              // For students
	Subject    string `json:"subject,omitempty" gorm:"size:255"`           // For teachers
	StudentID  string `json:"student_id,omitempty" gorm:"index;size:100"`  // Unique student ID
	EmployeeID string `json:"employee_id,omitempty" gorm:"index;size:100"` // Unique employee ID

	// Temporal control
	ActiveFrom  time.Time  `json:"active_from" gorm:"not null"`
	ActiveUntil *time.Time `json:"active_until,omitempty"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// Standard model fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for OrganizationalMembership model
func (OrganizationalMembership) TableName() string {
	return "organizational_membership"
}

// Department represents a subdivision within an organization.
type Department struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36)"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	ParentID       *string      `json:"parent_id,omitempty" gorm:"index;type:varchar(36)"` // For hierarchical structures
	Name           string       `json:"name" gorm:"not null;size:255"`
	Description    string       `json:"description" gorm:"size:500"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}

// =============================================================================
// ROLES Y PERMISOS
// =============================================================================

// Role defines a set of permissions within an organization.
type Role struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"type:varchar(36);not null;uniqueIndex:idx_org_role_name,priority:1"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	Name           string       `json:"name" gorm:"not null;size:100;uniqueIndex:idx_org_role_name,priority:2"`
	DisplayName    string       `json:"display_name" gorm:"size:255"`
	Description    string       `json:"description" gorm:"size:500"`
	HierarchyLevel int          `json:"hierarchy_level" gorm:"default:0;type:smallint"`
	IsSystemRole   bool         `json:"is_system_role" gorm:"default:false"`
	Permissions    []Permission `json:"permissions" gorm:"many2many:role_permission;"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "role"
}

// Permission defines a specific action that can be performed on a resource.
type Permission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Resource  string    `json:"resource" gorm:"not null;size:100;uniqueIndex:idx_resource_action_scope"` // e.g., "users", "products"
	Action    string    `json:"action" gorm:"not null;size:100;uniqueIndex:idx_resource_action_scope"`   // e.g., "create", "read", "update", "delete"
	Scope     string    `json:"scope" gorm:"not null;size:50;uniqueIndex:idx_resource_action_scope"`     // e.g., "own", "department", "organization", "all"
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permission"
}

//=============================================================================
// SISTEMA DE CONFIGURACIÓN DINÁMICO DE RECURSOS POR ORGANIZACIÓN
// =============================================================================

// OrganizationModuleConfig permite a cada organización personalizar módulos
type OrganizationModuleConfig struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36);not null"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	ModuleName     string       `json:"module_name" gorm:"not null;size:100;index"`
	IsEnabled      bool         `json:"is_enabled" gorm:"default:true;index"`

	// Configuración JSON del módulo
	Resources    string `json:"resources" gorm:"type:json"`     // {"expenses": ["create", "read", "approve"]}
	DefaultRoles string `json:"default_roles" gorm:"type:json"` // Roles predefinidos del módulo
	Settings     string `json:"settings" gorm:"type:json"`      // Configuración específica

	// Metadata
	Version     string   `json:"version" gorm:"size:20;default:'1.0.0'"`
	InstalledBy string   `json:"installed_by" gorm:"type:varchar(36)"`
	Installer   Identity `json:"installer,omitempty" gorm:"foreignKey:InstalledBy"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName especifica el nombre de tabla
func (OrganizationModuleConfig) TableName() string {
	return "organization_module_config"
}

//=============================================================================
// RECURSOS BASE DEL SISTEMA - MANTIENE COMPATIBILIDAD
// =============================================================================

// SystemResources define recursos base del sistema (MANTENER PARA COMPATIBILIDAD)
var SystemResources = map[string][]string{
	"users":        {"create", "read", "update", "delete", "invite", "suspend"},
	"students":     {"read", "update", "grade", "report", "communicate"},
	"teachers":     {"read", "update", "assign", "evaluate"},
	"courses":      {"create", "read", "update", "delete", "enroll"},
	"grades":       {"create", "read", "update", "approve", "publish"},
	"finances":     {"read", "create", "update", "approve", "report"},
	"departments":  {"create", "read", "update", "delete", "manage"},
	"products":     {"create", "read", "update", "delete", "price", "inventory"},
	"orders":       {"create", "read", "update", "process", "refund"},
	"customers":    {"read", "update", "communicate", "discount"},
	"audit":        {"read", "export"},
	"settings":     {"read", "update"},
	"integrations": {"read", "configure"},
}

// ModuleResources organiza recursos por módulos para configuración dinámica
var ModuleResources = map[string]map[string][]string{
	"authentication": {
		"users":       {"create", "read", "update", "delete", "invite", "suspend"},
		"students":    {"read", "update", "grade", "report", "communicate"},
		"teachers":    {"read", "update", "assign", "evaluate"},
		"departments": {"create", "read", "update", "delete", "manage"},
		"audit":       {"read", "export"},
		"settings":    {"read", "update"},
	},
	"ecommerce": {
		"products":  {"create", "read", "update", "delete", "price", "inventory"},
		"orders":    {"create", "read", "update", "process", "refund"},
		"customers": {"read", "update", "communicate", "discount"},
		"inventory": {"read", "update", "count", "transfer", "adjust"},
	},
	"education": {
		"students": {"read", "update", "grade", "report", "communicate"},
		"teachers": {"read", "update", "assign", "evaluate"},
		"courses":  {"create", "read", "update", "delete", "enroll"},
		"grades":   {"create", "read", "update", "approve", "publish"},
	},
	"expenses": {
		"expenses":   {"create", "read", "update", "delete", "approve", "reject"},
		"budgets":    {"create", "read", "update", "approve", "monitor", "report"},
		"categories": {"create", "read", "update", "delete", "assign"},
		"approvals":  {"view", "approve", "reject", "delegate"},
	},
	"inventory": {
		"inventory":  {"create", "read", "update", "delete", "transfer", "adjust", "audit"},
		"warehouses": {"create", "read", "update", "delete", "manage", "assign"},
		"stock":      {"read", "update", "reserve", "release", "count", "audit"},
		"transfers":  {"create", "read", "approve", "cancel", "receive", "dispatch"},
		"purchases":  {"create", "read", "update", "approve", "receive", "cancel"},
		"suppliers":  {"create", "read", "update", "delete", "evaluate", "block"},
	},
	"reports": {
		"reports":    {"inventory", "financial", "sales", "expenses", "custom", "export"},
		"analytics":  {"view", "create", "export", "schedule"},
		"dashboards": {"view", "create", "update", "share"},
	},
}

// AuthorizationScopes defines the different levels of access.
var AuthorizationScopes = []string{
	"own",
	"department",
	"organization",
	"all",
}

// =============================================================================
// PERFILES DE CLIENTE Y E-COMMERCE
// =============================================================================

// CustomerProfile holds e-commerce specific data for an Identity.
type CustomerProfile struct {
	ID         string   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string   `json:"identity_id" gorm:"uniqueIndex;type:varchar(36);not null"`
	Identity   Identity `json:"identity" gorm:"foreignKey:IdentityID"`

	// E-commerce specific data
	CustomerNumber         string  `json:"customer_number" gorm:"uniqueIndex;size:191;not null"`
	PreferredPaymentMethod string  `json:"preferred_payment_method,omitempty" gorm:"size:100"`
	CreditLimit            float64 `json:"credit_limit" gorm:"default:0;type:decimal(15,2)"`
	TotalSpent             float64 `json:"total_spent" gorm:"default:0;type:decimal(15,2)"`
	LoyaltyPoints          int     `json:"loyalty_points" gorm:"default:0;type:int"`

	// Marketing preferences
	AcceptsMarketing  bool   `json:"accepts_marketing" gorm:"default:false"`
	PreferredLanguage string `json:"preferred_language" gorm:"default:'es';size:10"`

	// Standard model fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// CustomerPreferences stores various preferences for a customer.
type CustomerPreferences struct {
	ID                 string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CustomerProfileID  string          `json:"customer_profile_id" gorm:"uniqueIndex;type:varchar(36);not null"`
	CustomerProfile    CustomerProfile `json:"customer_profile" gorm:"foreignKey:CustomerProfileID"`
	Theme              string          `json:"theme,omitempty" gorm:"default:'light';size:20"`
	Language           string          `json:"language,omitempty" gorm:"default:'es';size:5"`
	TimeZone           string          `json:"time_zone,omitempty" gorm:"default:'UTC';size:50"`
	EmailNotifications bool            `json:"email_notifications" gorm:"default:true"`
	SmsNotifications   bool            `json:"sms_notifications" gorm:"default:false"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// TableName overrides the table name used by CustomerPreferences to ensure singular naming
func (CustomerPreferences) TableName() string {
	return "customer_preference"
}

// ShippingAddress represents a shipping address for a customer.
type ShippingAddress struct {
	ID                string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CustomerProfileID string          `json:"customer_profile_id" gorm:"index;type:varchar(36)"`
	CustomerProfile   CustomerProfile `json:"customer_profile" gorm:"foreignKey:CustomerProfileID"`
	AddressLine1      string          `json:"address_line_1" gorm:"not null;size:255"`
	AddressLine2      string          `json:"address_line_2,omitempty" gorm:"size:255"`
	City              string          `json:"city" gorm:"not null;size:100"`
	State             string          `json:"state" gorm:"not null;size:100"`
	PostalCode        string          `json:"postal_code" gorm:"not null;size:20"`
	Country           string          `json:"country" gorm:"not null;size:100"`
	IsDefault         bool            `json:"is_default" gorm:"default:false"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
}

// =============================================================================
// PERFIL DE USUARIO EXTENDIDO
// =============================================================================

// UserProfile contains extended, non-essential information for an identity.
type UserProfile struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"uniqueIndex;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Bio        string    `json:"bio,omitempty" gorm:"size:1000"`
	Location   string    `json:"location,omitempty" gorm:"size:255"`
	Website    string    `json:"website,omitempty" gorm:"size:255"`
	Socials    string    `json:"socials,omitempty" gorm:"type:json"` // JSON blob for social links
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// =============================================================================
// AUDITORÍA Y SEGURIDAD
// =============================================================================

// AuditLog records an action performed by a user in the system.
type AuditLog struct {
	ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID *string   `json:"organization_id,omitempty" gorm:"index;type:varchar(36)"` // Optional, for system-level actions
	IdentityID     string    `json:"identity_id" gorm:"index;type:varchar(36);not null"`
	Identity       Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Action         string    `json:"action" gorm:"not null;size:255;index"` // e.g., "user.login", "organization.create"
	Resource       string    `json:"resource" gorm:"size:100;index"`        // e.g., "organization"
	ResourceID     string    `json:"resource_id" gorm:"index;size:100"`     // e.g., the ID of the created organization
	Status         string    `json:"status" gorm:"size:50"`                 // e.g., "success", "failure"
	IPAddress      string    `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent      string    `json:"user_agent,omitempty" gorm:"size:1000"`
	Details        string    `json:"details,omitempty" gorm:"type:json"`
	Timestamp      time.Time `json:"timestamp" gorm:"not null;index"`
}

// =============================================================================
// TOKENS DE AUTENTICACIÓN
// =============================================================================

// RefreshToken stores a refresh token for an identity.
type RefreshToken struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Token      string    `json:"-" gorm:"uniqueIndex;size:512"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsRevoked  bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}

// PasswordResetToken stores a token for the password reset process.
type PasswordResetToken struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Token      string    `json:"-" gorm:"uniqueIndex;size:512"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// EmailVerificationToken stores a token for the email verification process.
type EmailVerificationToken struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	IdentityID string    `json:"identity_id" gorm:"index;type:varchar(36)"`
	Identity   Identity  `json:"identity" gorm:"foreignKey:IdentityID"`
	Token      string    `json:"-" gorm:"uniqueIndex;size:512"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// =============================================================================
// INVITACIONES Y SESIONES DE INVITADOS
// =============================================================================

// Invitation represents an invitation for a user to join an organization.
type Invitation struct {
	ID             string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OrganizationID string       `json:"organization_id" gorm:"index;type:varchar(36);not null"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	InviterID      string       `json:"inviter_id" gorm:"index;type:varchar(36);not null"` // The user who sent the invitation
	Inviter        Identity     `json:"inviter" gorm:"foreignKey:InviterID"`
	Email          string       `json:"email" gorm:"not null;size:191"`
	RoleID         string       `json:"role_id" gorm:"index;type:varchar(36);not null"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`
	Token          string       `json:"-" gorm:"uniqueIndex;size:512;not null"`
	Status         string       `json:"status" gorm:"default:'pending';size:50"` // e.g., pending, accepted, expired
	ExpiresAt      time.Time    `json:"expires_at" gorm:"not null"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// GuestSession represents a temporary session for a user who has not registered.
type GuestSession struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(36)"`
	SessionToken string `json:"session_token" gorm:"uniqueIndex;size:512;not null"`

	// Temporary guest data
	Email     string `json:"email,omitempty" gorm:"size:191"`
	FirstName string `json:"first_name,omitempty" gorm:"size:100"`
	LastName  string `json:"last_name,omitempty" gorm:"size:100"`
	Phone     string `json:"phone,omitempty" gorm:"size:20"`

	// Activity tracking
	CartData     *string   `json:"cart_data,omitempty" gorm:"type:json"` // JSON blob for cart
	LastActivity time.Time `json:"last_activity" gorm:"not null"`
	IPAddress    string    `json:"ip_address,omitempty" gorm:"size:45"` // IPv6 support
	UserAgent    string    `json:"user_agent,omitempty" gorm:"size:1000"`

	// Auto-expiration
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// =============================================================================
// RESUMEN DE MODELOS PRINCIPALES
// =============================================================================
/*
MODELOS CORE:
1. Identity (User) - Identidad central del usuario
2. Organization - Organizaciones (empresas, escuelas, etc.)
3. OrganizationalMembership - Relación usuario-organización-rol
4. Role - Roles dentro de organizaciones
5. Permission - Permisos específicos

MODELOS DE PERFIL:
6. CustomerProfile - Perfil de cliente e-commerce
7. CustomerPreferences - Preferencias del cliente
8. UserProfile - Perfil extendido del usuario
9. ShippingAddress - Direcciones de envío

MODELOS ORGANIZACIONALES:
10. Department - Departamentos dentro de organizaciones

MODELOS DE SEGURIDAD:
11. AuditLog - Registro de auditoría
12. RefreshToken - Tokens de actualización
13. PasswordResetToken - Tokens de reseteo de contraseña

MODELOS DE INVITACIÓN:
14. Invitation - Invitaciones a organizaciones
15. GuestSession - Sesiones de usuarios invitados

RELACIONES PRINCIPALES:
- Identity (1) <-> (0..1) CustomerProfile
- Identity (1) <-> (0..1) UserProfile
- Identity (1) <-> (0..*) OrganizationalMembership
- Organization (1) <-> (0..*) OrganizationalMembership
- Organization (1) <-> (0..*) Role
- Organization (1) <-> (0..*) Department
- CustomerProfile (1) <-> (0..1) CustomerPreferences
- CustomerProfile (1) <-> (0..*) ShippingAddress
- Role (1) <-> (0..*) Permission (many-to-many)
*/
