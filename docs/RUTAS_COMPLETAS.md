# Rutas Completas del Módulo de Autenticación

## 📋 RESUMEN EJECUTIVO

Este documento define todas las rutas que implementaremos en el módulo de autenticación basado en la arquitectura de submódulos definida en `GUIA_DESARROLLO_MODULO_AUTH.md`.

**Total de rutas**: 89 endpoints organizados en 13 categorías principales.

---

## 🔐 1. AUTENTICACIÓN CORE (`/api/v1/auth`)

### Autenticación Básica
```
POST   /api/v1/auth/register              # Registro de usuario empresarial
POST   /api/v1/auth/login                 # Login empresarial
POST   /api/v1/auth/logout                # Logout con invalidación de token
POST   /api/v1/auth/refresh               # Renovar access token con refresh token
GET    /api/v1/auth/me                    # Obtener perfil del usuario autenticado
PUT    /api/v1/auth/me                    # Actualizar perfil básico
```

### Recuperación de Contraseña
```
POST   /api/v1/auth/forgot-password       # Enviar email de recuperación
POST   /api/v1/auth/reset-password        # Restablecer contraseña con token
POST   /api/v1/auth/change-password       # Cambiar contraseña (usuario autenticado)
```

### Verificación de Email
```
POST   /api/v1/auth/send-verification     # Reenviar email de verificación
POST   /api/v1/auth/verify-email          # Verificar email con token
```

---

## 🏢 2. AUTENTICACIÓN ORGANIZACIONAL (`/api/v1/org/:slug/auth`)

### Contexto Organizacional
```
POST   /api/v1/org/:slug/auth/login       # Login con contexto organizacional
POST   /api/v1/org/:slug/auth/logout      # Logout del contexto organizacional
GET    /api/v1/org/:slug/auth/me          # Perfil en contexto organizacional
POST   /api/v1/org/:slug/auth/refresh     # Renovar token organizacional
PUT    /api/v1/org/:slug/auth/switch      # Cambiar rol dentro de la organización
```

---

## 👥 3. GESTIÓN DE USUARIOS EMPRESARIALES (`/api/v1/users`)

### CRUD de Usuarios
```
GET    /api/v1/users                      # Listar usuarios (con filtros)
POST   /api/v1/users                      # Crear usuario empresarial
GET    /api/v1/users/:id                  # Obtener usuario específico
PUT    /api/v1/users/:id                  # Actualizar usuario completo
PATCH  /api/v1/users/:id                  # Actualización parcial
DELETE /api/v1/users/:id                  # Eliminar usuario (soft delete)
```

### Gestión de Estado
```
PUT    /api/v1/users/:id/status           # Cambiar estado (activo/inactivo/suspendido)
PUT    /api/v1/users/:id/password         # Cambiar contraseña de usuario
POST   /api/v1/users/:id/unlock           # Desbloquear usuario después de intentos fallidos
```

### Bulk Operations
```
POST   /api/v1/users/bulk-create          # Creación masiva de usuarios
PUT    /api/v1/users/bulk-update          # Actualización masiva
DELETE /api/v1/users/bulk-delete          # Eliminación masiva
POST   /api/v1/users/import               # Importar usuarios desde CSV/Excel
```

---

## 🛒 4. CLIENTES E-COMMERCE (`/api/v1/customers`)

### Autenticación de Clientes
```
POST   /api/v1/customers/register         # Registro de cliente e-commerce
POST   /api/v1/customers/login            # Login de cliente
POST   /api/v1/customers/logout           # Logout de cliente
POST   /api/v1/customers/refresh          # Renovar token de cliente
```

### Gestión de Perfil
```
GET    /api/v1/customers/profile          # Obtener perfil de cliente
PUT    /api/v1/customers/profile          # Actualizar perfil completo
PATCH  /api/v1/customers/profile          # Actualización parcial de perfil
DELETE /api/v1/customers/profile          # Eliminar cuenta de cliente
```

