# Estructura Completa del Módulo Authentication

## 📋 Tabla de Contenidos
1. [Estructura de Archivos Completa](#estructura-de-archivos-completa)
2. [Mapeo Exacto: Submódulos ↔ Entidades ER](#mapeo-exacto-submódulos--entidades-er)
3. [Verificación de Consistencia Arquitectónica](#verificación-de-consistencia-arquitectónica)
4. [Flujos de Trabajo Validados](#flujos-de-trabajo-validados)

## 🏗️ Estructura de Archivos Completa (✅ Verificada vs Diagrama ER)

```
proyecto/
├── cmd/main.go                           # ✅ MANTENER (con modificaciones menores)
├── config/config.go                      # ✅ EXTENDER para JWT settings
├── database/database.go                  # ✅ MANTENER
├── registry/                             # ✅ USAR sistema existente
│   ├── registry.go                      # Sistema de registro de módulos
│   ├── migrations.go                    # Sistema de migraciones
│   ├── cache/redis.go                   # Cache Redis existente
│   └── queue/rabbitmq.go                # Cola RabbitMQ existente
├── module/                               # ✅ ESTRUCTURA PRINCIPAL
│   └── authentication/                  # 🆕 MÓDULO DE AUTENTICACIÓN COMPLETO
│       ├── init.go                      # 🆕 Registro del módulo
│       ├── routes.go                    # 🆕 Configuración de rutas centralizadas
│       ├── migrations.go                # 🆕 Migraciones del módulo (✅ Refleja ER exacto)
│       │
│       ├── auth/                        # 🔐 MANEJA: IDENTITY + REFRESH_TOKEN
│       │   ├── handler.go              # Endpoints JWT y autenticación
│       │   ├── service.go              # Lógica de tokens y validación
│       │   ├── repository.go           # Consultas IDENTITY + REFRESH_TOKEN
│       │   └── dto.go                  # LoginDTO, TokenResponseDTO
│       │
│       ├── organization/                # 🏢 MANEJA: ORGANIZATION (multi-tenant)
│       │   ├── handler.go              # CRUD organizaciones
│       │   ├── service.go              # Lógica multi-tenancy + store config
│       │   ├── repository.go           # Consultas ORGANIZATION
│       │   └── dto.go                  # OrganizationDTO, StoreConfigDTO
│       │
│       ├── user/                        # 👥 MANEJA: ORGANIZATIONAL_MEMBERSHIP (tabla puente clave)
│       │   ├── handler.go              # Gestión de miembros organizacionales
│       │   ├── service.go              # Lógica de membresías y contextos
│       │   ├── repository.go           # Consultas ORGANIZATIONAL_MEMBERSHIP complejas
│       │   └── dto.go                  # MembershipDTO, ContextDTO
│       │
│       ├── customer/                    # 🛒 MANEJA: CUSTOMER_PROFILE + SHIPPING_ADDRESS + CUSTOMER_PREFERENCES
│       │   ├── handler.go              # E-commerce: registro, perfil, direcciones
│       │   ├── service.go              # Lógica de clientes y conversiones
│       │   ├── repository.go           # Consultas CUSTOMER_PROFILE + relaciones
│       │   └── dto.go                  # CustomerDTO, AddressDTO, PreferencesDTO
│       │
│       ├── guest/                       # 👤 MANEJA: GUEST_SESSION (independiente)
│       │   ├── handler.go              # Sesiones temporales sin FK
│       │   ├── service.go              # Lógica de expiración y conversión
│       │   ├── repository.go           # Consultas GUEST_SESSION + limpieza
│       │   └── dto.go                  # GuestSessionDTO, ConversionDTO
│       │
│       ├── profile/                     # 📊 MANEJA: USER_PROFILE (extensión de IDENTITY)
│       │   ├── handler.go              # Perfiles extendidos y avatares
│       │   ├── service.go              # Validación de perfiles multimedia
│       │   ├── repository.go           # Consultas USER_PROFILE
│       │   └── dto.go                  # ProfileDTO, AvatarDTO
│       │
│       ├── role/                        # 🎭 MANEJA: ROLE + PERMISSION (RBAC+ jerárquico)
│       │   ├── handler.go              # CRUD roles + asignación permisos
│       │   ├── service.go              # Lógica RBAC+ con herencia
│       │   ├── repository.go           # Consultas ROLE + PERMISSION complejas
│       │   └── dto.go                  # RoleDTO, PermissionDTO, HierarchyDTO
│       │
│       ├── department/                  # 🏛️ MANEJA: DEPARTMENT (auto-referencial)
│       │   ├── handler.go              # Estructuras organizacionales jerárquicas
│       │   ├── service.go              # Lógica de jerarquías + managers
│       │   ├── repository.go           # Consultas DEPARTMENT + auto-referencia
│       │   └── dto.go                  # DepartmentDTO, HierarchyDTO
│       │
│       ├── invitation/                  # 📧 MANEJA: INVITATION (multi-FK)
│       │   ├── handler.go              # Sistema completo de invitaciones
│       │   ├── service.go              # Tokens temporales + expiración
│       │   ├── repository.go           # Consultas INVITATION + relaciones
│       │   └── dto.go                  # InvitationDTO, AcceptDTO
│       │
│       ├── audit/                       # 📋 MANEJA: AUDIT_LOG (inmutable)
│       │   ├── handler.go              # Consulta de logs + reportes
│       │   ├── service.go              # Registro automático + clasificación
│       │   ├── repository.go           # Consultas AUDIT_LOG optimizadas
│       │   └── dto.go                  # AuditLogDTO, ReportDTO
│       │
│       ├── unified/                     # 🔗 MANEJA: Múltiples entidades (cross-submódulo)
│       │   ├── handler.go              # Dashboard + búsquedas globales
│       │   ├── service.go              # Lógica que combina entidades
│       │   ├── repository.go           # Consultas cross-entity complejas
│       │   └── dto.go                  # UnifiedDTO, DashboardDTO
│       │
│       ├── models/                      # 📊 MODELOS ER EXACTOS
│       │   ├── identity.go             # ✅ IDENTITY (núcleo central)
│       │   ├── organization.go         # ✅ ORGANIZATION (multi-tenant)
│       │   ├── organizational_membership.go # ✅ Tabla puente principal
│       │   ├── customer_profile.go     # ✅ CUSTOMER_PROFILE (e-commerce)
│       │   ├── user_profile.go         # ✅ USER_PROFILE (extensión)
│       │   ├── guest_session.go        # ✅ GUEST_SESSION (independiente)
│       │   ├── role.go                 # ✅ ROLE (RBAC+)
│       │   ├── permission.go           # ✅ PERMISSION (granular)
│       │   ├── department.go           # ✅ DEPARTMENT (jerárquico)
│       │   ├── invitation.go           # ✅ INVITATION (multi-FK)
│       │   ├── shipping_address.go     # ✅ SHIPPING_ADDRESS
│       │   ├── customer_preferences.go # ✅ CUSTOMER_PREFERENCES
│       │   ├── refresh_token.go        # ✅ REFRESH_TOKEN
│       │   └── audit_log.go            # ✅ AUDIT_LOG (inmutable)
│       │
│       ├── middleware/                  # 🛡️ UTILIZA MÚLTIPLES ENTIDADES
│       │   ├── auth.go                 # ✅ IDENTITY + REFRESH_TOKEN
│       │   ├── smart_auth.go           # ✅ IDENTITY + ORG_MEMBERSHIP + CUSTOMER_PROFILE + GUEST_SESSION
│       │   ├── org_context.go          # ✅ ORGANIZATION + ORG_MEMBERSHIP
│       │   ├── permissions.go          # ✅ ROLE + PERMISSION + ORG_MEMBERSHIP
│       │   ├── rate_limit.go           # ✅ Basado en tipo de IDENTITY
│       │   ├── audit.go                # ✅ Crea AUDIT_LOG automáticamente
│       │   └── cors.go                 # ✅ Configuración e-commerce
│       │
│       ├── utils/                       # 🔧 UTILIDADES TRANSVERSALES
│       │   ├── jwt.go                  # ✅ Claims multi-contexto
│       │   ├── password.go             # ✅ Para IDENTITY.password_hash
│       │   ├── validation.go           # ✅ Para todos los modelos
│       │   ├── response.go             # ✅ Respuestas estándar HTTP
│       │   ├── cache.go                # ✅ Cache de IDENTITY + ROLE
│       │   ├── email.go                # ✅ Para INVITATION
│       │   ├── upload.go               # ✅ Para avatares en IDENTITY/USER_PROFILE
│       │   └── geo.go                  # ✅ Para SHIPPING_ADDRESS
│       │
│       ├── integration/                 # 🔗 HOOKS EXTERNOS
│       │   ├── ecommerce.go            # ✅ Eventos desde CUSTOMER_PROFILE
│       │   ├── delivery.go             # ✅ Integración con SHIPPING_ADDRESS
│       │   ├── payment.go              # ✅ Hooks desde CUSTOMER_PROFILE
│       │   ├── notification.go         # ✅ Para INVITATION + AUDIT_LOG
│       │   └── events.go               # ✅ Sistema de eventos general
│       │
│       └── test/                        # 🧪 TESTS COMPLETOS
│           ├── setup_test.go           # ✅ Configuración de BD de prueba
│           ├── identity_test.go        # ✅ Tests de IDENTITY
│           ├── organization_test.go    # ✅ Tests de ORGANIZATION
│           ├── membership_test.go      # ✅ Tests de ORG_MEMBERSHIP
│           ├── customer_test.go        # ✅ Tests de CUSTOMER_PROFILE
│           ├── guest_test.go           # ✅ Tests de GUEST_SESSION
│           ├── role_test.go            # ✅ Tests de ROLE + PERMISSION
│           ├── invitation_test.go      # ✅ Tests de INVITATION
│           ├── audit_test.go           # ✅ Tests de AUDIT_LOG
│           ├── middleware_test.go      # ✅ Tests de middleware
│           ├── jwt_test.go             # ✅ Tests de tokens
│           └── integration_test.go     # ✅ Tests end-to-end
│
└── addModules/modules.go                # ✅ REGISTRAR módulo aquí
```

## 🎯 **Mapeo Exacto: Submódulos ↔ Entidades ER**

### **Entidades Principales → Submódulo