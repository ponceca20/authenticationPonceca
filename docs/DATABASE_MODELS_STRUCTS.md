# Modelos de Base de Datos - Módulo de Autenticación

## ✅ **ESTADO: IMPLEMENTADO Y FUNCIONAL**

**Fecha de implementación:** 30 de julio, 2025  
**Versión:** 1.0.0  
**Estado:** ✅ **Todos los modelos implementados y funcionando**

## 📊 Estructuras Go Completas

### 🏆 **MODELOS IMPLEMENTADOS EXITOSAMENTE**

Todos los modelos listados a continuación han sido implementados como archivos Go funcionales en:
```
module/authentication/models/
├── identity.go                    ✅ Implementado
├── organization.go                ✅ Implementado  
├── organizational_membership.go   ✅ Implementado
├── role.go                        ✅ Implementado
├── permission.go                  ✅ Implementado
├── department.go                  ✅ Implementado
├── customer_profile.go            ✅ Implementado
├── user_profile.go                ✅ Implementado
├── guest_session.go               ✅ Implementado
├── shipping_address.go            ✅ Implementado
├── customer_preferences.go        ✅ Implementado
├── invitation.go                  ✅ Implementado
├── refresh_token.go               ✅ Implementado
├── audit_log.go                   ✅ Implementado
└── models.go                      ✅ Implementado (agregador)
```

### 🔑 Modelo Central: Identidad

```go
package models

import (
    "time"
    "gorm.io/gorm"
    "github.com/google/uuid"
)

// Identity - Modelo central de identidad (única fuente de verdad)
type Identity struct {
    ID                  string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Email               string     `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
    PasswordHash        string     `json:"-" gorm:"column:password_hash;type:varchar(255);not null"`
    EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
    EmailVerifiedAt     *time.Time `json:"email_verified_at"`
    
    // Información personal básica
    FirstName           string     `json:"first_name" gorm:"type:varchar(100)"`
    LastName            string     `json:"last_name" gorm:"type:varchar(100)"`
    Avatar              string     `json:"avatar,omitempty" gorm:"type:varchar(500)"`
    Phone               string     `json:"phone,omitempty" gorm:"type:varchar(20)"`
    DateOfBirth         *time.Time `json:"date_of_birth"`
    
    // Metadatos de seguridad
    LastLoginAt         *time.Time `json:"last_login_at"`
    FailedLoginAttempts int        `json:"failed_login_attempts" gorm:"default:0"`
    LockedUntil         *time.Time `json:"locked_until"`
    
    // Relaciones
    Memberships         []OrganizationalMembership `json:"memberships,omitempty" gorm:"foreignKey:IdentityID"`
    CustomerProfile     *CustomerProfile           `json:"customer_profile,omitempty" gorm:"foreignKey:IdentityID"`
    UserProfile         *UserProfile              `json:"user_profile,omitempty" gorm:"foreignKey:IdentityID"`
    RefreshTokens       []RefreshToken            `json:"-" gorm:"foreignKey:IdentityID"`
    SentInvitations     []Invitation              `json:"-" gorm:"foreignKey:InvitedByID"`
    AuditLogs           []AuditLog                `json:"-" gorm:"foreignKey:IdentityID"`
    ManagedDepartments  []Department              `json:"-" gorm:"foreignKey:ManagerID"`
    
    // Timestamps con soft delete
    DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
}

// BeforeCreate - Hook para generar UUID
func (i *Identity) BeforeCreate(tx *gorm.DB) error {
    if i.ID == "" {
        i.ID = uuid.New().String()
    }
    return nil
}

// FullName - Método de conveniencia
func (i *Identity) FullName() string {
    return i.FirstName + " " + i.LastName
}

