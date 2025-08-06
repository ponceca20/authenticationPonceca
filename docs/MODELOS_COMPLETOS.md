# DOCUMENTACIÓN COMPLETA DE MODELOS - SISTEMA DE AUTENTICACIÓN

## 📋 ÍNDICE DE CONTENIDO

1. [Resumen General](#resumen-general)
2. [Modelos Core](#modelos-core)
3. [Modelos de Perfil](#modelos-de-perfil)
4. [Modelos Organizacionales](#modelos-organizacionales)
5. [Modelos de Seguridad](#modelos-de-seguridad)
6. [Modelos de Invitación](#modelos-de-invitación)
7. [Diagrama de Relaciones](#diagrama-de-relaciones)
8. [Casos de Uso](#casos-de-uso)

---

## 🎯 RESUMEN GENERAL

Este sistema de autenticación está diseñado con una arquitectura multi-tenant que soporta:
- **Identidades unificadas**: Un usuario puede pertenecer a múltiples organizaciones
- **E-commerce**: Perfiles de cliente con características comerciales
- **Instituciones educativas**: Gestión de estudiantes, profesores y empleados
- **Empresas**: Gestión de empleados y departamentos
- **Auditoría completa**: Registro de todas las acciones del sistema
- **Sesiones de invitados**: Para usuarios no registrados

---

## 🔑 MODELOS CORE

### 1. Identity (User)
**Propósito**: Identidad central y única de una persona en el sistema.

**Campos principales**:
- `ID`: Identificador único (UUID)
- `Email`: Email único del usuario
- `PasswordHash`: Hash de la contraseña
- `FirstName`, `LastName`: Nombres
- `EmailVerified`: Estado de verificación
- `LastLoginAt`: Último login
- `FailedLoginAttempts`: Intentos fallidos de login
- `LockedUntil`: Fecha de bloqueo

**Uso**:
```go
// User es un alias de Identity para compatibilidad
type User = Identity
```

### 2. Organization
**Propósito**: Representa organizaciones multi-tenant (empresas, escuelas, etc.).

**Campos principales**:
- `ID`: Identificador único
- `Name`: Nombre de la organización
- `Slug`: Identificador único en URL
- `Type`: Tipo (company, educational_institution)
- `IsActive`: Estado activo/inactivo

### 3. OrganizationalMembership
**Propósito**: Vincula una Identity con una Organization mediante un Role específico.

**Campos principales**:
- `IdentityID`: Referencia al usuario
- `OrganizationID`: Referencia a la organización
- `RoleID`: Referencia al rol
- `Department`: Departamento (para empleados)
- `Grade`: Grado (para estudiantes)
- `StudentID`, `EmployeeID`: IDs específicos
- `ActiveFrom`, `ActiveUntil`: Periodo de actividad

---

## 👤 MODELOS DE PERFIL

### 4. CustomerProfile
**Propósito**: Información específica para clientes de e-commerce.

**Campos principales**:
- `CustomerNumber`: Número único de cliente
- `PreferredPaymentMethod`: Método de pago preferido
- `CreditLimit`: Límite de crédito
- `TotalSpent`: Total gastado
- `LoyaltyPoints`: Puntos de lealtad
- `AcceptsMarketing`: Acepta marketing

### 5. CustomerPreferences
**Propósito**: Preferencias específicas del cliente.

**Campos principales**:
- `Theme`: Tema de la interfaz
- `Language`: Idioma preferido
- `TimeZone`: Zona horaria
- `EmailNotifications`: Notificaciones por email
- `SmsNotifications`: Notificaciones por SMS

### 6. UserProfile
**Propósito**: Información extendida no esencial del usuario.

**Campos principales**:
- `Bio`: Biografía
- `Location`: Ubicación
- `Website`: Sitio web personal
- `Socials`: Enlaces sociales (JSON)

### 7. ShippingAddress
**Propósito**: Direcciones de envío para clientes.

**Campos principales**:
- `AddressLine1`, `AddressLine2`: Líneas de dirección
- `City`, `State`, `PostalCode`, `Country`: Datos de ubicación
- `IsDefault`: Si es la dirección por defecto

---

## 🏢 MODELOS ORGANIZACIONALES

### 8. Department
**Propósito**: Subdivisiones dentro de organizaciones.

**Campos principales**:
- `OrganizationID`: Organización padre
- `ParentID`: Departamento padre (estructura jerárquica)
- `Name`: Nombre del departamento
- `Description`: Descripción

### 9. Role
**Propósito**: Roles específicos dentro de organizaciones.

**Campos principales**:
- `OrganizationID`: Organización
- `Name`: Nombre único dentro de la organización
- `DisplayName`: Nombre para mostrar
- `HierarchyLevel`: Nivel jerárquico
- `IsSystemRole`: Si es un rol del sistema
- `Permissions`: Lista de permisos

### 10. Permission
**Propósito**: Permisos específicos sobre recursos.

**Campos principales**:
- `Resource`: Recurso (users, products, etc.)
- `Action`: Acción (create, read, update, delete)
- `Scope`: Alcance (own, department, organization, all)

---

## 🔒 MODELOS DE SEGURIDAD

### 11. AuditLog
**Propósito**: Registro de auditoría de todas las acciones.

**Campos principales**:
- `IdentityID`: Usuario que realizó la acción
- `OrganizationID`: Organización (opcional)
- `Action`: Acción realizada
- `Resource`: Recurso afectado
- `ResourceID`: ID del recurso
- `Status`: Estado (success, failure)
- `IPAddress`, `UserAgent`: Información de la sesión
- `Details`: Detalles adicionales (JSON)

### 12. RefreshToken
**Propósito**: Tokens para refrescar autenticación.

**Campos principales**:
- `IdentityID`: Usuario propietario
- `Token`: Token encriptado
- `ExpiresAt`: Fecha de expiración
- `IsRevoked`: Si está revocado

### 13. PasswordResetToken
**Propósito**: Tokens para reseteo de contraseñas.

**Campos principales**:
- `IdentityID`: Usuario
- `Token`: Token de reseteo
- `ExpiresAt`: Fecha de expiración

---

## 📨 MODELOS DE INVITACIÓN

### 14. Invitation
**Propósito**: Invitaciones para unirse a organizaciones.

**Campos principales**:
- `OrganizationID`: Organización de destino
- `InviterID`: Usuario que invita
- `Email`: Email del invitado
- `RoleID`: Rol asignado
- `Token`: Token de invitación
- `Status`: Estado (pending, accepted, expired)
- `ExpiresAt`: Fecha de expiración

### 15. GuestSession
**Propósito**: Sesiones temporales para usuarios no registrados.

**Campos principales**:
- `SessionToken`: Token de sesión
- `Email`, `FirstName`, `LastName`: Datos temporales
- `CartData`: Datos del carrito (JSON)
- `LastActivity`: Última actividad
- `ExpiresAt`: Fecha de expiración

---

## 🔄 DIAGRAMA DE RELACIONES

```
Identity (User)
├── 1:0..1 ── CustomerProfile
│             ├── 1:0..1 ── CustomerPreferences
│             └── 1:0..* ── ShippingAddress
├── 1:0..1 ── UserProfile
├── 1:0..* ── OrganizationalMembership
│             ├── *:1 ── Organization
│             │         ├── 1:0..* ── Department
│             │         ├── 1:0..* ── Role
│             │         └── 1:0..* ── Invitation
│             └── *:1 ── Role
│                       └── *:* ── Permission
├── 1:0..* ── RefreshToken
├── 1:0..* ── PasswordResetToken
└── 1:0..* ── AuditLog

GuestSession (independiente)
```

---

## 📝 CASOS DE USO

### Caso 1: Usuario E-commerce
```go
// Un usuario puede tener perfil de cliente
user := Identity{...}
customerProfile := CustomerProfile{IdentityID: user.ID, ...}
preferences := CustomerPreferences{CustomerProfileID: customerProfile.ID, ...}
address := ShippingAddress{CustomerProfileID: customerProfile.ID, ...}
```

### Caso 2: Empleado de Empresa
```go
// Un usuario puede ser empleado de una organización
user := Identity{...}
company := Organization{Type: "company", ...}
employeeRole := Role{OrganizationID: company.ID, Name: "employee", ...}
membership := OrganizationalMembership{
    IdentityID: user.ID,
    OrganizationID: company.ID,
    RoleID: employeeRole.ID,
    EmployeeID: "EMP001",
}
```

### Caso 3: Estudiante en Escuela
```go
// Un usuario puede ser estudiante
user := Identity{...}
school := Organization{Type: "educational_institution", ...}
studentRole := Role{OrganizationID: school.ID, Name: "student", ...}
membership := OrganizationalMembership{
    IdentityID: user.ID,
    OrganizationID: school.ID,
    RoleID: studentRole.ID,
    Grade: "10th",
    StudentID: "STU2025001",
}
```

### Caso 4: Usuario Multi-organización
```go
// Un usuario puede pertenecer a múltiples organizaciones
user := Identity{...}

// Como empleado en empresa
companyMembership := OrganizationalMembership{
    IdentityID: user.ID,
    OrganizationID: company.ID,
    RoleID: employeeRole.ID,
}

// Como cliente en tienda
customerProfile := CustomerProfile{IdentityID: user.ID, ...}

// Como estudiante en curso online
schoolMembership := OrganizationalMembership{
    IdentityID: user.ID,
    OrganizationID: onlineSchool.ID,
    RoleID: studentRole.ID,
}
```

---

## 🛠️ RECURSOS DEL SISTEMA

### Recursos Disponibles:
- **users**: Gestión de usuarios
- **students**: Gestión de estudiantes
- **teachers**: Gestión de profesores
- **courses**: Gestión de cursos
- **grades**: Gestión de calificaciones
- **finances**: Gestión financiera
- **departments**: Gestión de departamentos
- **products**: Gestión de productos
- **orders**: Gestión de pedidos
- **customers**: Gestión de clientes
- **audit**: Auditoría
- **settings**: Configuraciones
- **integrations**: Integraciones

### Acciones Disponibles:
- **create**: Crear recursos
- **read**: Leer/consultar recursos
- **update**: Actualizar recursos
- **delete**: Eliminar recursos
- **invite**: Invitar usuarios
- **suspend**: Suspender usuarios
- **grade**: Calificar estudiantes
- **approve**: Aprobar elementos
- **process**: Procesar elementos

### Alcances de Autorización:
- **own**: Solo recursos propios
- **department**: Recursos del departamento
- **organization**: Recursos de la organización
- **all**: Todos los recursos del sistema

---

## 📁 ARCHIVOS DE REFERENCIA

- **Archivo consolidado**: `all_models.go` - Contiene todos los modelos en un solo archivo
- **Archivos individuales**: Cada modelo en su archivo separado en `/models/`
- **Documentación**: Este archivo para referencia completa

Este sistema proporciona una base sólida y flexible para manejar autenticación multi-tenant con soporte completo para e-commerce, instituciones educativas y empresas.
