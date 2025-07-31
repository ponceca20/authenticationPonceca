# Diagrama Entidad-Relación - Módulo de Autenticación

## 🗂️ Diagrama ER Completo (Actualizado para reflejar la arquitectura exacta)

```mermaid
erDiagram
    %% ENTIDAD CENTRAL DE IDENTIDAD
    IDENTITY {
        string id PK
        string email UK
        string password_hash
        boolean email_verified
        timestamp email_verified_at
        string first_name
        string last_name
        string avatar
        string phone
        date date_of_birth
        timestamp last_login_at
        int failed_login_attempts
        timestamp locked_until
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% ORGANIZACIONES MULTI-TENANT
    ORGANIZATION {
        string id PK
        string name
        string slug UK
        string type
        string address
        string phone
        string website
        string email
        int year_founded
        string ceo
        string industry
        string principal
        json levels
        boolean store_enabled
        string store_name
        json store_categories
        string logo
        string description
        boolean is_active
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% MEMBRESÍAS ORGANIZACIONALES (TABLA PUENTE PRINCIPAL)
    ORGANIZATIONAL_MEMBERSHIP {
        string id PK
        string identity_id FK
        string organization_id FK
        string role_id FK
        string department
        string grade
        string subject
        string student_id
        string employee_id
        timestamp active_from
        timestamp active_until
        boolean is_active
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% PERFILES DE CLIENTE E-COMMERCE
    CUSTOMER_PROFILE {
        string id PK
        string identity_id FK
        string customer_number UK
        string preferred_payment_method
        decimal credit_limit
        decimal total_spent
        int loyalty_points
        boolean accepts_marketing
        string preferred_language
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% PERFILES EXTENDIDOS DE USUARIO
    USER_PROFILE {
        string id PK
        string identity_id FK
        text bio
        json social_links
        json skills
        json interests
        string location
        string website
        json custom_attributes
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% SESIONES DE INVITADOS (SIN IDENTIDAD PERSISTENTE)
    GUEST_SESSION {
        string id PK
        string session_token UK
        string email
        string first_name
        string last_name
        string phone
        longtext cart_data
        timestamp last_activity
        string ip_address
        text user_agent
        timestamp expires_at
        timestamp created_at
    }

    %% ROLES CON JERARQUÍAS
    ROLE {
        string id PK
        string organization_id FK
        string name
        string display_name
        text description
        int hierarchy_level
        boolean is_system_role
        json permissions
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% PERMISOS GRANULARES
    PERMISSION {
        string id PK
        string role_id FK
        string resource
        string action
        string scope
        json conditions
        timestamp created_at
        timestamp updated_at
    }

    %% DEPARTAMENTOS JERÁRQUICOS
    DEPARTMENT {
        string id PK
        string organization_id FK
        string parent_id FK
        string name
        text description
        string manager_id FK
        int hierarchy_level
        boolean is_active
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% SISTEMA DE INVITACIONES
    INVITATION {
        string id PK
        string organization_id FK
        string invited_by_id FK
        string role_id FK
        string email
        string token UK
        string status
        string department
        string grade
        string subject
        text message
        json metadata
        timestamp expires_at
        timestamp accepted_at
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% DIRECCIONES DE ENVÍO
    SHIPPING_ADDRESS {
        string id PK
        string customer_profile_id FK
        string type
        string first_name
        string last_name
        string company
        string address_line_1
        string address_line_2
        string city
        string state
        string postal_code
        string country
        string phone
        boolean is_default
        timestamp deleted_at
        timestamp created_at
        timestamp updated_at
    }

    %% PREFERENCIAS DE CLIENTES
    CUSTOMER_PREFERENCES {
        string id PK
        string customer_profile_id FK
        json notification_preferences
        json privacy_settings
        json marketing_preferences
        string preferred_currency
        string timezone
        json custom_fields
        timestamp created_at
        timestamp updated_at
    }

    %% TOKENS DE REFRESH
    REFRESH_TOKEN {
        string id PK
        string identity_id FK
        string token UK
        text device_info
        string ip_address
        text user_agent
        boolean is_active
        timestamp expires_at
        timestamp last_used_at
        timestamp created_at
        timestamp updated_at
    }

    %% LOGS DE AUDITORÍA (INMUTABLES)
    AUDIT_LOG {
        string id PK
        string organization_id FK
        string identity_id FK
        string action
        string resource
        string resource_id
        json old_values
        json new_values
        string ip_address
        text user_agent
        json metadata
        string status
        timestamp created_at
    }

    %% RELACIONES PRINCIPALES EXACTAS SEGÚN LA ARQUITECTURA
    
    %% 1. Identidad como núcleo central
    IDENTITY ||--o{ ORGANIZATIONAL_MEMBERSHIP : "puede_tener_multiples_membresías"
    IDENTITY ||--o| CUSTOMER_PROFILE : "puede_ser_cliente"
    IDENTITY ||--o| USER_PROFILE : "puede_tener_perfil_extendido"
    IDENTITY ||--o{ REFRESH_TOKEN : "tiene_tokens_activos"
    IDENTITY ||--o{ INVITATION : "puede_invitar_otros"
    IDENTITY ||--o{ AUDIT_LOG : "genera_logs"
    IDENTITY ||--o{ DEPARTMENT : "puede_gestionar_departamentos"
    
    %% 2. Organización como contenedor multi-tenant
    ORGANIZATION ||--o{ ORGANIZATIONAL_MEMBERSHIP : "contiene_miembros"
    ORGANIZATION ||--o{ ROLE : "define_roles"
    ORGANIZATION ||--o{ DEPARTMENT : "estructura_departamental"
    ORGANIZATION ||--o{ INVITATION : "emite_invitaciones"
    ORGANIZATION ||--o{ AUDIT_LOG : "registra_actividades"
    
    %% 3. Membresía como tabla puente principal (CRUCIAL)
    ORGANIZATIONAL_MEMBERSHIP }o--|| IDENTITY : "pertenece_a"
    ORGANIZATIONAL_MEMBERSHIP }o--|| ORGANIZATION : "en_organizacion"
    ORGANIZATIONAL_MEMBERSHIP }o--|| ROLE : "con_rol"
    
    %% 4. Roles y permisos jerárquicos
    ROLE ||--o{ PERMISSION : "tiene_permisos"
    ROLE ||--o{ ORGANIZATIONAL_MEMBERSHIP : "asignado_a_miembros"
    ROLE ||--o{ INVITATION : "para_invitar_con_rol"
    ROLE }o--|| ORGANIZATION : "pertenece_a_org"
    
    %% 5. Estructura departamental jerárquica
    DEPARTMENT ||--o{ DEPARTMENT : "puede_tener_subdepartamentos"
    DEPARTMENT }o--|| ORGANIZATION : "pertenece_a_org"
    DEPARTMENT }o--o| IDENTITY : "gestionado_por"
    
    %% 6. Sistema de invitaciones completo
    INVITATION }o--|| ORGANIZATION : "para_organizacion"
    INVITATION }o--|| IDENTITY : "enviada_por"
    INVITATION }o--|| ROLE : "con_rol_especifico"
    
    %% 7. Cliente e-commerce y direcciones
    CUSTOMER_PROFILE ||--o{ SHIPPING_ADDRESS : "tiene_direcciones"
    CUSTOMER_PROFILE ||--o| CUSTOMER_PREFERENCES : "configuracion_personalizada"
    CUSTOMER_PROFILE }o--|| IDENTITY : "extension_de_identidad"
    
    %% 8. Perfil extendido de usuario
    USER_PROFILE }o--|| IDENTITY : "extension_profesional"
    
    %% 9. Tokens de sesión
    REFRESH_TOKEN }o--|| IDENTITY : "para_identidad"
    
    %% 10. Auditoría completa
    AUDIT_LOG }o--o| ORGANIZATION : "en_contexto_org"
    AUDIT_LOG }o--o| IDENTITY : "realizada_por"
    
    %% 11. Sesiones temporales (independientes)
    %% GUEST_SESSION es independiente - no tiene relaciones FK
```