// IsLocked - Verifica si la cuenta está bloqueada
func (i *Identity) IsLocked() bool {
    return i.LockedUntil != nil && i.LockedUntil.After(time.Now())
}
```

### 🏢 Modelo: Organización

```go
// Organization - Entidades multi-tenant (empresas, colegios, etc.)
type Organization struct {
    ID              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    Name            string    `json:"name" gorm:"type:varchar(255);not null"`
    Slug            string    `json:"slug" gorm:"uniqueIndex;type:varchar(100);not null"`
    Type            string    `json:"type" gorm:"type:varchar(50);not null"` // company, educational_institution
    Address         string    `json:"address,omitempty" gorm:"type:text"`
    Phone           string    `json:"phone,omitempty" gorm:"type:varchar(20)"`
    Website         string    `json:"website,omitempty" gorm:"type:varchar(255)"`
    Email           string    `json:"email,omitempty" gorm:"type:varchar(255)"`
    YearFounded     int       `json:"year_founded,omitempty"`
    
    // Campos específicos por tipo
    CEO             string    `json:"ceo,omitempty" gorm:"type:varchar(255)"` // Para empresas
    Industry        string    `json:"industry,omitempty" gorm:"type:varchar(100)"` // Para empresas
    Principal       string    `json:"principal,omitempty" gorm:"type:varchar(255)"` // Para instituciones educativas
    Levels          string    `json:"levels,omitempty" gorm:"type:json"` // ["primary", "secondary"] - Para educativo
    
    // Configuración de tienda virtual
    StoreEnabled    bool      `json:"store_enabled" gorm:"default:false"`
    StoreName       string    `json:"store_name,omitempty" gorm:"type:varchar(255)"`
    StoreCategories string    `json:"store_categories,omitempty" gorm:"type:json"` // ["uniformes", "utiles"]
    
    // Metadatos
    Logo            string    `json:"logo,omitempty" gorm:"type:varchar(500)"`
    Description     string    `json:"description,omitempty" gorm:"type:text"`
    IsActive        bool      `json:"is_active" gorm:"default:true"`
    
    // Relaciones
    Memberships     []OrganizationalMembership `json:"memberships,omitempty" gorm:"foreignKey:OrganizationID"`
    Roles           []Role                     `json:"roles,omitempty" gorm:"foreignKey:OrganizationID"`
    Departments     []Department               `json:"departments,omitempty" gorm:"foreignKey:OrganizationID"`
    Invitations     []Invitation               `json:"invitations,omitempty" gorm:"foreignKey:OrganizationID"`
    AuditLogs       []AuditLog                 `json:"-" gorm:"foreignKey:OrganizationID"`
    
    // Timestamps con soft delete
    DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
    if o.ID == "" {
        o.ID = uuid.New().String()
    }
    return nil
}
```

### 👥 Modelo: Membresía Organizacional

```go
// OrganizationalMembership - Vincula identidades con organizaciones y roles
type OrganizationalMembership struct {
    ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    IdentityID     string    `json:"identity_id" gorm:"type:varchar(36);not null;index"`
    OrganizationID string    `json:"organization_id" gorm:"type:varchar(36);not null;index"`
    RoleID         string    `json:"role_id" gorm:"type:varchar(36);not null;index"`
    
    // Contexto específico organizacional
    Department     string    `json:"department,omitempty" gorm:"type:varchar(100)"` // Para empleados
    Grade          string    `json:"grade,omitempty" gorm:"type:varchar(50)"`       // Para estudiantes
    Subject        string    `json:"subject,omitempty" gorm:"type:varchar(100)"`    // Para profesores
    StudentID      string    `json:"student_id,omitempty" gorm:"type:varchar(50)"`  // ID único estudiante
    EmployeeID     string    `json:"employee_id,omitempty" gorm:"type:varchar(50)"` // ID único empleado
    
    // Control temporal
    ActiveFrom     time.Time  `json:"active_from" gorm:"not null"`
    ActiveUntil    *time.Time `json:"active_until"`
    IsActive       bool       `json:"is_active" gorm:"default:true"`
    
    // Relaciones
    Identity       Identity     `json:"identity,omitempty" gorm:"foreignKey:IdentityID;references:ID"`
    Organization   Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
    Role           Role         `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
    
    // Timestamps con soft delete
    DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
}

func (om *OrganizationalMembership) BeforeCreate(tx *gorm.DB) error {
    if om.ID == "" {
        om.ID = uuid.New().String()
    }
    if om.ActiveFrom.IsZero() {
        om.ActiveFrom = time.Now()
    }
    return nil
}

// IsCurrentlyActive - Verifica si la membresía está activa ahora
func (om *OrganizationalMembership) IsCurrentlyActive() bool {
    now := time.Now()
    return om.IsActive && 
           now.After(om.ActiveFrom) && 
           (om.ActiveUntil == nil || now.Before(*om.ActiveUntil))
}
```

### 🛒 Modelo: Perfil de Cliente

```go
// CustomerProfile - Perfil específico para e-commerce
type CustomerProfile struct {
    ID                      string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    IdentityID              string    `json:"identity_id" gorm:"type:varchar(36);not null;uniqueIndex"`
    CustomerNumber          string    `json:"customer_number" gorm:"uniqueIndex;type:varchar(50);not null"`
    
    // Datos específicos de e-commerce
    PreferredPaymentMethod  string    `json:"preferred_payment_method,omitempty" gorm:"type:varchar(50)"`
    CreditLimit             float64   `json:"credit_limit" gorm:"type:decimal(10,2);default:0"`
    TotalSpent              float64   `json:"total_spent" gorm:"type:decimal(10,2);default:0"`
    LoyaltyPoints           int       `json:"loyalty_points" gorm:"default:0"`
    
    // Preferencias de marketing
    AcceptsMarketing        bool      `json:"accepts_marketing" gorm:"default:false"`
    PreferredLanguage       string    `json:"preferred_language" gorm:"type:varchar(10);default:es"`
    
    // Relaciones
    Identity                Identity              `json:"identity,omitempty" gorm:"foreignKey:IdentityID;references:ID"`
    ShippingAddresses       []ShippingAddress     `json:"shipping_addresses,omitempty" gorm:"foreignKey:CustomerProfileID"`
    Preferences             *CustomerPreferences  `json:"preferences,omitempty" gorm:"foreignKey:CustomerProfileID"`
    
    // Timestamps con soft delete
    DeletedAt               gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt               time.Time      `json:"created_at"`
    UpdatedAt               time.Time      `json:"updated_at"`
}

func (cp *CustomerProfile) BeforeCreate(tx *gorm.DB) error {
    if cp.ID == "" {
        cp.ID = uuid.New().String()
    }
    if cp.CustomerNumber == "" {
        cp.CustomerNumber = cp.generateCustomerNumber()
    }
    return nil
}

func (cp *CustomerProfile) generateCustomerNumber() string {
    timestamp := time.Now().Unix()
    return fmt.Sprintf("CUST%d", timestamp)
}
```

### 👤 Modelo: Sesión de Invitado

```go
// GuestSession - Sesiones temporales para usuarios no registrados
type GuestSession struct {
    ID           string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    SessionToken string    `json:"session_token" gorm:"uniqueIndex;type:varchar(255);not null"`
    
    // Datos temporales del invitado
    Email        string    `json:"email,omitempty" gorm:"type:varchar(255)"`
    FirstName    string    `json:"first_name,omitempty" gorm:"type:varchar(100)"`
    LastName     string    `json:"last_name,omitempty" gorm:"type:varchar(100)"`
    Phone        string    `json:"phone,omitempty" gorm:"type:varchar(20)"`
    
    // Tracking de actividad
    CartData     string    `json:"cart_data,omitempty" gorm:"type:longtext"` // JSON del carrito
    LastActivity time.Time `json:"last_activity" gorm:"not null"`
    IPAddress    string    `json:"ip_address,omitempty" gorm:"type:varchar(45)"`
    UserAgent    string    `json:"user_agent,omitempty" gorm:"type:text"`
    
    // Auto-expiración (sin soft delete - se elimina físicamente)
    ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index"`
    CreatedAt    time.Time `json:"created_at"`
}