### Direcciones de Envío
```
GET    /api/v1/customers/addresses        # Listar direcciones del cliente
POST   /api/v1/customers/addresses        # Agregar nueva dirección
GET    /api/v1/customers/addresses/:id    # Obtener dirección específica
PUT    /api/v1/customers/addresses/:id    # Actualizar dirección
DELETE /api/v1/customers/addresses/:id    # Eliminar dirección
PUT    /api/v1/customers/addresses/:id/default  # Establecer como dirección por defecto
```

### Preferencias
```
GET    /api/v1/customers/preferences      # Obtener preferencias del cliente
PUT    /api/v1/customers/preferences      # Actualizar preferencias
GET    /api/v1/customers/loyalty          # Estado del programa de lealtad
```

---

## 👤 5. SESIONES DE INVITADOS (`/api/v1/guests`)

### Gestión de Sesiones
```
POST   /api/v1/guests/session             # Crear nueva sesión de invitado
GET    /api/v1/guests/session             # Obtener sesión actual
PUT    /api/v1/guests/session             # Actualizar datos de sesión
DELETE /api/v1/guests/session             # Eliminar sesión de invitado
```

### Conversiones
```
POST   /api/v1/guests/convert-to-customer # Convertir invitado a cliente registrado
POST   /api/v1/guests/convert-to-user     # Convertir invitado a usuario empresarial
```

### Tracking
```
POST   /api/v1/guests/track-activity      # Registrar actividad del invitado
GET    /api/v1/guests/activity            # Obtener historial de actividad
```

---

## 🏛️ 6. ORGANIZACIONES (`/api/v1/organizations`)

### CRUD de Organizaciones
```
GET    /api/v1/organizations              # Listar organizaciones (con permisos)
POST   /api/v1/organizations              # Crear nueva organización
GET    /api/v1/organizations/:id          # Obtener organización específica
PUT    /api/v1/organizations/:id          # Actualizar organización
DELETE /api/v1/organizations/:id          # Eliminar organización
```

### Gestión de Membresías
```
GET    /api/v1/organizations/:id/members  # Listar miembros de la organización
POST   /api/v1/organizations/:id/members  # Agregar miembro a organización
DELETE /api/v1/organizations/:id/members/:userId  # Remover miembro
```

### Configuración
```
GET    /api/v1/organizations/:id/settings # Obtener configuración de organización
PUT    /api/v1/organizations/:id/settings # Actualizar configuración
GET    /api/v1/organizations/:id/stats    # Estadísticas de la organización
```

---

## 🎭 7. ROLES Y PERMISOS (`/api/v1/org/:slug/roles`)

### CRUD de Roles
```
GET    /api/v1/org/:slug/roles            # Listar roles de la organización
POST   /api/v1/org/:slug/roles            # Crear nuevo rol
GET    /api/v1/org/:slug/roles/:id        # Obtener rol específico
PUT    /api/v1/org/:slug/roles/:id        # Actualizar rol
DELETE /api/v1/org/:slug/roles/:id        # Eliminar rol
```

### Gestión de Permisos
```
GET    /api/v1/org/:slug/roles/:id/permissions     # Listar permisos del rol
PUT    /api/v1/org/:slug/roles/:id/permissions     # Asignar permisos al rol
POST   /api/v1/org/:slug/roles/:id/permissions/add # Agregar permisos específicos
DELETE /api/v1/org/:slug/roles/:id/permissions/:permissionId  # Remover permiso
```

### Asignación de Roles
```
GET    /api/v1/org/:slug/roles/:id/users  # Usuarios con este rol
POST   /api/v1/org/:slug/roles/:id/users  # Asignar rol a usuarios
DELETE /api/v1/org/:slug/roles/:id/users/:userId  # Remover rol de usuario
```