## 📊 Verificación de Consistencia Arquitectónica

### ✅ **Consistencias Confirmadas**

#### **1. Identidad Centralizada**
- ✅ **Diagrama ER**: IDENTITY como entidad central con relaciones 1:N y 1:1
- ✅ **Estructura**: `models/identity.go` como modelo principal
- ✅ **Arquitectura**: Una sola fuente de verdad para usuarios

#### **2. Multi-Tenancy Organizacional**
- ✅ **Diagrama ER**: ORGANIZATION como contenedor con ORGANIZATIONAL_MEMBERSHIP como puente
- ✅ **Estructura**: Submódulo `organization/` y tabla puente bien definida
- ✅ **Arquitectura**: Aislamiento total entre organizaciones

#### **3. RBAC+ Jerárquico**
- ✅ **Diagrama ER**: ROLE → PERMISSION con hierarchy_level
- ✅ **Estructura**: Submódulo `role/` con sistema granular
- ✅ **Arquitectura**: Permisos heredados por jerarquía

#### **4. E-commerce Integrado**
- ✅ **Diagrama ER**: CUSTOMER_PROFILE como extensión opcional de IDENTITY
- ✅ **Estructura**: Submódulo `customer/` con direcciones y preferencias
- ✅ **Arquitectura**: Contexto dual (organizacional + comercial)

### 🔧 **Correcciones Aplicadas al Diagrama**

#### **1. Tabla Puente Principal**
```mermaid
ORGANIZATIONAL_MEMBERSHIP {
    string identity_id FK      -- ✅ Conecta con IDENTITY
    string organization_id FK  -- ✅ Conecta con ORGANIZATION  
    string role_id FK         -- ✅ Conecta con ROLE
    string department         -- ✅ Contexto específico
    string grade             -- ✅ Para estudiantes
    string subject           -- ✅ Para profesores
}
```