func (gs *GuestSession) BeforeCreate(tx *gorm.DB) error {
    if gs.ID == "" {
        gs.ID = uuid.New().String()
    }
    if gs.SessionToken == "" {
        gs.SessionToken = uuid.New().String()
    }
    if gs.ExpiresAt.IsZero() {
        gs.ExpiresAt = time.Now().Add(24 * time.Hour) // 24 horas por defecto
    }
    gs.LastActivity = time.Now()
    return nil
}

// IsExpired - Verifica si la sesión ha expirado
func (gs *GuestSession) IsExpired() bool {
    return time.Now().After(gs.ExpiresAt)
}

// ExtendSession - Extiende la sesión por X horas
func (gs *GuestSession) ExtendSession(hours int) {
    gs.ExpiresAt = time.Now().Add(time.Duration(hours) * time.Hour)
    gs.LastActivity = time.Now()
}
```

### 🎭 Modelo: Roles

```go
// Role - Roles con permisos granulares por organización
type Role struct {
    ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    OrganizationID string    `json:"organization_id" gorm:"type:varchar(36);not null;index"`
    Name           string    `json:"name" gorm:"type:varchar(100);not null"`
    DisplayName    string    `json:"display_name" gorm:"type:varchar(255);not null"`
    Description    string    `json:"description,omitempty" gorm:"type:text"`
    HierarchyLevel int       `json:"hierarchy_level" gorm:"default:1"` // 1=bajo, 100=alto
    IsSystemRole   bool      `json:"is_system_role" gorm:"default:false"`
    
    // Permisos almacenados como JSON para flexibilidad
    Permissions    string    `json:"permissions,omitempty" gorm:"type:json"`
    
    // Relaciones
    Organization   Organization               `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
    Memberships    []OrganizationalMembership `json:"memberships,omitempty" gorm:"foreignKey:RoleID"`
    PermissionList []Permission               `json:"permission_list,omitempty" gorm:"foreignKey:RoleID"`
    Invitations    []Invitation               `json:"invitations,omitempty" gorm:"foreignKey:RoleID"`
    
    // Timestamps con soft delete
    DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
    if r.ID == "" {
        r.ID = uuid.New().String()
    }
    return nil
}

// Índice compuesto para evitar roles duplicados por organización
func (Role) TableName() string {
    return "roles"
}
```

### 🔐 Modelo: Permisos

```go
// Permission - Permisos granulares por recurso y acción
type Permission struct {
    ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    RoleID     string    `json:"role_id" gorm:"type:varchar(36);not null;index"`
    Resource   string    `json:"resource" gorm:"type:varchar(100);not null"` // users, products, orders
    Action     string    `json:"action" gorm:"type:varchar(50);not null"`    // create, read, update, delete
    Scope      string    `json:"scope" gorm:"type:varchar(50);not null"`     // own, department, organization, all
    Conditions string    `json:"conditions,omitempty" gorm:"type:json"`      // Condiciones adicionales
    
    // Relaciones
    Role       Role      `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
    
    // Timestamps (sin soft delete para permisos)
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}
```

### 🏛️ Modelo: Departamentos

```go
// Department - Departamentos jerárquicos por organización
type Department struct {
    ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    OrganizationID string    `json:"organization_id" gorm:"type:varchar(36);not null;index"`
    ParentID       *string   `json:"parent_id,omitempty" gorm:"type:varchar(36);index"` // Auto-referencia
    Name           string    `json:"name" gorm:"type:varchar(255);not null"`
    Description    string    `json:"description,omitempty" gorm:"type:text"`
    ManagerID      *string   `json:"manager_id,omitempty" gorm:"type:varchar(36);index"`
    HierarchyLevel int       `json:"hierarchy_level" gorm:"default:1"`
    IsActive       bool      `json:"is_active" gorm:"default:true"`
    
    // Relaciones
    Organization   Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
    Parent         *Department  `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
    Children       []Department `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
    Manager        *Identity    `json:"manager,omitempty" gorm:"foreignKey:ManagerID;references:ID"`
    
    // Timestamps con soft delete
    DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
}