### Sistema de Roles
```
GET    /api/v1/roles/system               # Listar roles del sistema
GET    /api/v1/permissions                # Listar todos los permisos disponibles
GET    /api/v1/permissions/resources      # Listar recursos protegidos
```

---

## 🏛️ 8. DEPARTAMENTOS (`/api/v1/org/:slug/departments`)

### CRUD de Departamentos
```
GET    /api/v1/org/:slug/departments      # Listar departamentos
POST   /api/v1/org/:slug/departments      # Crear departamento
GET    /api/v1/org/:slug/departments/:id  # Obtener departamento
PUT    /api/v1/org/:slug/departments/:id  # Actualizar departamento
DELETE /api/v1/org/:slug/departments/:id  # Eliminar departamento
```

### Jerarquía Departamental
```
GET    /api/v1/org/:slug/departments/:id/children    # Sub-departamentos
GET    /api/v1/org/:slug/departments/:id/parent      # Departamento padre
PUT    /api/v1/org/:slug/departments/:id/move        # Mover departamento en jerarquía
```

### Gestión de Personal
```
GET    /api/v1/org/:slug/departments/:id/users       # Usuarios del departamento
POST   /api/v1/org/:slug/departments/:id/users       # Asignar usuarios
DELETE /api/v1/org/:slug/departments/:id/users/:userId  # Remover usuario
```

---

## 📧 9. INVITACIONES (`/api/v1/org/:slug/invitations` y `/api/v1/invitations`)

### Gestión de Invitaciones Organizacionales
```
GET    /api/v1/org/:slug/invitations      # Listar invitaciones de la organización
POST   /api/v1/org/:slug/invitations      # Crear nueva invitación
GET    /api/v1/org/:slug/invitations/:id  # Obtener invitación específica
PUT    /api/v1/org/:slug/invitations/:id  # Actualizar invitación
DELETE /api/v1/org/:slug/invitations/:id  # Cancelar invitación
POST   /api/v1/org/:slug/invitations/:id/resend  # Reenviar invitación
```

### Procesamiento de Invitaciones (Sin contexto organizacional)
```
GET    /api/v1/invitations/token/:token           # Validar token de invitación
POST   /api/v1/invitations/token/:token/accept    # Aceptar invitación
POST   /api/v1/invitations/token/:token/decline   # Rechazar invitación
GET    /api/v1/invitations/token/:token/info      # Info de la invitación
```

### Invitaciones Masivas
```
POST   /api/v1/org/:slug/invitations/bulk         # Envío masivo de invitaciones
POST   /api/v1/org/:slug/invitations/import       # Importar lista de invitaciones
```

---

## 📊 10. PERFILES EXTENDIDOS (`/api/v1/profiles`)

### Gestión de Perfiles
```
GET    /api/v1/profiles/me                # Mi perfil completo
PUT    /api/v1/profiles/me                # Actualizar perfil completo
PATCH  /api/v1/profiles/me                # Actualización parcial
```

### Archivos y Media
```
POST   /api/v1/profiles/avatar            # Subir avatar
DELETE /api/v1/profiles/avatar            # Eliminar avatar
POST   /api/v1/profiles/documents         # Subir documentos del perfil
GET    /api/v1/profiles/documents         # Listar documentos
DELETE /api/v1/profiles/documents/:id     # Eliminar documento
```

### Preferencias
```
GET    /api/v1/profiles/preferences       # Obtener preferencias del usuario
PUT    /api/v1/profiles/preferences       # Actualizar preferencias
GET    /api/v1/profiles/privacy           # Configuración de privacidad
PUT    /api/v1/profiles/privacy           # Actualizar configuración de privacidad
```

---

## 📋 11. AUDITORÍA (`/api/v1/org/:slug/audits`)

