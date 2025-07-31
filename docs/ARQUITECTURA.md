# Arquitectura del Sistema - Proyecto Go
## Stack Tecnológico: MySQL + GORM + Fiber

## Tabla de Contenidos
1. [Visión General](#visión-general)
2. [Stack Tecnológico](#stack-tecnológico)
3. [Arquitectura de Alto Nivel](#arquitectura-de-alto-nivel)
4. [Estructura de Directorios](#estructura-de-directorios)
5. [Componentes del Sistema](#componentes-del-sistema)
6. [Patrón de Módulos](#patrón-de-módulos)
7. [Guía de Implementación de Nuevos Módulos](#guía-de-implementación-de-nuevos-módulos)
8. [Flujo de Datos](#flujo-de-datos)
9. [Configuración y Variables de Entorno](#configuración-y-variables-de-entorno)
10. [Base de Datos y Migraciones](#base-de-datos-y-migraciones)
11. [Autenticación y Autorización](#autenticación-y-autorización)
12. [Caché y Colas](#caché-y-colas)
13. [Mejores Prácticas](#mejores-prácticas)

## Visión General

Este proyecto implementa una **arquitectura modular moderna** basada en Go, utilizando **Fiber** como framework web de alto rendimiento, **GORM** como ORM avanzado para **MySQL**, y un sistema de módulos auto-registrables para máxima escalabilidad.

### Características Principales
- **🏗️ Arquitectura Modular**: Cada funcionalidad es un módulo independiente auto-contenido
- **🔄 Registro Dinámico**: Los módulos se registran automáticamente usando el patrón `init()`
- **⚙️ Configuración Centralizada**: Gestión robusta de configuración con variables de entorno
- **🗄️ Base de Datos MySQL + GORM**: ORM moderno con migraciones automáticas y relaciones complejas
- **🔐 Autenticación JWT Avanzada**: Sistema completo con refresh tokens, roles y contextos organizacionales
- **⚡ Fiber Framework**: Web framework ultra-rápido inspirado en Express.js
- **📦 Caché Redis**: Sistema de caché distribuido con RediSearch
- **🐰 Colas RabbitMQ**: Sistema de mensajería asíncrona para operaciones en background
- **🛡️ Middleware Extensible**: Sistema de middleware configurable y personalizable

## Stack Tecnológico

### Core Technologies
- **🔷 Go 1.23+**: Lenguaje principal con toolchain 1.24.4
- **🚀 Fiber v2**: Framework web de alto rendimiento
- **🗄️ MySQL 8.0+**: Base de datos relacional principal
- **📦 GORM v1.30**: ORM moderno con soporte completo para MySQL

### Dependencias Principales
```go
// Web Framework & HTTP
github.com/gofiber/fiber/v2 v2.52.6

// Database & ORM  
gorm.io/driver/mysql v1.5.7
gorm.io/gorm v1.30.0

// Caché & Search
github.com/redis/go-redis/v9 v9.7.1
github.com/RediSearch/redisearch-go v1.1.1

// Message Queue
github.com/rabbitmq/amqp091-go v1.10.0

// Authentication & Security
github.com/golang-jwt/jwt/v5 v5.3.0
golang.org/x/crypto v0.40.0

// Configuration & Utils
github.com/joho/godotenv v1.5.1
github.com/google/uuid v1.6.0
```

## Arquitectura de Alto Nivel

```
┌─────────────────────────────────────────────────────────┐
│                   HTTP/HTTPS Requests                   │
│                (REST API Endpoints)                     │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 🚀 Fiber Web Server                     │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │    CORS     │ │   Logger    │ │  Security Headers   ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │JWT Validation│ │Rate Limiting│ │   Custom Middleware ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│            📦 Dynamic Module Registry                   │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │Authentication│ │   User      │ │   Organization      ││
│  │   Module     │ │  Module     │ │     Module          ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │  Customer   │ │   Role      │ │    Module N         ││
│  │   Module    │ │  Module     │ │   (Auto-Register)   ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│              🏢 Business Logic Layer                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │   Handlers  │ │  Services   │ │   Repositories      ││
│  │(Controllers)│ │(Use Cases)  │ │  (Data Access)      ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐│
│  │    DTOs     │ │ Validations │ │    Models           ││
│  │(Data Transfer)│ │(Business)  │ │  (GORM Entities)    ││
│  └─────────────┘ └─────────────┘ └─────────────────────┘│
└──┬────────────────┬─────────────────┬──────────────────┘
   │                │                 │
┌──▼──────────────┐ ┌▼───────────────┐ ┌▼─────────────────┐
│  🗄️ MySQL      │ │  🔴 Redis       │ │  🐰 RabbitMQ     │
│   Database      │ │   Cache         │ │    Queue         │
│                 │ │                 │ │                  │
│ • Users         │ │ • Sessions      │ │ • Email Queue    │
│ • Organizations │ │ • Cache Data    │ │ • Notifications  │
│ • Roles         │ │ • Search Index  │ │ • Background     │
│ • Permissions   │ │ • Rate Limits   │ │   Tasks          │
│ • Audit Logs    │ │                 │ │                  │
└─────────────────┘ └─────────────────┘ └──────────────────┘
```

### Flujo de Procesamiento de Requests

```mermaid
sequenceDiagram
    participant C as Client
    participant F as Fiber Server
    participant M as Middleware
    participant R as Router/Module
    participant H as Handler
    participant S as Service
    participant D as Database/GORM
    participant A as Audit Log

    C->>F: HTTP Request
    F->>M: Apply Middleware Stack
    M->>M: CORS, Logging, Auth
    M->>R: Route to Module
    R->>H: Execute Handler
    H->>H: Parse & Validate DTO
    H->>S: Call Service Method
    S->>S: Business Logic
    S->>D: Database Operations
    D->>S: Return Data
    S->>A: Log Audit Event
    S->>H: Return Result
    H->>C: HTTP Response
```

## Estructura de Directorios

```
proyecto/ (practicev2)
├── cmd/                          # 🚀 Punto de entrada de la aplicación
│   ├── main.go                   # Configuración principal del servidor Fiber
│   └── uploads/                  # Almacenamiento de archivos subidos
│       └── images/               # Imágenes de usuarios, productos, etc.
│
├── config/                       # ⚙️ Configuración centralizada
│   └── config.go                 # Variables de entorno, JWT, cookies
│
├── database/                     # 🗄️ Conexión y configuración MySQL + GORM
│   └── database.go               # DSN, pool de conexiones, configuración GORM
│
├── registry/                     # 📦 Sistema de registro de módulos y servicios
│   ├── registry.go               # Registro dinámico de módulos con init()
│   ├── migrations.go             # Sistema de migraciones con prioridades
│   ├── cache/                    # 🔴 Servicios de caché
│   │   ├── redis.go              # Cliente Redis optimizado
│   │   └── redisearch.go         # Búsqueda avanzada con RediSearch
│   └── queue/                    # 🐰 Servicios de colas
│       └── rabbitmq.go           # Cliente RabbitMQ con reconexión automática
│
├── module/                       # 🏢 Módulos de la aplicación
│   ├── authentication/           # 🔐 Sistema completo de autenticación
│   │   ├── authentication.go     # Registro del módulo
│   │   ├── migrations.go         # Migraciones específicas
│   │   ├── README.md             # Documentación del módulo
│   │   ├── auth/                 # Autenticación core
│   │   │   ├── dto.go            # Data Transfer Objects
│   │   │   ├── handler.go        # Controllers HTTP/REST
│   │   │   ├── repository.go     # Acceso a datos con GORM
│   │   │   ├── routes.go         # Definición de rutas Fiber
│   │   │   └── service.go        # Lógica de negocio
│   │   ├── models/               # 📊 Modelos GORM
│   │   │   ├── identity.go       # Usuario principal
│   │   │   ├── role.go           # Roles y permisos
│   │   │   ├── organization.go   # Organizaciones
│   │   │   ├── refresh_token.go  # Tokens de renovación
│   │   │   ├── audit_log.go      # Logs de auditoría
│   │   │   └── [otros modelos]   # Perfiles, membresías, etc.
│   │   ├── utils/                # 🛠️ Utilidades del módulo
│   │   │   ├── jwt.go            # Generación y validación JWT
│   │   │   ├── password.go       # Hash y verificación
│   │   │   ├── validation.go     # Validaciones de negocio
│   │   │   ├── response.go       # Respuestas HTTP estandarizadas
│   │   │   └── [otras utils]     # Cache, email, upload, etc.
│   │   └── [otros submódulos]/   # customer/, user/, role/, etc.
│   │
│   └── middleware/               # 🛡️ Middleware personalizado
│       └── middleware.go         # Logger, CORS, Auth, Rate Limiting
│
├── addModules/                   # 🔄 Importación automática de módulos  
│   └── modules.go                # Registro automático usando init()
│
├── handlers/                     # 📝 Handlers HTTP globales (legacy/shared)
├── docs/                         # 📚 Documentación del proyecto
│   ├── ARQUITECTURA.md           # Este documento
│   ├── DATABASE_ER_DIAGRAM.md    # Diagramas de base de datos
│   ├── GUIA_DESARROLLO_MODULO_AUTH.md # Guía específica de auth
│   └── [otros docs]/             # Documentación técnica adicional
│
├── arquitectura/                 # 🏗️ Diagramas y documentación técnica
├── bin/                          # 📦 Binarios compilados
│   └── server                    # Ejecutable principal
├── docker-compose.yml            # 🐳 Servicios de desarrollo (MySQL, Redis, RabbitMQ)
├── go.mod                        # 📋 Dependencias del proyecto
└── go.sum                        # 🔒 Checksums de dependencias
```

### Convenciones de Nombrado

- **📁 Directorios**: snake_case (ej: `user_profile`, `organization_membership`)
- **📄 Archivos**: snake_case (ej: `user_handler.go`, `auth_service.go`) 
- **🏷️ Packages**: lowercase, single word (ej: `auth`, `user`, `organization`)
- **🏛️ Structs**: PascalCase (ej: `UserProfile`, `OrganizationalMembership`)
- **🔧 Functions**: camelCase (ej: `getUserByID`, `createOrganization`)
- **📊 DB Tables**: snake_case (ej: `user_profiles`, `organizational_memberships`)

## Componentes del Sistema

### 1. 🚀 Core Application (`cmd/main.go`)
El punto de entrada principal que orquesta toda la aplicación:

```go
func main() {
    // 1. Cargar configuración desde .env
    godotenv.Load()
    config.Init()
    
    // 2. Conectar a MySQL con GORM
    database.ConnectDatabase()
    
    // 3. Ejecutar migraciones ordenadas
    registry.RunMigrations(database.DBconn)
    
    // 4. Configurar Fiber con middleware stack
    app := fiber.New(fiber.Config{
        BodyLimit: 10 * 1024 * 1024, // 10MB
    })
    
    // 5. Configurar CORS para desarrollo/producción
    app.Use(cors.New(corsConfig))
    app.Use(middleware.ConfigMiddleware())
    
    // 6. Registrar todos los módulos dinámicamente
    addModules.RegisterAllModules(app)
    
    // 7. Iniciar servidor HTTP/HTTPS
    app.Listen(addr)
}
```

### 2. 📦 Dynamic Registry System (`registry/`)
Sistema centralizado para el registro automático de:

#### Módulos (`registry.go`)
```go
type ModuleRegistrar func(*fiber.App)

var moduleRegistrars []ModuleRegistrar

func RegisterModule(registrar ModuleRegistrar) {
    moduleRegistrars = append(moduleRegistrars, registrar)
}

func RegisterAllModules(app *fiber.App) {
    for _, registrar := range moduleRegistrars {
        registrar(app)
    }
}
```

#### Migraciones (`migrations.go`)  
```go
type Migration struct {
    Order int
    Func  func(*gorm.DB) error
}

func RegisterMigration(order int, migrationFunc func(*gorm.DB) error) {
    migrations = append(migrations, Migration{Order: order, Func: migrationFunc})
}

func RunMigrations(db *gorm.DB) {
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Order < migrations[j].Order
    })
    // Ejecutar migraciones en orden...
}
```

### 3. ⚙️ Configuration (`config/`)
Gestión centralizada y robusta de configuración:

#### Variables de Entorno
- **Database**: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- **JWT**: `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `JWT_ISSUER`
- **Server**: `APP_HOST`, `APP_PORT`, `ENV`
- **Cache**: `REDIS_URL`
- **Queue**: `RABBITMQ_URL`, `RABBITMQ_EXCHANGE`

#### Funciones de Configuración JWT
```go
func GetJWTAccessSecret() string    // Token de acceso (15 min)
func GetJWTRefreshSecret() string   // Token de refresh (7 días)
func GetJWTAccessTTL() time.Duration
func GetJWTRefreshTTL() time.Duration
func IsProductionCookie() bool      // Para configuración de cookies
```

### 4. 🗄️ Database Layer (`database/`)
Conexión optimizada a MySQL usando GORM:

#### Configuración de Conexión
```go
func ConnectDatabase() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        user, password, host, port, name)
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // Desactiva pluralización
        },
    })
    
    DBconn = db // Variable global
}
```

#### Características GORM
- **Singular Tables**: Nombres de tabla sin pluralización
- **Soft Deletes**: Eliminación lógica con `DeletedAt`
- **Migrations**: Auto-migración de esquemas
- **Relations**: Soporte completo para relaciones complejas
- **Hooks**: Before/After Create, Update, Delete
- **Connection Pool**: Pool de conexiones optimizado

### 5. 🏢 Module System (`module/`)
Arquitectura de módulos auto-contenidos:

#### Estructura Estándar de un Módulo
```
module/nombre_modulo/
├── init.go              # Registro automático
├── routes.go            # Definición de rutas Fiber
├── handler.go           # Controllers HTTP
├── service.go           # Lógica de negocio  
├── repository.go        # Acceso a datos GORM
├── dto.go               # Data Transfer Objects
├── models.go            # Entidades GORM
├── validation.go        # Validaciones de negocio
└── migrations.go        # Migraciones específicas
```

#### Módulo de Autenticación (`module/authentication/`)
Módulo completo de autenticación con:
- **JWT Access & Refresh Tokens**
- **Roles y Permisos**
- **Contextos Organizacionales**
- **Auditoría de Accesos**
- **Multi-tenant Support**
- **Customer Profiles**

### 6. 🛡️ Middleware Stack (`module/middleware/`)
Sistema de middleware extensible:

#### Middleware Configurado
```go
// CORS para desarrollo y producción
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://localhost:3000,http://localhost:3000",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
    AllowCredentials: true,
}))

// Logger personalizado
app.Use(logger.New(logger.Config{
    Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
}))

// Rate Limiting (configurar según necesidades)
// Security Headers
// Request ID
// Error Handling
```

### 7. 🔄 Auto-Registration (`addModules/`)
Sistema de importación automática:

```go
// addModules/modules.go
package modules

import (
    // Importar módulos para activar init()
    _ "practicev2/module/authentication"
    _ "practicev2/module/user"
    _ "practicev2/module/organization" 
    // ... más módulos
)

func RegisterAllModules(app *fiber.App) {
    registry.RegisterAllModules(app)
}
```

## Patrón de Módulos

### Arquitectura de un Módulo

Cada módulo sigue una **estructura hexagonal** (Clean Architecture) adaptada para Go + Fiber + GORM:

```
module/
└── nombre_modulo/
    ├── init.go                # 🔄 Registro automático del módulo
    ├── routes.go              # 🛣️ Definición de rutas Fiber  
    ├── handler.go             # 🎮 Controllers/Handlers HTTP
    ├── service.go             # 🏢 Lógica de negocio (Use Cases)
    ├── repository.go          # 🗄️ Acceso a datos con GORM
    ├── models.go              # 📊 Entidades/Modelos GORM
    ├── dto.go                 # 📦 Data Transfer Objects
    ├── validation.go          # ✅ Validaciones de negocio
    ├── migrations.go          # 🔄 Migraciones específicas del módulo
    └── [submódulos]/          # 📁 Organización por subdominios
```

### Principios de Diseño SOLID

1. **🔒 Single Responsibility**: Cada archivo tiene una responsabilidad específica
2. **🔓 Open/Closed**: Extensible via interfaces, cerrado para modificación
3. **🔄 Liskov Substitution**: Interfaces intercambiables
4. **⚡ Interface Segregation**: Interfaces pequeñas y específicas
5. **🔀 Dependency Inversion**: Dependencias via interfaces, no implementaciones

### Flujo de Datos en un Módulo

```
HTTP Request
    ↓
🛣️ Router (routes.go)
    ↓
🎮 Handler (handler.go)
    ↓ (parse DTO)
✅ Validation (validation.go)
    ↓
🏢 Service (service.go)
    ↓ (business logic)
🗄️ Repository (repository.go)
    ↓ (GORM operations)
📊 Model (models.go)
    ↓ (MySQL)
🗄️ Database
```

### Ejemplo Completo de Módulo

#### 1. Registro del Módulo (`init.go`)
```go
package nombre_modulo

import (
    "practicev2/registry"
    "github.com/gofiber/fiber/v2"
)

func init() {
    // Auto-registrar rutas
    registry.RegisterModule(RegisterRoutes)
    
    // Auto-registrar migraciones
    registry.RegisterMigration(100, MigrateTables)
}

func RegisterRoutes(app *fiber.App) {
    // Crear grupo de rutas
    api := app.Group("/api/entities")
    
    // Inicializar handler
    handler := NewHandler()
    
    // Definir rutas CRUD
    api.Get("/", handler.List)        // GET /api/entities
    api.Post("/", handler.Create)     // POST /api/entities
    api.Get("/:id", handler.Get)      // GET /api/entities/:id
    api.Put("/:id", handler.Update)   // PUT /api/entities/:id
    api.Delete("/:id", handler.Delete) // DELETE /api/entities/:id
    
    // Rutas específicas del dominio
    api.Get("/:id/details", handler.GetDetails)
    api.Post("/:id/activate", handler.Activate)
}
```

#### 2. Modelo GORM (`models.go`)
```go
package nombre_modulo

import (
    "time"
    "gorm.io/gorm"
)

// Entity representa la entidad principal del módulo
type Entity struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null;size:100;index"`
    Description string         `json:"description" gorm:"type:text"`
    Status      string         `json:"status" gorm:"default:'active';index"`
    Priority    int            `json:"priority" gorm:"default:0"`
    
    // Metadata
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relaciones
    UserID      string         `json:"user_id" gorm:"index"`
    User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Tags        []Tag          `json:"tags,omitempty" gorm:"many2many:entity_tags;"`
}

// TableName especifica el nombre de la tabla (evita pluralización)
func (Entity) TableName() string {
    return "entities"
}

// Hooks GORM
func (e *Entity) BeforeCreate(tx *gorm.DB) error {
    // Lógica antes de crear
    if e.Status == "" {
        e.Status = "active"
    }
    return nil
}

func (e *Entity) AfterCreate(tx *gorm.DB) error {
    // Log de auditoría, envío de eventos, etc.
    return nil
}

// Métodos de negocio
func (e *Entity) IsActive() bool {
    return e.Status == "active"
}

func (e *Entity) GetDisplayName() string {
    if e.Name != "" {
        return e.Name
    }
    return fmt.Sprintf("Entity #%d", e.ID)
}
```

#### 3. DTOs para API (`dto.go`)
```go
package nombre_modulo

// CreateEntityDTO para creación
type CreateEntityDTO struct {
    Name        string   `json:"name" validate:"required,min=2,max=100"`
    Description string   `json:"description" validate:"max=500"`
    Priority    int      `json:"priority" validate:"min=0,max=10"`
    TagIDs      []uint   `json:"tag_ids" validate:"dive,min=1"`
}

// UpdateEntityDTO para actualización
type UpdateEntityDTO struct {
    Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
    Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
    Priority    *int     `json:"priority,omitempty" validate:"omitempty,min=0,max=10"`
}

// EntityResponseDTO para respuestas
type EntityResponseDTO struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Priority    int       `json:"priority"`
    CreatedAt   string    `json:"created_at"`
    UpdatedAt   string    `json:"updated_at"`
    User        *UserDTO  `json:"user,omitempty"`
    Tags        []TagDTO  `json:"tags,omitempty"`
}

// ListEntitiesQuery para filtros y paginación
type ListEntitiesQuery struct {
    Page     int    `query:"page" validate:"min=1"`
    Limit    int    `query:"limit" validate:"min=1,max=100"`
    Status   string `query:"status" validate:"omitempty,oneof=active inactive"`
    Search   string `query:"search" validate:"max=100"`
    SortBy   string `query:"sort_by" validate:"omitempty,oneof=name created_at priority"`
    SortDir  string `query:"sort_dir" validate:"omitempty,oneof=asc desc"`
}
```

#### 4. Repository con GORM (`repository.go`)
```go
package nombre_modulo

import (
    "errors"
    "practicev2/database"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository() *Repository {
    return &Repository{db: database.DBconn}
}

// GetAll obtiene entidades con filtros y paginación
func (r *Repository) GetAll(query ListEntitiesQuery) ([]Entity, int64, error) {
    var entities []Entity
    var total int64
    
    // Construir query base
    baseQuery := r.db.Model(&Entity{})
    
    // Aplicar filtros
    if query.Status != "" {
        baseQuery = baseQuery.Where("status = ?", query.Status)
    }
    
    if query.Search != "" {
        baseQuery = baseQuery.Where("name ILIKE ? OR description ILIKE ?", 
            "%"+query.Search+"%", "%"+query.Search+"%")
    }
    
    // Contar total
    if err := baseQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Aplicar paginación y ordenamiento
    offset := (query.Page - 1) * query.Limit
    orderClause := "created_at DESC"
    if query.SortBy != "" {
        orderClause = query.SortBy
        if query.SortDir != "" {
            orderClause += " " + query.SortDir
        }
    }
    
    err := baseQuery.
        Preload("User").
        Preload("Tags").
        Order(orderClause).
        Offset(offset).
        Limit(query.Limit).
        Find(&entities).Error
    
    return entities, total, err
}

// GetByID obtiene una entidad por ID
func (r *Repository) GetByID(id uint) (*Entity, error) {
    var entity Entity
    err := r.db.
        Preload("User").
        Preload("Tags").
        First(&entity, id).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New("entidad no encontrada")
    }
    
    return &entity, err
}

// Create crea una nueva entidad
func (r *Repository) Create(entity *Entity) error {
    return r.db.Create(entity).Error
}

// Update actualiza una entidad
func (r *Repository) Update(entity *Entity) error {
    return r.db.Save(entity).Error
}

// Delete elimina una entidad (soft delete)
func (r *Repository) Delete(id uint) error {
    result := r.db.Delete(&Entity{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return errors.New("entidad no encontrada")
    }
    return nil
}

// Métodos específicos del dominio
func (r *Repository) GetByStatus(status string) ([]Entity, error) {
    var entities []Entity
    err := r.db.Where("status = ?", status).Find(&entities).Error
    return entities, err
}

func (r *Repository) UpdateStatus(id uint, status string) error {
    return r.db.Model(&Entity{}).Where("id = ?", id).Update("status", status).Error
}

// Transacciones complejas
func (r *Repository) CreateWithTags(entity *Entity, tagIDs []uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Crear entidad
        if err := tx.Create(entity).Error; err != nil {
            return err
        }
        
        // Asociar tags
        if len(tagIDs) > 0 {
            var tags []Tag
            if err := tx.Find(&tags, tagIDs).Error; err != nil {
                return err
            }
            if err := tx.Model(entity).Association("Tags").Append(tags); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## Guía de Implementación de Nuevos Módulos

### Paso 1: Crear la Estructura del Módulo

```bash
# Crear directorio del módulo
mkdir -p module/nombre_modulo
```

### Paso 2: Implementar el Archivo de Registro (`init.go`)

```go
// module/nombre_modulo/init.go
package nombre_modulo

import (
    "practicev2/registry"
    "github.com/gofiber/fiber/v2"
)

func init() {
    // Registrar las rutas del módulo
    registry.RegisterModule(registerRoutes)
    
    // Registrar migraciones si es necesario
    registry.RegisterMigration(100, migrateTables) // orden 100
}

func registerRoutes(app *fiber.App) {
    // Grupo de rutas para el módulo
    api := app.Group("/api/nombre_modulo")
    
    // Definir rutas
    api.Get("/", listHandler)
    api.Post("/", createHandler)
    api.Get("/:id", getHandler)
    api.Put("/:id", updateHandler)
    api.Delete("/:id", deleteHandler)
}
```

### Paso 3: Definir Modelos (`models.go`)

```go
// module/nombre_modulo/models.go
package nombre_modulo

import (
    "time"
    "gorm.io/gorm"
)

type Entity struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    Status    bool           `json:"status" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (Entity) TableName() string {
    return "entities"
}
```

### Paso 4: Implementar Handlers (`handlers.go`)

```go
// module/nombre_modulo/handlers.go
package nombre_modulo

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
)

func listHandler(c *fiber.Ctx) error {
    entities, err := GetAll()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error al obtener registros",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": entities,
    })
}

func createHandler(c *fiber.Ctx) error {
    var dto CreateEntityDTO
    
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Datos inválidos",
        })
    }
    
    if err := ValidateCreateEntity(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    entity, err := Create(dto)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error al crear registro",
        })
    }
    
    return c.Status(201).JSON(fiber.Map{
        "data": entity,
    })
}