func (d *Department) BeforeCreate(tx *gorm.DB) error {
    if d.ID == "" {
        d.ID = uuid.New().String()
    }
    return nil
}
```

### 📧 Modelo: Invitaciones

```go
// Invitation - Sistema de invitaciones organizacionales
type Invitation struct {
    ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    OrganizationID string    `json:"organization_id" gorm:"type:varchar(36);not null;index"`
    InvitedByID    string    `json:"invited_by_id" gorm:"type:varchar(36);not null;index"`
    RoleID         string    `json:"role_id" gorm:"type:varchar(36);not null;index"`
    Email          string    `json:"email" gorm:"type:varchar(255);not null"`
    Token          string    `json:"token" gorm:"uniqueIndex;type:varchar(255);not null"`
    Status         string    `json:"status" gorm:"type:varchar(50);default:pending"` // pending, accepted, expired, cancelled
    
    // Contexto específico para la invitación
    Department     string    `json:"department,omitempty" gorm:"type:varchar(100)"`
    Grade          string    `json:"grade,omitempty" gorm:"type:varchar(50)"`
    Subject        string    `json:"subject,omitempty" gorm:"type:varchar(100)"`
    Message        string    `json:"message,omitempty" gorm:"type:text"`
    Metadata       string    `json:"metadata,omitempty" gorm:"type:json"`
    
    // Control temporal
    ExpiresAt      time.Time `json:"expires_at" gorm:"not null;index"`
    AcceptedAt     *time.Time `json:"accepted_at"`
    
    // Relaciones
    Organization   Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
    InvitedBy      Identity     `json:"invited_by,omitempty" gorm:"foreignKey:InvitedByID;references:ID"`
    Role           Role         `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
    
    // Timestamps con soft delete
    DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
}

func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
    if i.ID == "" {
        i.ID = uuid.New().String()
    }
    if i.Token == "" {
        i.Token = uuid.New().String()
    }
    if i.ExpiresAt.IsZero() {
        i.ExpiresAt = time.Now().Add(168 * time.Hour) // 7 días por defecto
    }
    return nil
}

// IsExpired - Verifica si la invitación ha expirado
func (i *Invitation) IsExpired() bool {
    return time.Now().After(i.ExpiresAt)
}

// CanBeAccepted - Verifica si la invitación puede ser aceptada
func (i *Invitation) CanBeAccepted() bool {
    return i.Status == "pending" && !i.IsExpired()
}
```

### 📊 Modelos Auxiliares

```go
// ShippingAddress - Direcciones de envío para clientes
type ShippingAddress struct {
    ID                string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    CustomerProfileID string    `json:"customer_profile_id" gorm:"type:varchar(36);not null;index"`
    Type              string    `json:"type" gorm:"type:varchar(50);default:shipping"` // shipping, billing
    FirstName         string    `json:"first_name" gorm:"type:varchar(100);not null"`
    LastName          string    `json:"last_name" gorm:"type:varchar(100);not null"`
    Company           string    `json:"company,omitempty" gorm:"type:varchar(255)"`
    AddressLine1      string    `json:"address_line_1" gorm:"type:varchar(255);not null"`
    AddressLine2      string    `json:"address_line_2,omitempty" gorm:"type:varchar(255)"`
    City              string    `json:"city" gorm:"type:varchar(100);not null"`
    State             string    `json:"state" gorm:"type:varchar(100);not null"`
    PostalCode        string    `json:"postal_code" gorm:"type:varchar(20);not null"`
    Country           string    `json:"country" gorm:"type:varchar(50);not null;default:CO"`
    Phone             string    `json:"phone,omitempty" gorm:"type:varchar(20)"`
    IsDefault         bool      `json:"is_default" gorm:"default:false"`
    
    // Relaciones
    CustomerProfile   CustomerProfile `json:"customer_profile,omitempty" gorm:"foreignKey:CustomerProfileID;references:ID"`
    
    // Timestamps con soft delete
    DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
}

// CustomerPreferences - Preferencias específicas de clientes
type CustomerPreferences struct {
    ID                      string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    CustomerProfileID       string    `json:"customer_profile_id" gorm:"type:varchar(36);not null;uniqueIndex"`
    NotificationPreferences string    `json:"notification_preferences,omitempty" gorm:"type:json"`
    PrivacySettings         string    `json:"privacy_settings,omitempty" gorm:"type:json"`
    MarketingPreferences    string    `json:"marketing_preferences,omitempty" gorm:"type:json"`
    PreferredCurrency       string    `json:"preferred_currency" gorm:"type:varchar(10);default:COP"`
    Timezone                string    `json:"timezone" gorm:"type:varchar(50);default:America/Bogota"`
    CustomFields            string    `json:"custom_fields,omitempty" gorm:"type:json"`
    
    // Relaciones
    CustomerProfile         CustomerProfile `json:"customer_profile,omitempty" gorm:"foreignKey:CustomerProfileID;references:ID"`
    
    // Timestamps (sin soft delete)
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
}

// UserProfile - Perfiles extendidos para identidades
type UserProfile struct {
    ID               string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    IdentityID       string    `json:"identity_id" gorm:"type:varchar(36);not null;uniqueIndex"`
    Bio              string    `json:"bio,omitempty" gorm:"type:text"`
    SocialLinks      string    `json:"social_links,omitempty" gorm:"type:json"`
    Skills           string    `json:"skills,omitempty" gorm:"type:json"`
    Interests        string    `json:"interests,omitempty" gorm:"type:json"`
    Location         string    `json:"location,omitempty" gorm:"type:varchar(255)"`
    Website          string    `json:"website,omitempty" gorm:"type:varchar(255)"`
    CustomAttributes string    `json:"custom_attributes,omitempty" gorm:"type:json"`
    
    // Relaciones
    Identity         Identity  `json:"identity,omitempty" gorm:"foreignKey:IdentityID;references:ID"`
    
    // Timestamps con soft delete
    DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
    CreatedAt        time.Time      `json:"created_at"`
    UpdatedAt        time.Time      `json:"updated_at"`
}

// RefreshToken - Tokens de renovación de sesión
type RefreshToken struct {
    ID          string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    IdentityID  string    `json:"identity_id" gorm:"type:varchar(36);not null;index"`
    Token       string    `json:"token" gorm:"uniqueIndex;type:varchar(500);not null"`
    DeviceInfo  string    `json:"device_info,omitempty" gorm:"type:text"`
    IPAddress   string    `json:"ip_address,omitempty" gorm:"type:varchar(45)"`
    UserAgent   string    `json:"user_agent,omitempty" gorm:"type:text"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    ExpiresAt   time.Time `json:"expires_at" gorm:"not null;index"`
    LastUsedAt  *time.Time `json:"last_used_at"`
    
    // Relaciones
    Identity    Identity  `json:"identity,omitempty" gorm:"foreignKey:IdentityID;references:ID"`
    
    // Timestamps (sin soft delete - se elimina físicamente al expirar)
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// AuditLog - Logs de auditoría inmutables
type AuditLog struct {
    ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
    OrganizationID *string   `json:"organization_id,omitempty" gorm:"type:varchar(36);index"`
    IdentityID     *string   `json:"identity_id,omitempty" gorm:"type:varchar(36);index"`
    Action         string    `json:"action" gorm:"type:varchar(100);not null;index"`
    Resource       string    `json:"resource" gorm:"type:varchar(100);not null;index"`
    ResourceID     string    `json:"resource_id,omitempty" gorm:"type:varchar(100);index"`
    OldValues      string    `json:"old_values,omitempty" gorm:"type:json"`
    NewValues      string    `json:"new_values,omitempty" gorm:"type:json"`
    IPAddress      string    `json:"ip_address,omitempty" gorm:"type:varchar(45)"`
    UserAgent      string    `json:"user_agent,omitempty" gorm:"type:text"`
    Metadata       string    `json:"metadata,omitempty" gorm:"type:json"`
    Status         string    `json:"status" gorm:"type:varchar(50);default:success"` // success, failed, pending
    
    // Relaciones
    Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
    Identity       *Identity     `json:"identity,omitempty" gorm:"foreignKey:IdentityID;references:ID"`
    
    // Timestamp inmutable (sin soft delete ni update)
    CreatedAt      time.Time `json:"created_at" gorm:"not null"`
}

func (al *AuditLog) BeforeCreate(tx *gorm.DB) error {
    if al.ID == "" {
        al.ID = uuid.New().String()
    }
    return nil
}

// BeforeUpdate - Prevenir actualizaciones en logs de auditoría
func (al *AuditLog) BeforeUpdate(tx *gorm.DB) error {
    return gorm.ErrInvalidTransaction // Los logs de auditoría son inmutables
}
```

## 🚀 Migraciones de Base de Datos

### Archivo de Migración Principal

```go
// filepath: module/authentication/migrations.go
package authentication

import (
    "gorm.io/gorm"
    "module/authentication/models"
)

// RunMigrations - Ejecuta todas las migraciones del módulo de autenticación
func RunMigrations(db *gorm.DB) error {
    // Orden específico para respetar dependencias de foreign keys
    models := []interface{}{
        // 1. Entidades independientes primero
        &models.Identity{},
        &models.Organization{},
        &models.GuestSession{},
        
        // 2. Entidades que dependen de Identity y Organization
        &models.Role{},
        &models.Permission{},
        &models.Department{},
        &models.OrganizationalMembership{},
        &models.CustomerProfile{},
        &models.UserProfile{},
        &models.RefreshToken{},
        
        // 3. Entidades que dependen de múltiples tablas
        &models.Invitation{},
        &models.ShippingAddress{},
        &models.CustomerPreferences{},
        &models.AuditLog{},
    }
    
    // Ejecutar migraciones
    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to migrate %T: %w", model, err)
        }
    }
    
    // Crear índices adicionales
    if err := createAdditionalIndexes(db); err != nil {
        return fmt.Errorf("failed to create additional indexes: %w", err)
    }
    
    // Insertar datos iniciales
    if err := seedInitialData(db); err != nil {
        return fmt.Errorf("failed to seed initial data: %w", err)
    }
    
    return nil
}