### Consulta de Logs
```
GET    /api/v1/org/:slug/audits           # Listar logs de auditoría (con filtros)
GET    /api/v1/org/:slug/audits/:id       # Obtener log específico
GET    /api/v1/org/:slug/audits/user/:userId      # Logs por usuario
GET    /api/v1/org/:slug/audits/actions/:action   # Logs por tipo de acción
GET    /api/v1/org/:slug/audits/resources/:resource  # Logs por recurso
```

### Reportes y Exportación
```
GET    /api/v1/org/:slug/audits/report    # Generar reporte de auditoría
GET    /api/v1/org/:slug/audits/export    # Exportar logs (CSV/PDF)
GET    /api/v1/org/:slug/audits/stats     # Estadísticas de auditoría
```

### Auditoría del Sistema
```
GET    /api/v1/audit/system               # Logs del sistema (solo super admin)
GET    /api/v1/audit/security             # Eventos de seguridad
GET    /api/v1/audit/failed-logins        # Intentos de login fallidos
```

---

## 🔗 12. RUTAS UNIFICADAS (`/api/v1/unified`)

### Dashboard y Búsqueda
```
GET    /api/v1/unified/dashboard          # Dashboard unificado multi-contexto
GET    /api/v1/unified/search             # Búsqueda global en todos los contextos
GET    /api/v1/unified/notifications      # Notificaciones unificadas
PUT    /api/v1/unified/notifications/:id/read  # Marcar notificación como leída
```

### Operaciones Masivas
```
POST   /api/v1/unified/bulk-actions       # Acciones masivas cross-módulo
GET    /api/v1/unified/contexts           # Listar todos los contextos del usuario
POST   /api/v1/unified/switch-context     # Cambiar contexto activo
```

### Analytics Unificados
```
GET    /api/v1/unified/analytics          # Analytics unificados del usuario
GET    /api/v1/unified/activity-feed      # Feed de actividad global
GET    /api/v1/unified/recommendations    # Recomendaciones personalizadas
```

---

## 🔗 13. INTEGRACIÓN CON OTROS MÓDULOS (`/api/v1/integration`)

### Validación de Tokens
```
POST   /api/v1/integration/validate-token    # Validar token para otros módulos
GET    /api/v1/integration/user-info/:id     # Información de usuario para integración
POST   /api/v1/integration/check-permission  # Verificar permisos específicos
GET    /api/v1/integration/user-context/:id  # Contexto completo del usuario
```

### Eventos y Webhooks
```
POST   /api/v1/integration/webhooks       # Registrar webhook de otros módulos
GET    /api/v1/integration/webhooks       # Listar webhooks registrados
DELETE /api/v1/integration/webhooks/:id   # Eliminar webhook
POST   /api/v1/integration/events         # Enviar evento desde otro módulo
```

### Sincronización
```
POST   /api/v1/integration/sync-user      # Sincronizar datos de usuario
POST   /api/v1/integration/sync-org       # Sincronizar datos organizacionales
GET    /api/v1/integration/health         # Estado de salud del módulo de auth
```

---

## 🔧 14. ADMINISTRACIÓN DEL SISTEMA (`/api/v1/admin`)

### Gestión de Sistema
```
GET    /api/v1/admin/stats                # Estadísticas generales del sistema
GET    /api/v1/admin/users                # Gestión de usuarios (super admin)
GET    /api/v1/admin/organizations        # Gestión de organizaciones
POST   /api/v1/admin/maintenance          # Activar modo mantenimiento
```

### Configuración Global
```
GET    /api/v1/admin/config               # Configuración global del sistema
PUT    /api/v1/admin/config               # Actualizar configuración global
GET    /api/v1/admin/features             # Feature flags del sistema
PUT    /api/v1/admin/features             # Actualizar feature flags
```

---

## 🏥 15. HEALTH CHECK Y MONITOREO

### Health Checks
```
GET    /health                            # Estado básico del servicio
GET    /health/detailed                   # Estado detallado con dependencias
GET    /health/ready                      # Readiness probe (Kubernetes)
GET    /health/live                       # Liveness probe (Kubernetes)
```