func getHandler(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "ID inválido",
        })
    }
    
    entity, err := GetByID(uint(id))
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "Registro no encontrado",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": entity,
    })
}

func updateHandler(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "ID inválido",
        })
    }
    
    var dto UpdateEntityDTO
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Datos inválidos",
        })
    }
    
    if err := ValidateUpdateEntity(dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    entity, err := Update(uint(id), dto)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error al actualizar registro",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": entity,
    })
}

func deleteHandler(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "ID inválido",
        })
    }
    
    if err := Delete(uint(id)); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Error al eliminar registro",
        })
    }
    
    return c.JSON(fiber.Map{
        "message": "Registro eliminado correctamente",
    })
}
```

### Paso 5: Implementar Servicios (`services.go`)

```go
// module/nombre_modulo/services.go
package nombre_modulo

import (
    "errors"
    "practicev2/database"
)

func GetAll() ([]Entity, error) {
    var entities []Entity
    result := database.DBconn.Find(&entities)
    return entities, result.Error
}

func GetByID(id uint) (*Entity, error) {
    var entity Entity
    result := database.DBconn.First(&entity, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &entity, nil
}

func Create(dto CreateEntityDTO) (*Entity, error) {
    entity := Entity{
        Name:   dto.Name,
        Status: true,
    }
    
    result := database.DBconn.Create(&entity)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &entity, nil
}

func Update(id uint, dto UpdateEntityDTO) (*Entity, error) {
    entity, err := GetByID(id)
    if err != nil {
        return nil, err
    }
    
    entity.Name = dto.Name
    if dto.Status != nil {
        entity.Status = *dto.Status
    }
    
    result := database.DBconn.Save(entity)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return entity, nil
}

func Delete(id uint) error {
    result := database.DBconn.Delete(&Entity{}, id)
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return errors.New("registro no encontrado")
    }
    
    return nil
}
```

### Paso 6: Definir DTOs (`dto.go`)

```go
// module/nombre_modulo/dto.go
package nombre_modulo