#### **2. Sesiones de Invitados Independientes**
```mermaid
GUEST_SESSION {
    -- ✅ SIN foreign keys - completamente independiente
    -- ✅ Convertible a IDENTITY vía servicio de conversión
    -- ✅ Auto-expirable sin soft delete
}
```

#### **3. Auditoría Inmutable**
```mermaid
AUDIT_LOG {
    -- ✅ Solo timestamp created_at (inmutable)
    -- ✅ Referencias opcionales a organization e identity
    -- ✅ Campos old_values y new_values para compliance
}
```

#### **4. Jerarquías Departamentales**
```mermaid
DEPARTMENT {
    string parent_id FK       -- ✅ Auto-referencia para jerarquías
    string manager_id FK      -- ✅ Apunta a IDENTITY
    int hierarchy_level       -- ✅ Nivel en la jerarquía
}
```

## 🎯 **Mapeo Exacto: Diagrama ↔ Estructura**

### **Submódulos → Entidades Principales**

| Submódulo | Entidad Principal | Entidades Relacionadas |
|-----------|------------------|------------------------|
| `auth/` | IDENTITY | REFRESH_TOKEN, AUDIT_LOG |
| `organization/` | ORGANIZATION | ORGANIZATIONAL_MEMBERSHIP |
| `user/` | ORGANIZATIONAL_MEMBERSHIP | IDENTITY, ROLE, DEPARTMENT |
| `customer/` | CUSTOMER_PROFILE | SHIPPING_ADDRESS, CUSTOMER_PREFERENCES |
| `guest/` | GUEST_SESSION | *(independiente)* |
| `profile/` | USER_PROFILE | IDENTITY |
| `role/` | ROLE | PERMISSION |
| `department/` | DEPARTMENT | *(auto-referencial)* |
| `invitation/` | INVITATION | ORGANIZATION, ROLE, IDENTITY |
| `audit/` | AUDIT_LOG | *(inmutable)* |

### **Middleware → Entidades de Apoyo**

| Middleware | Entidades que Utiliza |
|------------|----------------------|
| `smart_auth.go` | IDENTITY, ORGANIZATIONAL_MEMBERSHIP, CUSTOMER_PROFILE, GUEST_SESSION |
| `permissions.go` | ROLE, PERMISSION, ORGANIZATIONAL_MEMBERSHIP |
| `org_context.go` | ORGANIZATION, ORGANIZATIONAL_MEMBERSHIP |
| `audit.go` | AUDIT_LOG |

### **Utils → Operaciones Transversales**

| Utilidad | Entidades Afectadas |
|----------|-------------------|
| `jwt.go` | IDENTITY, REFRESH_TOKEN |
| `validation.go` | *(todas las entidades)* |
| `cache.go` | IDENTITY, ORGANIZATION, ROLE |

## 🔄 **Flujos de Datos Verificados**

### **1. Registro de Usuario Empresarial**
```
1. Crear IDENTITY
2. Crear ORGANIZATION  
3. Crear ROLE (si no existe)
4. Crear ORGANIZATIONAL_MEMBERSHIP
5. Crear AUDIT_LOG
```

### **2. Conversión Guest → Customer**
```
1. Leer GUEST_SESSION
2. Crear IDENTITY
3. Crear CUSTOMER_PROFILE
4. Transferir cart_data
5. Eliminar GUEST_SESSION
6. Crear AUDIT_LOG
```

### **3. Invitación Organizacional**
```
1. Validar ORGANIZATION y ROLE
2. Crear INVITATION
3. Enviar email
4. Al aceptar: crear ORGANIZATIONAL_MEMBERSHIP
5. Crear AUDIT_LOG
```

## ✅ **Conclusión de Verificación**

### **DIAGRAMA ER ✅ PERFECTAMENTE ALINEADO**
- ✅ Todas las entidades del diagrama tienen su modelo Go correspondiente
- ✅ Todas las relaciones están correctamente implementadas
- ✅ Los submódulos mapean exactamente a las entidades principales
- ✅ La arquitectura RBAC+ está completamente reflejada
- ✅ El multi-tenancy está correctamente modelado

### **ESTRUCTURA DEL MÓDULO ✅ CONSISTENTE**
- ✅ Cada submódulo maneja las entidades que le corresponden
- ✅ Los middleware utilizan las entidades apropiadas
- ✅ Las utilidades son transversales como debe ser
- ✅ Los flujos de datos siguen las relaciones del diagrama

### **ARQUITECTURA ✅ ROBUSTA**
- ✅ Identidad centralizada como núcleo
- ✅ Multi-tenancy seguro vía ORGANIZATION
- ✅ RBAC+ jerárquico completamente implementado
- ✅ E-commerce integrado sin conflictos
- ✅ Auditoría completa e inmutable

**El diagrama ER y la estructura del módulo están en perfecta sincronía. La arquitectura es consistente y escalable.**