### Métricas
```
GET    /metrics                           # Métricas en formato Prometheus
GET    /api/v1/monitoring/performance     # Métricas de performance
GET    /api/v1/monitoring/usage           # Métricas de uso
```

---

## 📊 RESUMEN POR CATEGORÍAS

| Categoría | Endpoints | Descripción |
|-----------|-----------|-------------|
| Autenticación Core | 9 | Login, registro, recuperación de contraseña |
| Auth Organizacional | 5 | Contexto organizacional específico |
| Usuarios Empresariales | 11 | CRUD y gestión de usuarios empresariales |
| Clientes E-commerce | 12 | Gestión completa de clientes |
| Invitados | 6 | Sesiones temporales y conversiones |
| Organizaciones | 8 | CRUD y gestión organizacional |
| Roles y Permisos | 12 | Sistema RBAC+ completo |
| Departamentos | 8 | Jerarquías departamentales |
| Invitaciones | 9 | Sistema de invitaciones |
| Perfiles | 9 | Gestión de perfiles extendidos |
| Auditoría | 8 | Logs y reportes de auditoría |
| Rutas Unificadas | 8 | Dashboard y operaciones cross-módulo |
| Integración | 9 | APIs para otros módulos |
| Administración | 7 | Gestión del sistema |
| Health Check | 6 | Monitoreo y métricas |

**TOTAL: 127 endpoints**

---

## 🎯 PRIORIDAD DE IMPLEMENTACIÓN

### Fase 1 - Core Crítico (Semana 1)
1. **Autenticación Core** - Login/registro básico
2. **Organizaciones** - CRUD básico
3. **Usuarios Empresariales** - CRUD básico
4. **Health Check** - Monitoreo básico

### Fase 2 - Funcionalidad Empresarial (Semana 2)
1. **Auth Organizacional** - Contexto multi-tenant
2. **Roles y Permisos** - RBAC básico
3. **Invitaciones** - Sistema de invitaciones
4. **Departamentos** - Estructura organizacional

### Fase 3 - E-commerce y Avanzado (Semana 3)
1. **Clientes E-commerce** - Sistema de clientes
2. **Invitados** - Sesiones temporales
3. **Perfiles** - Gestión avanzada
4. **Auditoría** - Logging y reportes

### Fase 4 - Integración y Unificación (Semana 4)
1. **Rutas Unificadas** - Dashboard global
2. **Integración** - APIs para otros módulos
3. **Administración** - Gestión del sistema
4. **Optimización** - Performance y seguridad

---

## 🔐 CONSIDERACIONES DE SEGURIDAD

### Autenticación Requerida
- **Todas las rutas** excepto health checks requieren autenticación
- **JWT unificado** maneja múltiples contextos
- **Rate limiting** específico por tipo de usuario

### Niveles de Autorización
- **Public**: Health checks, documentación
- **Guest**: Sesiones temporales, conversión
- **Customer**: E-commerce y perfil básico
- **User**: Funcionalidad empresarial/educativa
- **Admin**: Gestión organizacional
- **Super Admin**: Administración del sistema

### Middleware Stack
```
Request → Rate Limiting → CORS → JWT Validation → RBAC → Audit → Handler
```

---

## 📝 NOTAS DE IMPLEMENTACIÓN

1. **Versionado**: Todas las rutas principales bajo `/api/v1/`
2. **Consistencia**: DTOs estandarizados para todas las operaciones
3. **Documentación**: Swagger/OpenAPI auto-generado
4. **Testing**: Suite de tests completa por categoría
5. **Monitoreo**: Logs estructurados y métricas automáticas

Este documento sirve como la **especificación completa** para la implementación del módulo de autenticación. Cada endpoint está diseñado para soportar los casos de uso empresariales, educativos y de e-commerce definidos en la guía de desarrollo.