type CreateEntityDTO struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateEntityDTO struct {
    Name   string `json:"name" validate:"required,min=2,max=100"`
    Status *bool  `json:"status,omitempty"`
}
```

### Paso 7: Implementar Validaciones (`validation.go`)

```go
// module/nombre_modulo/validation.go
package nombre_modulo

import (
    "errors"
    "strings"
)

func ValidateCreateEntity(dto CreateEntityDTO) error {
    if strings.TrimSpace(dto.Name) == "" {
        return errors.New("el nombre es requerido")
    }
    
    if len(dto.Name) < 2 || len(dto.Name) > 100 {
        return errors.New("el nombre debe tener entre 2 y 100 caracteres")
    }
    
    return nil
}

func ValidateUpdateEntity(dto UpdateEntityDTO) error {
    return ValidateCreateEntity(CreateEntityDTO{Name: dto.Name})
}
```

### Paso 8: Implementar Migraciones (`migrations.go`)

```go
// module/nombre_modulo/migrations.go
package nombre_modulo

import (
    "gorm.io/gorm"
)

func migrateTables(db *gorm.DB) error {
    return db.AutoMigrate(&Entity{})
}
```

### Paso 9: Registrar el Módulo

Agregar la importación en `addModules/modules.go`:

```go
// addModules/modules.go
package modules