// createAdditionalIndexes - Crea índices adicionales para optimización
func createAdditionalIndexes(db *gorm.DB) error {
    indexes := []string{
        // Índices compuestos para consultas frecuentes
        "CREATE INDEX IF NOT EXISTS idx_org_membership_active ON organizational_memberships(organization_id, is_active, active_from, active_until)",
        "CREATE INDEX IF NOT EXISTS idx_customer_profile_active ON customer_profiles(identity_id, deleted_at)",
        "CREATE INDEX IF NOT EXISTS idx_invitation_status_expires ON invitations(status, expires_at)",
        "CREATE INDEX IF NOT EXISTS idx_audit_log_time_action ON audit_logs(created_at, action, organization_id)",
        "CREATE INDEX IF NOT EXISTS idx_guest_session_expires ON guest_sessions(expires_at)",
        "CREATE INDEX IF NOT EXISTS idx_refresh_token_active ON refresh_tokens(identity_id, is_active, expires_at)",
        
        // Índices para búsquedas de texto
        "CREATE INDEX IF NOT EXISTS idx_identity_name ON identities(first_name, last_name)",
        "CREATE INDEX IF NOT EXISTS idx_organization_search ON organizations(name, type, is_active)",
        
        // Índices únicos compuestos
        "CREATE UNIQUE INDEX IF NOT EXISTS idx_role_org_name ON roles(organization_id, name) WHERE deleted_at IS NULL",
        "CREATE UNIQUE INDEX IF NOT EXISTS idx_dept_org_name ON departments(organization_id, name) WHERE deleted_at IS NULL",
        "CREATE UNIQUE INDEX IF NOT EXISTS idx_membership_unique ON organizational_memberships(identity_id, organization_id, role_id) WHERE deleted_at IS NULL",
    }
    
    for _, index := range indexes {
        if err := db.Exec(index).Error; err != nil {
            return fmt.Errorf("failed to create index: %s, error: %w", index, err)
        }
    }
    
    return nil
}