import (
    // Módulos existentes
    _ "practicev2/module/authentication"
    
    // Nuevo módulo
    _ "practicev2/module/nombre_modulo"
)
```

## Flujo de Datos

### 1. Inicialización de la Aplicación

```
main.go
├── Cargar variables de entorno (.env)
├── config.Init() - Validar configuración
├── database.ConnectDatabase() - Conectar a MySQL
├── registry.RunMigrations() - Ejecutar migraciones
├── Configurar Fiber con middleware
├── registry.RegisterAllModules() - Registrar módulos
└── app.Listen() - Iniciar servidor
```

### 2. Procesamiento de Requests

```
HTTP Request
├── Middleware CORS
├── Middleware Logger
├── Middleware Personalizado
├── Router (Módulo específico)
├── Handler
├── Validación
├── Service Layer
├── Repository/Database
└── HTTP Response
```

### 3. Flujo de Registro de Módulos

```
Inicio de Aplicación
├── Importación de addModules/modules.go
├── Ejecución de init() de cada módulo
│   ├── registry.RegisterModule(routes)
│   └── registry.RegisterMigration(order, migration)
├── registry.RunMigrations() - Ordenar y ejecutar migraciones
└── registry.RegisterAllModules() - Registrar rutas
```

## Configuración y Variables de Entorno

### Variables Requeridas

Crear un archivo `.env` en la raíz del proyecto:

```bash
# Configuración de la aplicación
APP_HOST=localhost
APP_PORT=8080

# Base de datos MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=database_name

# JWT
JWT_KEY=tu-clave-jwt-muy-segura-de-al-menos-32-caracteres

# Redis (opcional)
REDIS_URL=redis://localhost:6379

# RabbitMQ (opcional)
RABBITMQ_URL=amqp://user:password@localhost:5672/
RABBITMQ_EXCHANGE=app_exchange

# Configuración de entorno
ENV=development
COOKIE_CONFIG=false  # true para producción
```

### Gestión de Configuración

El paquete `config` centraliza toda la configuración:

```go
// Ejemplo de uso en un módulo
import "practicev2/config"

// Obtener clave JWT
jwtKey := config.GetJWTKey()

// Verificar si estamos en producción
isProduction := config.IsProductionCookie()

// Obtener dominio de cookies
domain := config.GetCookieDomain()
```

## Base de Datos y Migraciones

### Configuración MySQL + GORM

#### Conexión Optimizada
```go
func ConnectDatabase() {
    // DSN para MySQL con configuraciones optimizadas
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        user, password, host, port, name)
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // No pluralizar nombres de tabla
        },
        // Configuraciones adicionales para producción
        PrepareStmt: true,              // Cache de prepared statements
        DisableForeignKeyConstraintWhenMigrating: false, // Mantener FK
    })
    
    // Configuración del pool de conexiones
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)           // Conexiones idle
    sqlDB.SetMaxOpenConns(100)          // Conexiones máximas
    sqlDB.SetConnMaxLifetime(time.Hour) // Tiempo de vida
}
```

#### Características GORM Utilizadas

- **🔗 Relaciones Complejas**: HasOne, HasMany, BelongsTo, Many2Many
- **🗑️ Soft Deletes**: Eliminación lógica con `DeletedAt`
- **🔄 Hooks**: Before/After Create, Update, Delete, Find
- **📊 Preloading**: Carga eager de relaciones
- **🔍 Scopes**: Queries reutilizables
- **🏷️ Tagging**: Sistema de etiquetas flexible
- **📈 Índices**: Optimización de consultas

### Sistema de Migraciones con Prioridades

```go
type Migration struct {
    Order int                    // Orden de ejecución
    Name  string                // Nombre descriptivo
    Func  func(*gorm.DB) error  // Función de migración
}

func RegisterMigration(order int, name string, migrationFunc func(*gorm.DB) error) {
    migrations = append(migrations, Migration{
        Order: order,
        Name:  name, 
        Func:  migrationFunc,
    })
}

func RunMigrations(db *gorm.DB) error {
    // Ordenar por prioridad
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Order < migrations[j].Order
    })
    
    // Ejecutar en orden
    for _, migration := range migrations {
        log.Printf("Ejecutando migración: %s (orden: %d)", migration.Name, migration.Order)
        if err := migration.Func(db); err != nil {
            return fmt.Errorf("error en migración %s: %w", migration.Name, err)
        }
    }
    
    return nil
}
```

#### Convención de Órdenes de Migración

| Orden | Propósito | Ejemplos |
|-------|-----------|----------|
| **1-99** | 🏗️ Infraestructura Base | Users, Roles, Organizations |
| **100-199** | 🏢 Módulos de Negocio | Products, Orders, Customers |
| **200-299** | 🔗 Relaciones y Referencias | Foreign Keys, Junction Tables |
| **300-399** | 📊 Índices y Optimizaciones | Database Indexes, Views |
| **400-499** | 📈 Datos de Configuración | Default Roles, System Settings |
| **500+** | 🧪 Datos de Prueba | Test Data, Demo Content |

### Ejemplo de Migración de Módulo

```go
// module/authentication/migrations.go
func init() {
    registry.RegisterMigration(10, "CreateIdentityTables", CreateIdentityTables)
    registry.RegisterMigration(11, "CreateRolesTables", CreateRolesTables)  
    registry.RegisterMigration(12, "CreateOrganizationTables", CreateOrganizationTables)
    registry.RegisterMigration(15, "CreateRefreshTokenTables", CreateRefreshTokenTables)
    registry.RegisterMigration(20, "CreateAuditLogTables", CreateAuditLogTables)
}

func CreateIdentityTables(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.Identity{},
        &models.UserProfile{},
        &models.CustomerProfile{},
    )
}

func CreateRolesTables(db *gorm.DB) error {
    if err := db.AutoMigrate(&models.Role{}, &models.Permission{}); err != nil {
        return err
    }
    
    // Crear índices personalizados
    return db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_roles_organization_name 
        ON roles (organization_id, name);
        
        CREATE INDEX IF NOT EXISTS idx_permissions_resource_action 
        ON permissions (resource, action);
    `).Error
}
```

### Mejores Prácticas para GORM

#### 1. Modelos con Relaciones
```go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex;not null"`
    
    // Relación One-to-One
    Profile  UserProfile `gorm:"foreignKey:UserID"`
    
    // Relación One-to-Many
    Orders   []Order `gorm:"foreignKey:UserID"`
    
    // Relación Many-to-Many
    Roles    []Role `gorm:"many2many:user_roles;"`
}

// Preload selectivo para optimización
func (r *Repository) GetUserWithDetails(id uint) (*User, error) {
    var user User
    return &user, r.db.
        Preload("Profile").
        Preload("Orders", "status = ?", "active").
        Preload("Roles.Permissions").
        First(&user, id).Error
}
```

#### 2. Scopes Reutilizables
```go
// Scopes para queries comunes
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

func WithProfile(db *gorm.DB) *gorm.DB {
    return db.Preload("Profile")
}

// Uso de scopes
users, err := r.db.Scopes(ActiveUsers, WithProfile).Find(&users).Error
```

#### 3. Transacciones Complejas
```go
func (r *Repository) CreateUserWithProfile(userData CreateUserDTO) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Crear usuario
        user := &User{Email: userData.Email}
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // Crear perfil
        profile := &UserProfile{
            UserID:    user.ID,
            FirstName: userData.FirstName,
            LastName:  userData.LastName,
        }
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        // Asignar rol por defecto
        var defaultRole Role
        if err := tx.Where("name = ?", "user").First(&defaultRole).Error; err != nil {
            return err
        }
        
        return tx.Model(user).Association("Roles").Append(&defaultRole)
    })
}
```

## Autenticación y Autorización

### Sistema JWT Avanzado

Implementación completa de autenticación con **Access + Refresh Tokens**, **roles**, **permisos** y **contextos organizacionales**.

#### Arquitectura de Tokens

```
┌─────────────────────────────────────────────────────────┐
│                    Token Strategy                       │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────▼─────────────┐
        │      Access Token         │
        │   • Corta duración (15m)  │
        │   • Contiene claims       │
        │   • Para autorización     │
        └─────────────┬─────────────┘
                      │
        ┌─────────────▼─────────────┐
        │     Refresh Token         │
        │   • Larga duración (7d)   │
        │   • Almacenado en BD      │
        │   • Para renovación       │
        └───────────────────────────┘