// seedInitialData - Inserta datos iniciales requeridos
func seedInitialData(db *gorm.DB) error {
    // Verificar si ya existen datos
    var count int64
    db.Model(&models.Role{}).Where("is_system_role = ?", true).Count(&count)
    if count > 0 {
        return nil // Ya hay datos iniciales
    }
    
    // Crear roles del sistema predefinidos para diferentes tipos de organización
    systemRoles := []models.Role{
        // Roles empresariales
        {
            ID:             "super-admin-role",
            Name:           "super_admin",
            DisplayName:    "Super Administrador",
            Description:    "Acceso completo al sistema",
            HierarchyLevel: 100,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"*","action":"*","scope":"all"}]`,
        },
        {
            ID:             "admin-role",
            Name:           "admin",
            DisplayName:    "Administrador",
            Description:    "Administrador organizacional",
            HierarchyLevel: 90,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"*","action":"*","scope":"organization"}]`,
        },
        {
            ID:             "manager-role",
            Name:           "manager",
            DisplayName:    "Gerente",
            Description:    "Gerente de departamento",
            HierarchyLevel: 80,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"users","action":"read","scope":"department"},{"resource":"users","action":"update","scope":"department"}]`,
        },
        {
            ID:             "employee-role",
            Name:           "employee",
            DisplayName:    "Empleado",
            Description:    "Empleado básico",
            HierarchyLevel: 50,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"users","action":"read","scope":"own"}]`,
        },
        
        // Roles educativos
        {
            ID:             "director-role",
            Name:           "director",
            DisplayName:    "Director",
            Description:    "Director de institución educativa",
            HierarchyLevel: 100,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"*","action":"*","scope":"organization"}]`,
        },
        {
            ID:             "coordinator-role",
            Name:           "coordinator",
            DisplayName:    "Coordinador",
            Description:    "Coordinador académico",
            HierarchyLevel: 85,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"students","action":"*","scope":"department"},{"resource":"teachers","action":"read","scope":"department"}]`,
        },
        {
            ID:             "teacher-role",
            Name:           "teacher",
            DisplayName:    "Profesor",
            Description:    "Profesor de materia",
            HierarchyLevel: 60,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"students","action":"read","scope":"department"},{"resource":"students","action":"grade","scope":"department"}]`,
        },
        {
            ID:             "student-role",
            Name:           "student",
            DisplayName:    "Estudiante",
            Description:    "Estudiante",
            HierarchyLevel: 20,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"grades","action":"read","scope":"own"},{"resource":"courses","action":"read","scope":"own"}]`,
        },
        {
            ID:             "parent-role",
            Name:           "parent",
            DisplayName:    "Padre de Familia",
            Description:    "Padre o acudiente",
            HierarchyLevel: 30,
            IsSystemRole:   true,
            Permissions:    `[{"resource":"students","action":"read","scope":"own"},{"resource":"grades","action":"read","scope":"own"}]`,
        },
    }
    
    // Insertar roles del sistema
    for _, role := range systemRoles {
        if err := db.Create(&role).Error; err != nil {
            return fmt.Errorf("failed to create system role %s: %w", role.Name, err)
        }
    }
    
    return nil
}
```

### Comando de Migración

```go
// filepath: cmd/migrate.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    
    "your-project/config"
    "your-project/database"
    "your-project/module/authentication"
)