```

#### Estructura de Claims JWT

```go
type AccessTokenClaims struct {
    IdentityID     string `json:"identity_id"`
    Email          string `json:"email"`
    OrganizationID string `json:"organization_id,omitempty"`
    RoleID         string `json:"role_id,omitempty"`
    RoleName       string `json:"role_name,omitempty"`
    IsCustomer     bool   `json:"is_customer"`
    CustomerID     string `json:"customer_id,omitempty"`
    SessionID      string `json:"session_id"`
    
    jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
    IdentityID     string `json:"identity_id"`
    SessionID      string `json:"session_id"`
    OrganizationID string `json:"organization_id,omitempty"`
    RoleID         string `json:"role_id,omitempty"`
    RoleName       string `json:"role_name,omitempty"`
    CustomerID     string `json:"customer_id,omitempty"`
    
    jwt.RegisteredClaims
}
```

### Flujo de Autenticación

```mermaid
sequenceDiagram
    participant C as Client
    participant A as Auth API
    participant D as Database
    participant R as Redis Cache

    C->>A: POST /api/auth/login
    A->>D: Validate credentials
    D->>A: User data + memberships
    A->>A: Generate Access + Refresh tokens
    A->>D: Store refresh token
    A->>R: Cache user session
    A->>C: Return tokens + user data

    Note over C,A: Access token expires after 15 minutes

    C->>A: POST /api/auth/refresh
    A->>D: Validate refresh token
    A->>A: Generate new Access token
    A->>D: Update refresh token
    A->>C: Return new tokens

    Note over C,A: Refresh token rotates on each use
```

### Sistema de Roles y Permisos

#### Modelo de Datos
```go
type Role struct {
    ID             string       `json:"id" gorm:"primaryKey"`
    Name           string       `json:"name" gorm:"not null"`
    Description    string       `json:"description"`
    OrganizationID string       `json:"organization_id" gorm:"index"`
    IsSystemRole   bool         `json:"is_system_role" gorm:"default:false"`
    
    // Relaciones
    Permissions    []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
    Users          []User       `json:"users" gorm:"many2many:user_roles;"`
}

type Permission struct {
    ID          string `json:"id" gorm:"primaryKey"`
    Resource    string `json:"resource" gorm:"not null;index"`  // "users", "orders"
    Action      string `json:"action" gorm:"not null;index"`    // "create", "read", "update", "delete"
    Description string `json:"description"`
}

type OrganizationalMembership struct {
    ID             string    `json:"id" gorm:"primaryKey"`
    IdentityID     string    `json:"identity_id" gorm:"index"`
    OrganizationID string    `json:"organization_id" gorm:"index"`
    RoleID         string    `json:"role_id" gorm:"index"`
    
    // Metadatos específicos del contexto
    Department     string    `json:"department,omitempty"`
    Position       string    `json:"position,omitempty"`
    Grade          string    `json:"grade,omitempty"`          // Para contexto educativo
    Subject        string    `json:"subject,omitempty"`        // Para profesores
    StudentID      string    `json:"student_id,omitempty"`     // Para estudiantes
    EmployeeID     string    `json:"employee_id,omitempty"`    // Para empleados
    
    // Fechas de vigencia
    ActiveFrom     time.Time  `json:"active_from"`
    ActiveUntil    *time.Time `json:"active_until,omitempty"`
    IsActive       bool       `json:"is_active" gorm:"default:true;index"`
}
```

#### Middleware de Autorización

```go
// RequireAuth - Middleware que requiere autenticación
func (h *Handler) RequireAuth(c *fiber.Ctx) error {
    token, err := extractTokenFromHeader(c.Get("Authorization"))
    if err != nil {
        return sendUnauthorizedResponse(c, "Token requerido")
    }
    
    claims, err := validateAccessToken(token)
    if err != nil {
        return sendUnauthorizedResponse(c, "Token inválido")
    }
    
    // Establecer contexto de usuario
    c.Locals("identity_id", claims.IdentityID)
    c.Locals("email", claims.Email)
    c.Locals("organization_id", claims.OrganizationID)
    c.Locals("role_id", claims.RoleID)
    c.Locals("role_name", claims.RoleName)
    c.Locals("session_id", claims.SessionID)
    
    return c.Next()
}

// RequireRole - Middleware que requiere rol específico
func RequireRole(roleName string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role_name").(string)
        if userRole != roleName {
            return sendForbiddenResponse(c, "Rol insuficiente")
        }
        return c.Next()
    }
}

// RequirePermission - Middleware que requiere permiso específico
func RequirePermission(resource, action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        roleID := c.Locals("role_id").(string)
        
        hasPermission, err := checkUserPermission(roleID, resource, action)
        if err != nil || !hasPermission {
            return sendForbiddenResponse(c, "Permiso insuficiente")
        }
        
        return c.Next()
    }
}
```

#### Uso en Rutas

```go
func RegisterRoutes(app *fiber.App) {
    auth := NewAuthHandler()
    
    // Rutas públicas
    public := app.Group("/api/auth")
    public.Post("/login", auth.Login)
    public.Post("/register", auth.Register)
    public.Post("/refresh", auth.RefreshToken)
    
    // Rutas protegidas
    protected := app.Group("/api", auth.RequireAuth)
    protected.Get("/me", auth.Me)
    protected.Post("/logout", auth.Logout)
    
    // Rutas con roles específicos
    admin := protected.Group("/admin", RequireRole("admin"))
    admin.Get("/users", userHandler.ListAll)
    admin.Delete("/users/:id", userHandler.Delete)
    
    // Rutas con permisos específicos
    protected.Get("/orders", RequirePermission("orders", "read"), orderHandler.List)
    protected.Post("/orders", RequirePermission("orders", "create"), orderHandler.Create)
    protected.Put("/orders/:id", RequirePermission("orders", "update"), orderHandler.Update)
    protected.Delete("/orders/:id", RequirePermission("orders", "delete"), orderHandler.Delete)
}
```

### Multi-tenancy y Contextos Organizacionales

#### Cambio de Contexto
```go
// SwitchOrganizationContext cambia el contexto organizacional del usuario
func (s *Service) SwitchOrganizationContext(dto ContextSwitchDTO, identityID string) (*TokenPair, error) {
    // Verificar membresía activa
    membership, err := s.repo.GetActiveMembership(identityID, dto.OrganizationID)
    if err != nil {
        return nil, errors.New("sin acceso a esta organización")
    }
    
    // Obtener rol y permisos
    role, err := s.repo.GetRoleWithPermissions(membership.RoleID)
    if err != nil {
        return nil, errors.New("rol no encontrado")
    }
    
    // Generar nuevos tokens con nuevo contexto
    tokenPair, err := generateTokenPair(
        identityID,
        dto.OrganizationID,
        membership.RoleID,
        role.Name,
        uuid.New().String(), // nuevo session ID
    )
    
    return tokenPair, err
}
```

### Auditoría de Seguridad

```go
type AuditLog struct {
    ID             string                 `json:"id" gorm:"primaryKey"`
    IdentityID     string                 `json:"identity_id" gorm:"index"`
    OrganizationID string                 `json:"organization_id" gorm:"index"`
    Action         string                 `json:"action" gorm:"index"`        // "login", "logout", "access_denied"
    Resource       string                 `json:"resource,omitempty"`         // Recurso accedido
    IPAddress      string                 `json:"ip_address"`
    UserAgent      string                 `json:"user_agent"`
    Metadata       map[string]interface{} `json:"metadata" gorm:"type:json"`
    CreatedAt      time.Time              `json:"created_at"`
}

func LogSecurityEvent(identityID, organizationID, action, ipAddress, userAgent string, metadata map[string]interface{}) {
    auditLog := &AuditLog{
        ID:             uuid.New().String(),
        IdentityID:     identityID,
        OrganizationID: organizationID,
        Action:         action,
        IPAddress:      ipAddress,
        UserAgent:      userAgent,
        Metadata:       metadata,
        CreatedAt:      time.Now(),
    }
    
    database.DBconn.Create(auditLog)
}
```

## Configuración y Variables de Entorno

### Archivo de Configuración (.env)

```bash
# ===================================
# 🚀 CONFIGURACIÓN DE LA APLICACIÓN
# ===================================
APP_HOST=localhost
APP_PORT=8080
ENV=development  # development | production

# ===================================
# 🗄️ BASE DE DATOS MYSQL
# ===================================
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_secure_password
DB_NAME=practicev2_db

# ===================================
# 🔐 CONFIGURACIÓN JWT AVANZADA
# ===================================
JWT_KEY=your-main-jwt-key-must-be-at-least-32-characters-long
JWT_ACCESS_SECRET=your-access-token-secret-key-unique
JWT_REFRESH_SECRET=your-refresh-token-secret-key-different
JWT_ISSUER=practicev2-auth
JWT_ACCESS_TTL_MINUTES=15      # 15 minutos
JWT_REFRESH_TTL_HOURS=168      # 7 días

# ===================================
# 🍪 CONFIGURACIÓN DE COOKIES
# ===================================
COOKIE_CONFIG=false  # true para producción

# ===================================
# 🔴 REDIS CACHE (Opcional)
# ===================================
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# ===================================
# 🐰 RABBITMQ QUEUE (Opcional)
# ===================================
RABBITMQ_URL=amqp://user:password@localhost:5672/
RABBITMQ_EXCHANGE=practicev2_exchange
RABBITMQ_DEAD_LETTER_EXCHANGE=practicev2_dlx

# ===================================
# 📧 EMAIL CONFIGURATION (Opcional)
# ===================================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourapp.com

# ===================================
# 📁 FILE STORAGE (Opcional)
# ===================================
UPLOAD_MAX_SIZE=10485760  # 10MB en bytes
UPLOAD_ALLOWED_TYPES=jpg,jpeg,png,gif,pdf,doc,docx
STORAGE_PATH=./cmd/uploads
```

### Gestión de Configuración

#### Inicialización (`config/config.go`)
```go
type Config struct {
    // Server
    AppHost string
    AppPort string
    Env     string
    
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    
    // JWT
    JWTKey            string
    JWTAccessSecret   string
    JWTRefreshSecret  string
    JWTIssuer         string
    JWTAccessTTL      time.Duration
    JWTRefreshTTL     time.Duration
    
    // Redis
    RedisURL      string
    RedisPassword string
    RedisDB       int
    
    // RabbitMQ
    RabbitMQURL      string
    RabbitMQExchange string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    config := &Config{
        AppHost: getEnv("APP_HOST", "localhost"),
        AppPort: getEnv("APP_PORT", "8080"),
        Env:     getEnv("ENV", "development"),
        
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBUser:     mustGetEnv("DB_USER"),
        DBPassword: mustGetEnv("DB_PASSWORD"),
        DBName:     mustGetEnv("DB_NAME"),
        
        JWTKey:           mustGetEnv("JWT_KEY"),
        JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
        JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
        JWTIssuer:        getEnv("JWT_ISSUER", "practicev2"),
    }
    
    // Validar configuración crítica
    if err := config.Validate(); err != nil {
        log.Fatal("Configuration validation failed:", err)
    }
    
    return config
}

func (c *Config) Validate() error {
    if len(c.JWTKey) < 32 {
        return errors.New("JWT_KEY must be at least 32 characters")
    }
    
    if c.Env == "production" && c.AppHost == "localhost" {
        return errors.New("localhost not allowed in production")
    }
    
    return nil
}

func (c *Config) IsDevelopment() bool {
    return c.Env == "development"
}

func (c *Config) IsProduction() bool {
    return c.Env == "production"
}
```

### Configuración por Entorno

#### Desarrollo (development)
```go
if config.IsDevelopment() {
    // Configuración relajada para desarrollo
    corsConfig.AllowOrigins = []string{"*"}
    corsConfig.AllowCredentials = false
    
    // Logging detallado
    loggerConfig.Format = "[${time}] ${ip} ${status} - ${method} ${path} (${latency})\n"
    
    // Timeouts relajados
    serverConfig.ReadTimeout = 30 * time.Second
    serverConfig.WriteTimeout = 30 * time.Second
}
```

#### Producción (production)
```go
if config.IsProduction() {
    // Configuración estricta para producción
    corsConfig.AllowOrigins = []string{
        "https://yourdomain.com",
        "https://www.yourdomain.com",
    }
    corsConfig.AllowCredentials = true
    
    // Logging optimizado
    loggerConfig.Format = "${time} ${status} ${method} ${path}\n"
    
    // Timeouts optimizados
    serverConfig.ReadTimeout = 10 * time.Second
    serverConfig.WriteTimeout = 10 * time.Second
    
    // Configuración de seguridad
    app.Use(helmet.New())
    app.Use(limiter.New(limiter.Config{
        Max: 100,
        Expiration: 1 * time.Minute,
    }))
}
```

## Caché y Colas

### Sistema Redis Cache

Sistema de caché distribuido optimizado para múltiples conexiones y alta concurrencia:

#### Configuración Redis (`registry/cache/redis.go`)
```go
type RedisCache struct {
    client     *redis.Client
    ctx        context.Context
    prefix     string
    defaultTTL time.Duration
}

func NewRedisCacheFromEnv() (*RedisCache, error) {
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        return nil, errors.New("REDIS_URL not configured")
    }
    
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, fmt.Errorf("invalid Redis URL: %w", err)
    }
    
    // Configuración optimizada para producción
    opts.PoolSize = 20
    opts.MinIdleConns = 5
    opts.MaxRetries = 3
    opts.DialTimeout = 5 * time.Second
    opts.ReadTimeout = 3 * time.Second
    opts.WriteTimeout = 3 * time.Second
    
    client := redis.NewClient(opts)
    
    // Test de conexión
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("Redis connection failed: %w", err)
    }
    
    return &RedisCache{
        client:     client,
        ctx:        ctx,
        prefix:     "practicev2:",
        defaultTTL: 1 * time.Hour,
    }, nil
}

// Operaciones básicas
func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = r.defaultTTL
    }
    
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return r.client.Set(r.ctx, r.prefix+key, data, ttl).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) error {
    data, err := r.client.Get(r.ctx, r.prefix+key).Result()
    if err != nil {
        if err == redis.Nil {
            return errors.New("key not found")
        }
        return err
    }
    
    return json.Unmarshal([]byte(data), dest)
}

func (r *RedisCache) Delete(key string) error {
    return r.client.Del(r.ctx, r.prefix+key).Err()
}

// Operaciones avanzadas
func (r *RedisCache) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
    data, err := json.Marshal(value)
    if err != nil {
        return false, err
    }
    
    result := r.client.SetNX(r.ctx, r.prefix+key, data, ttl)
    return result.Val(), result.Err()
}

func (r *RedisCache) Increment(key string) (int64, error) {
    return r.client.Incr(r.ctx, r.prefix+key).Result()
}

func (r *RedisCache) SetWithExpire(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = r.defaultTTL
    }
    return r.Set(key, value, ttl)
}

// Cache Pattern: GetOrSet
func (r *RedisCache) GetOrSet(key string, dest interface{}, fetchFunc func() (interface{}, error), ttl time.Duration) error {
    // Intentar obtener del cache
    err := r.Get(key, dest)
    if err == nil {
        return nil // Cache hit
    }
    
    // Cache miss - obtener datos
    data, err := fetchFunc()
    if err != nil {
        return err
    }
    
    // Almacenar en cache
    if err := r.Set(key, data, ttl); err != nil {
        log.Printf("Failed to cache data for key %s: %v", key, err)
    }
    
    // Copiar datos al destino
    jsonData, _ := json.Marshal(data)
    return json.Unmarshal(jsonData, dest)
}
```

### Sistema RabbitMQ Queue

Sistema de colas para procesamiento asíncrono con reconexión automática:

#### Configuración RabbitMQ (`registry/queue/rabbitmq.go`)
```go
type RabbitMQ struct {
    conn         *amqp.Connection
    channel      *amqp.Channel
    url          string
    exchange     string
    queues       map[string]bool
    reconnecting bool
    mu           sync.RWMutex
}

func NewRabbitMQFromEnv() (*RabbitMQ, error) {
    url := os.Getenv("RABBITMQ_URL")
    if url == "" {
        return nil, errors.New("RABBITMQ_URL not configured")
    }
    
    exchange := os.Getenv("RABBITMQ_EXCHANGE")
    if exchange == "" {
        exchange = "practicev2_exchange"
    }
    
    rmq := &RabbitMQ{
        url:      url,
        exchange: exchange,
        queues:   make(map[string]bool),
    }
    
    if err := rmq.connect(); err != nil {
        return nil, err
    }
    
    return rmq, nil
}

func (r *RabbitMQ) connect() error {
    conn, err := amqp.Dial(r.url)
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }
    
    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return fmt.Errorf("failed to open channel: %w", err)
    }
    
    // Declarar exchange
    err = ch.ExchangeDeclare(
        r.exchange, // name
        "topic",    // type
        true,       // durable
        false,      // auto-deleted
        false,      // internal
        false,      // no-wait
        nil,        // arguments
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return fmt.Errorf("failed to declare exchange: %w", err)
    }
    
    r.conn = conn
    r.channel = ch
    
    // Configurar reconexión automática
    go r.watchConnection()
    
    return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, message interface{}) error {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    if r.channel == nil {
        return errors.New("RabbitMQ channel not available")
    }
    
    body, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }
    
    return r.channel.Publish(
        r.exchange,  // exchange
        routingKey,  // routing key
        false,       // mandatory
        false,       // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            DeliveryMode: amqp.Persistent, // Persistir mensajes
            Timestamp:    time.Now(),
        },
    )
}

func (r *RabbitMQ) Subscribe(queueName, routingKey string, handler func([]byte) error) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Declarar cola
    queue, err := r.channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %w", err)
    }
    
    // Bind cola al exchange
    err = r.channel.QueueBind(
        queue.Name,  // queue name
        routingKey,  // routing key
        r.exchange,  // exchange
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue: %w", err)
    }
    
    // Configurar QoS
    err = r.channel.Qos(
        1,     // prefetch count
        0,     // prefetch size
        false, // global
    )
    if err != nil {
        return fmt.Errorf("failed to set QoS: %w", err)
    }
    
    // Consumir mensajes
    msgs, err := r.channel.Consume(
        queue.Name, // queue
        "",         // consumer
        false,      // auto-ack
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %w", err)
    }
    
    r.queues[queueName] = true
    
    // Procesar mensajes en goroutine
    go func() {
        for msg := range msgs {
            if err := handler(msg.Body); err != nil {
                log.Printf("Error processing message: %v", err)
                msg.Nack(false, true) // Requeue en caso de error
            } else {
                msg.Ack(false)
            }
        }
    }()
    
    log.Printf("Subscribed to queue: %s with routing key: %s", queueName, routingKey)
    return nil
}