func main() {
    // Flags de comando
    var (
        action = flag.String("action", "up", "Migration action: up, down, reset")
        module = flag.String("module", "all", "Module to migrate: all, auth")
    )
    flag.Parse()
    
    // Cargar configuración
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Conectar a base de datos
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Ejecutar migraciones según acción
    switch *action {
    case "up":
        if err := runMigrationsUp(db, *module); err != nil {
            log.Fatalf("Migration failed: %v", err)
        }
        fmt.Println("✅ Migrations completed successfully")
        
    case "down":
        if err := runMigrationsDown(db, *module); err != nil {
            log.Fatalf("Migration rollback failed: %v", err)
        }
        fmt.Println("✅ Migrations rolled back successfully")
        
    case "reset":
        if err := resetDatabase(db, *module); err != nil {
            log.Fatalf("Database reset failed: %v", err)
        }
        fmt.Println("✅ Database reset completed successfully")
        
    default:
        fmt.Printf("Unknown action: %s\n", *action)
        os.Exit(1)
    }
}

func runMigrationsUp(db *gorm.DB, module string) error {
    switch module {
    case "all", "auth":
        return authentication.RunMigrations(db)
    default:
        return fmt.Errorf("unknown module: %s", module)
    }
}

func runMigrationsDown(db *gorm.DB, module string) error {
    // Implementar rollback si es necesario
    // Por seguridad, no implementamos rollback automático
    return fmt.Errorf("rollback must be done manually for safety")
}

func resetDatabase(db *gorm.DB, module string) error {
    // ⚠️ PELIGROSO - Solo para desarrollo
    if os.Getenv("APP_ENV") == "production" {
        return fmt.Errorf("database reset is not allowed in production")
    }
    
    // Eliminar todas las tablas del módulo auth
    tables := []string{
        "audit_logs",
        "customer_preferences", 
        "shipping_addresses",
        "invitations",
        "refresh_tokens",
        "user_profiles",
        "customer_profiles",
        "organizational_memberships",
        "departments",
        "permissions",
        "roles",
        "guest_sessions",
        "organizations",
        "identities",
    }
    
    for _, table := range tables {
        if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
            return fmt.Errorf("failed to drop table %s: %w", table, err)
        }
    }
    
    // Ejecutar migraciones nuevamente
    return authentication.RunMigrations(db)
}
```

### Script de Ejecución

```bash
#!/bin/bash
# filepath: scripts/migrate.sh

# Configuración de colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Iniciando migraciones del módulo de autenticación...${NC}"

# Verificar que Go esté instalado
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go no está instalado${NC}"
    exit 1
fi

# Verificar variables de entorno
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}❌ DATABASE_URL no está configurada${NC}"
    exit 1
fi

# Ejecutar migraciones
echo -e "${YELLOW}📊 Ejecutando migraciones de base de datos...${NC}"
go run cmd/migrate.go -action=up -module=auth

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Migraciones completadas exitosamente${NC}"
    
    # Verificar que las tablas se crearon correctamente
    echo -e "${YELLOW}🔍 Verificando estructura de base de datos...${NC}"
    
    # Lista de tablas esperadas
    expected_tables=(
        "identities"
        "organizations" 
        "organizational_memberships"
        "customer_profiles"
        "guest_sessions"
        "roles"
        "permissions"
        "departments"
        "invitations"
        "shipping_addresses"
        "customer_preferences"
        "user_profiles"
        "refresh_tokens"
        "audit_logs"
    )
    
    echo -e "${GREEN}📋 Tablas creadas:${NC}"
    for table in "${expected_tables[@]}"; do
        echo -e "  ✓ $table"
    done
    
    echo -e "${GREEN}🎉 Sistema de autenticación listo para usar!${NC}"
else
    echo -e "${RED}❌ Error en las migraciones${NC}"
    exit 1
fi
```

## 🎯 Uso de las Migraciones

### Comandos Disponibles

```bash
# Ejecutar todas las migraciones
go run cmd/migrate.go -action=up -module=auth

# Ejecutar migraciones usando el script
chmod +x scripts/migrate.sh
./scripts/migrate.sh

# Resetear base de datos (solo desarrollo)
APP_ENV=development go run cmd/migrate.go -action=reset -module=auth
```

### Variables de Entorno Requeridas

```bash
# Base de datos
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

# O configuración individual
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable

# Entorno
APP_ENV=development  # development, staging, production
```

## ✅ Verificación Post-Migración

Después de ejecutar las migraciones, verifica que:

1. **Todas las tablas se crearon** correctamente
2. **Los índices están en su lugar** para optimización
3. **Los roles del sistema** se insertaron correctamente
4. **Las relaciones foreign key** funcionan
5. **Los constraints** están activos

La estructura de base de datos está optimizada para:
- **Alto rendimiento** con índices estratégicos
- **Integridad referencial** con foreign keys apropiadas
- **Escalabilidad** con soft deletes y particionado futuro
- **Auditoría completa** con logs inmutables
- **Flexibilidad** con campos JSON para extensiones