// Reconexión automática
func (r *RabbitMQ) watchConnection() {
    for {
        reason, ok := <-r.conn.NotifyClose(make(chan *amqp.Error))
        if !ok {
            log.Println("RabbitMQ connection closed")
            break
        }
        
        log.Printf("RabbitMQ connection lost: %v", reason)
        r.reconnect()
    }
}

func (r *RabbitMQ) reconnect() {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if r.reconnecting {
        return
    }
    
    r.reconnecting = true
    defer func() { r.reconnecting = false }()
    
    for {
        log.Println("Attempting to reconnect to RabbitMQ...")
        
        if err := r.connect(); err != nil {
            log.Printf("Reconnection failed: %v. Retrying in 5 seconds...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        
        log.Println("Successfully reconnected to RabbitMQ")
        break
    }
}
```

### Uso en Módulos

#### Cache en Services
```go
func (s *Service) GetUserByID(id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    var user User
    err := s.cache.GetOrSet(cacheKey, &user, func() (interface{}, error) {
        return s.repo.GetByID(id)
    }, 15*time.Minute)
    
    return &user, err
}

func (s *Service) InvalidateUserCache(id uint) {
    cacheKey := fmt.Sprintf("user:%d", id)
    s.cache.Delete(cacheKey)
}
```

#### Queue en Background Jobs
```go
// Publicar evento después de crear usuario
func (s *Service) CreateUser(dto CreateUserDTO) (*User, error) {
    user, err := s.repo.Create(dto)
    if err != nil {
        return nil, err
    }
    
    // Enviar email de bienvenida de forma asíncrona
    welcomeEvent := WelcomeEmailEvent{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.FirstName,
    }
    
    if err := s.queue.Publish(context.Background(), "email.welcome", welcomeEvent); err != nil {
        log.Printf("Failed to publish welcome email event: %v", err)
        // No fallar la operación principal por esto
    }
    
    return user, nil
}

// Suscribirse a eventos de email
func (s *EmailService) StartEmailWorker() error {
    return s.queue.Subscribe("email_queue", "email.*", s.handleEmailEvent)
}

func (s *EmailService) handleEmailEvent(body []byte) error {
    var event WelcomeEmailEvent
    if err := json.Unmarshal(body, &event); err != nil {
        return err
    }
    
    return s.SendWelcomeEmail(event.Email, event.Name)
}
```

## Caché y Colas

### Redis Cache

Sistema de caché distribuido optimizado para múltiples conexiones:

```go
// Usar caché en un módulo
import "practicev2/registry/cache"

// Crear cliente Redis
redisCache, err := cache.NewRedisCacheFromEnv()
if err != nil {
    // Manejar error
}

// Guardar en caché
err = redisCache.Set(ctx, "key", data, 1*time.Hour)

// Obtener de caché
var result MyStruct
err = redisCache.Get(ctx, "key", &result)
```

### RabbitMQ Queue

Sistema de colas para procesamiento asíncrono:

```go
// Usar colas en un módulo
import "practicev2/registry/queue"

// Crear cliente RabbitMQ
rabbitmq, err := queue.NewRabbitMQFromEnv()
if err != nil {
    // Manejar error
}

// Publicar mensaje
err = rabbitmq.Publish(ctx, "exchange", "routing.key", message)

// Suscribirse a cola
err = rabbitmq.Subscribe("queue_name", "routing.key", handleMessage)
```

## Mejores Prácticas

### 1. 🏗️ Organización de Código y Arquitectura

#### Estructura de Módulos Consistente
```go
// Cada módulo debe seguir esta estructura exacta
module/nombre_modulo/
├── init.go              # Registro automático del módulo
├── routes.go            # Definición de rutas Fiber
├── handler.go           # Controllers HTTP/REST
├── service.go           # Lógica de negocio (Use Cases)
├── repository.go        # Acceso a datos con GORM
├── models.go            # Entidades/Modelos GORM
├── dto.go               # Data Transfer Objects
├── validation.go        # Validaciones de negocio
├── migrations.go        # Migraciones específicas
└── [submódulos]/        # Organización por subdominios
```

#### Principios SOLID en Go
```go
// 1. Single Responsibility Principle
type UserService struct {
    repo UserRepository
    auth AuthService
    mail EmailService
}

// 2. Open/Closed Principle - Uso de interfaces
type PaymentProcessor interface {
    ProcessPayment(amount decimal.Decimal) error
}

type CreditCardProcessor struct{}
func (c *CreditCardProcessor) ProcessPayment(amount decimal.Decimal) error { /* */ }

type PayPalProcessor struct{}
func (p *PayPalProcessor) ProcessPayment(amount decimal.Decimal) error { /* */ }

// 3. Liskov Substitution Principle
func ProcessOrder(processor PaymentProcessor, amount decimal.Decimal) error {
    return processor.ProcessPayment(amount) // Cualquier implementación funciona
}

// 4. Interface Segregation Principle
type UserReader interface {
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
}

type UserWriter interface {
    Create(user *User) error
    Update(user *User) error
    Delete(id uint) error
}

// 5. Dependency Inversion Principle
type UserService struct {
    reader UserReader
    writer UserWriter
    cache  CacheInterface
}
```

### 2. 🗄️ Mejores Prácticas de Base de Datos con GORM

#### Optimización de Queries
```go
// ✅ CORRECTO: Preload selectivo
func (r *Repository) GetUserWithActiveOrders(id uint) (*User, error) {
    var user User
    return &user, r.db.
        Preload("Orders", "status = ? AND created_at > ?", "active", time.Now().AddDate(0, -1, 0)).
        Preload("Profile").
        First(&user, id).Error
}

// ❌ INCORRECTO: Preload innecesario
func (r *Repository) GetUser(id uint) (*User, error) {
    var user User
    return &user, r.db.
        Preload("Orders").           // Carga TODAS las órdenes
        Preload("Orders.Items").     // N+1 queries
        Preload("Profile").
        First(&user, id).Error
}

// ✅ CORRECTO: Paginación eficiente
func (r *Repository) GetUsersWithPagination(page, limit int) ([]User, int64, error) {
    var users []User
    var total int64
    
    // Contar total primero
    if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Query paginado
    offset := (page - 1) * limit
    err := r.db.
        Offset(offset).
        Limit(limit).
        Order("created_at DESC").
        Find(&users).Error
        
    return users, total, err
}

// ✅ CORRECTO: Uso de índices
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"uniqueIndex;not null"`                // Índice único
    Status   string `gorm:"index"`                               // Índice simple
    Name     string `gorm:"index:idx_name_status,priority:1"`    // Índice compuesto
    Status2  string `gorm:"index:idx_name_status,priority:2"`    // Índice compuesto
}
```

#### Transacciones Seguras
```go
// ✅ CORRECTO: Manejo completo de transacciones
func (r *Repository) CreateUserWithProfile(userData CreateUserDTO) (*User, error) {
    var user *User
    
    err := r.db.Transaction(func(tx *gorm.DB) error {
        // Crear usuario
        user = &User{
            Email:     userData.Email,
            FirstName: userData.FirstName,
            LastName:  userData.LastName,
        }
        
        if err := tx.Create(user).Error; err != nil {
            return fmt.Errorf("failed to create user: %w", err)
        }
        
        // Crear perfil
        profile := &UserProfile{
            UserID: user.ID,
            Bio:    userData.Bio,
        }
        
        if err := tx.Create(profile).Error; err != nil {
            return fmt.Errorf("failed to create profile: %w", err)
        }
        
        // Asignar rol por defecto
        defaultRole := &Role{}
        if err := tx.Where("name = ? AND is_default = ?", "user", true).First(defaultRole).Error; err != nil {
            return fmt.Errorf("default role not found: %w", err)
        }
        
        if err := tx.Model(user).Association("Roles").Append(defaultRole); err != nil {
            return fmt.Errorf("failed to assign role: %w", err)
        }
        
        return nil
    })
    
    return user, err
}
```

### 3. 🚀 Optimización de Performance con Fiber

#### Middleware Optimizado
```go
// ✅ CORRECTO: Middleware con configuración por entorno
func ConfigureMiddleware(app *fiber.App, config *Config) {
    // Compresión para respuestas grandes
    app.Use(compress.New(compress.Config{
        Level: compress.LevelBestSpeed, // Balance velocidad/tamaño
    }))
    
    // Rate limiting inteligente
    app.Use(limiter.New(limiter.Config{
        Max:        100,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            // Rate limiting por IP + User ID si está autenticado
            if userID := c.Locals("identity_id"); userID != nil {
                return fmt.Sprintf("%s:%s", c.IP(), userID)
            }
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many requests",
                "retry_after": 60,
            })
        },
    }))
    
    // CORS optimizado por entorno
    if config.IsDevelopment() {
        app.Use(cors.New(cors.Config{
            AllowOrigins:     "*",
            AllowCredentials: false,
        }))
    } else {
        app.Use(cors.New(cors.Config{
            AllowOrigins:     "https://yourdomain.com,https://www.yourdomain.com",
            AllowCredentials: true,
            MaxAge:           3600,
        }))
    }
}

// ✅ CORRECTO: Headers de seguridad en producción
if config.IsProduction() {
    app.Use(helmet.New(helmet.Config{
        XSSProtection:      "1; mode=block",
        ContentTypeNosniff: "nosniff",
        XFrameOptions:      "DENY",
        HSTSMaxAge:         31536000,
    }))
}
```

#### Manejo de Errores Estructurado
```go
// ✅ CORRECTO: Error handler global
func SetupErrorHandler(app *fiber.App) {
    app.Use(func(c *fiber.Ctx) error {
        err := c.Next()
        
        if err != nil {
            // Log del error con contexto
            log.Printf("[ERROR] %s %s - %v", c.Method(), c.Path(), err)
            
            // Determinar código de estado
            code := fiber.StatusInternalServerError
            message := "Error interno del servidor"
            
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
                message = e.Message
            }
            
            // Respuesta estructurada
            return c.Status(code).JSON(ErrorResponse{
                Error:     true,
                Message:   message,
                Code:      code,
                Timestamp: time.Now().Unix(),
                Path:      c.Path(),
                RequestID: c.Get("X-Request-ID"),
            })
        }
        
        return nil
    })
}

type ErrorResponse struct {
    Error     bool   `json:"error"`
    Message   string `json:"message"`
    Code      int    `json:"code"`
    Timestamp int64  `json:"timestamp"`
    Path      string `json:"path"`
    RequestID string `json:"request_id,omitempty"`
}
```

### 4. 🔐 Seguridad y Validación

#### Validación Robusta de DTOs
```go
// ✅ CORRECTO: Validación estructurada con mensajes claros
type CreateUserDTO struct {
    Email           string `json:"email" validate:"required,email,max=255"`
    Password        string `json:"password" validate:"required,min=8,max=72"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
    FirstName       string `json:"first_name" validate:"required,min=2,max=50,alpha"`
    LastName        string `json:"last_name" validate:"required,min=2,max=50,alpha"`
    Phone           string `json:"phone" validate:"omitempty,e164"` // Formato internacional
    DateOfBirth     string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
}

func (dto *CreateUserDTO) Validate() error {
    validator := validator.New()
    
    // Validaciones personalizadas
    validator.RegisterValidation("strong_password", validateStrongPassword)
    
    if err := validator.Struct(dto); err != nil {
        return formatValidationErrors(err)
    }
    
    // Validaciones de negocio
    if dto.Password != dto.ConfirmPassword {
        return errors.New("las contraseñas no coinciden")
    }
    
    return nil
}

func validateStrongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    // Al menos 1 mayúscula, 1 minúscula, 1 número, 1 símbolo
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber && hasSymbol
}
```

#### Sanitización de Datos
```go
// ✅ CORRECTO: Sanitización antes de procesar
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    var dto CreateUserDTO
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
    }
    
    // Sanitizar datos de entrada
    dto.Email = strings.ToLower(strings.TrimSpace(dto.Email))
    dto.FirstName = strings.TrimSpace(dto.FirstName)
    dto.LastName = strings.TrimSpace(dto.LastName)
    dto.Phone = strings.ReplaceAll(dto.Phone, " ", "")
    
    // Validar después de sanitizar
    if err := dto.Validate(); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    
    // Procesar...
}
```

### 5. 📊 Monitoring y Logging

#### Logging Estructurado
```go
// ✅ CORRECTO: Logging con niveles y contexto
type Logger struct {
    logger *log.Logger
    level  LogLevel
}

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

func (l *Logger) Info(message string, fields map[string]interface{}) {
    if l.level <= INFO {
        entry := LogEntry{
            Level:     "INFO",
            Message:   message,
            Fields:    fields,
            Timestamp: time.Now(),
        }
        l.logger.Println(entry.Format())
    }
}

func (l *Logger) Error(message string, err error, fields map[string]interface{}) {
    if l.level <= ERROR {
        if fields == nil {
            fields = make(map[string]interface{})
        }
        fields["error"] = err.Error()
        
        entry := LogEntry{
            Level:     "ERROR",
            Message:   message,
            Fields:    fields,
            Timestamp: time.Now(),
        }
        l.logger.Println(entry.Format())
    }
}

// Uso en handlers
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    logger := h.logger.WithFields(map[string]interface{}{
        "handler": "CreateUser",
        "ip":      c.IP(),
        "user_agent": c.Get("User-Agent"),
    })
    
    logger.Info("Creating new user", nil)
    
    // ... procesamiento ...
    
    if err != nil {
        logger.Error("Failed to create user", err, map[string]interface{}{
            "email": dto.Email,
        })
        return c.Status(500).JSON(fiber.Map{"error": "Error interno"})
    }
    
    logger.Info("User created successfully", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })
    
    return c.JSON(user)
}
```

### 6. 🧪 Testing Estratégico

#### Testing de Múltiples Capas
```go
// ✅ CORRECTO: Tests de Repository
func TestUserRepository_Create(t *testing.T) {
    // Setup de base de datos de prueba
    db := setupTestDB(t)
    repo := NewUserRepository(db)
    
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name: "valid user",
            user: &User{
                Email:     "test@example.com",
                FirstName: "John",
                LastName:  "Doe",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            user: &User{
                Email: "duplicate@example.com",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := repo.Create(tt.user)
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// ✅ CORRECTO: Tests de Service con mocks
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mockAuth := &MockAuthService{}
    service := NewUserService(mockRepo, mockAuth)
    
    dto := CreateUserDTO{
        Email:     "test@example.com",
        Password:  "SecurePass123!",
        FirstName: "John",
        LastName:  "Doe",
    }
    
    // Configurar mocks
    mockRepo.On("EmailExists", dto.Email).Return(false, nil)
    mockRepo.On("Create", mock.AnythingOfType("*User")).Return(nil)
    mockAuth.On("HashPassword", dto.Password).Return("hashed_password", nil)
    
    user, err := service.CreateUser(dto)
    
    assert.NoError(t, err)
    assert.Equal(t, dto.Email, user.Email)
    mockRepo.AssertExpectations(t)
    mockAuth.AssertExpectations(t)
}

// ✅ CORRECTO: Tests de integración con Fiber
func TestUserHandler_CreateUser(t *testing.T) {
    app := setupTestApp(t)
    
    payload := `{
        "email": "test@example.com",
        "password": "SecurePass123!",
        "confirm_password": "SecurePass123!",
        "first_name": "John",
        "last_name": "Doe"
    }`
    
    req := httptest.NewRequest("POST", "/api/users", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)
    
    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "test@example.com", response["email"])
}
```

### 7. 🚀 Deploy y DevOps

#### Docker para Desarrollo
```dockerfile
# Dockerfile para desarrollo
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/cmd/uploads ./cmd/uploads

EXPOSE 8080
CMD ["./main"]
```

#### Docker Compose para Stack Completo
```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=development
      - DB_HOST=mysql
      - REDIS_URL=redis://redis:6379
      - RABBITMQ_URL=amqp://rabbitmq:5672
    depends_on:
      - mysql
      - redis
      - rabbitmq
    volumes:
      - ./cmd/uploads:/root/cmd/uploads

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: practicev2_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: admin
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  mysql_data:
  redis_data:
  rabbitmq_data:
```

---

## 🎯 Conclusión

Esta arquitectura actualizada con **MySQL + GORM + Fiber** proporciona una base sólida y moderna para desarrollar aplicaciones web escalables en Go. Las principales ventajas incluyen:

### ✅ Beneficios Clave

1. **🚀 Alto Rendimiento**: Fiber ofrece rendimiento excepcional con su arquitectura basada en FastHTTP
2. **🗄️ ORM Moderno**: GORM proporciona un ORM maduro con soporte completo para relaciones complejas
3. **🔐 Seguridad Robusta**: Sistema de autenticación JWT avanzado con roles y permisos granulares
4. **📦 Modularidad**: Arquitectura de módulos auto-registrables que facilita el desarrollo en equipo
5. **⚡ Escalabilidad**: Soporte nativo para cache Redis y colas RabbitMQ
6. **🛠️ DevOps Ready**: Configuración completa para Docker y despliegue en producción

### 🔑 Factores de Éxito

- **Consistencia**: Seguir estrictamente las convenciones establecidas
- **Documentación**: Mantener documentación actualizada para cada módulo
- **Testing**: Implementar tests en todas las capas (unit, integration, e2e)
- **Monitoring**: Configurar logging y métricas desde el inicio
- **Security**: Aplicar principios de seguridad en cada capa

La clave del éxito está en mantener la disciplina arquitectónica y seguir las mejores prácticas establecidas. Cada nuevo módulo debe adherirse a los patrones definidos para garantizar la coherencia y mantenibilidad del sistema a largo plazo.

### 📚 Próximos Pasos

1. **Implementar Módulos Específicos**: User Management, Product Catalog, Order Management
2. **Configurar CI/CD**: Pipeline automatizado para testing y despliegue
3. **Métricas y Observabilidad**: Prometheus, Grafana, distributed tracing
4. **Optimización**: Profiling de performance y optimización de queries
5. **Documentación API**: Swagger/OpenAPI para documentación automática
