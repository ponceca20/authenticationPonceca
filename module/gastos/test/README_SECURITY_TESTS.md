# 🛡️ Tests de Seguridad - Módulo Gastos

Este archivo contiene tests completos de seguridad que evalúan tanto la **autenticación** como el **RBAC (Role-Based Access Control)** del módulo gastos contra el backend real.

## 📋 Características de los Tests

### 🔐 **Tests de Autenticación**
- ✅ Verificación de acceso sin token (debe denegar)
- ✅ Verificación de tokens inválidos (debe denegar)
- ✅ Verificación de tokens malformados (debe denegar)
- ✅ Verificación de tokens válidos (debe permitir acceso)

### 👥 **Tests de RBAC (Autorización)**
- ✅ **Administradores**: Acceso completo a todas las rutas
- ✅ **Usuarios Regulares**: Acceso limitado a sus propios recursos
- ✅ **Usuarios Sin Permisos**: Acceso denegado a recursos protegidos

### 🌐 **Tests por Endpoint**
- ✅ `GET /api/org/:slug/expenses` - Listar gastos
- ✅ `POST /api/org/:slug/expenses` - Crear gastos
- ✅ `PUT /api/org/:slug/expenses/:id` - Actualizar gastos
- ✅ `DELETE /api/org/:slug/expenses/:id` - Eliminar gastos
- ✅ `GET /api/org/:slug/admin/expenses/all` - Rutas administrativas

### 🔒 **Tests de Seguridad Avanzada**
- ✅ **Context Injection**: Prevención de acceso a organizaciones no autorizadas
- ✅ **Rate Limiting**: Prevención de ataques de fuerza bruta
- ✅ **Privilege Escalation**: Verificación de que los usuarios no puedan escalar privilegios

## 🚀 Cómo Ejecutar los Tests

### Opción 1: Script de PowerShell (Windows)
```powershell
.\run_security_tests.ps1
```

### Opción 2: Script de Bash (Linux/Mac)
```bash
chmod +x run_security_tests.sh
./run_security_tests.sh
```

### Opción 3: Comando Directo
```bash
# Configurar variables de entorno
export ENV=test
export RUN_SECURITY_TESTS=true

# Ejecutar tests
go test -v ./module/gastos/test/ -run TestSecuritySuite -timeout 300s
```

## 📊 Estructura de los Tests

```
TestSecuritySuite/
├── TestAuthentication_Security/          # Tests de autenticación básica
├── TestRBAC_AdminPermissions/            # Tests de permisos de administrador
├── TestRBAC_RegularUserPermissions/      # Tests de permisos de usuario regular
├── TestRBAC_UnauthorizedUserPermissions/ # Tests de usuarios sin permisos
├── TestEndpoint_*_Security/              # Tests específicos por endpoint
├── TestSecurity_ContextInjection/        # Tests de inyección de contexto
└── TestSecurity_RateLimiting/            # Tests de rate limiting
```

## 🎯 Casos de Test Específicos

### 🔐 Autenticación
```go
// Sin token - 401 Unauthorized
GET /api/org/test-org/expenses

// Token inválido - 401 Unauthorized
Authorization: Bearer invalid-token

// Token válido - 200 OK
Authorization: Bearer <valid-jwt-token>
```

### 👑 Permisos de Administrador
```go
// Admin puede acceder a rutas administrativas
GET /api/org/test-org/admin/expenses/all
Authorization: Bearer <admin-token>
// Esperado: 200 OK

// Admin puede crear, leer, actualizar, eliminar
POST|GET|PUT|DELETE /api/org/test-org/expenses/*
Authorization: Bearer <admin-token>
// Esperado: 200/201 OK
```

### 👤 Permisos de Usuario Regular
```go
// Usuario puede ver gastos organizacionales
GET /api/org/test-org/expenses
Authorization: Bearer <user-token>
// Esperado: 200 OK

// Usuario NO puede acceder a rutas admin
GET /api/org/test-org/admin/expenses/all
Authorization: Bearer <user-token>
// Esperado: 403 Forbidden
```

### 🚫 Usuario Sin Permisos
```go
// Usuario sin permisos NO puede acceder a gastos
GET /api/org/test-org/expenses
Authorization: Bearer <unauthorized-token>
// Esperado: 403 Forbidden

// No puede crear gastos
POST /api/org/test-org/expenses
Authorization: Bearer <unauthorized-token>
// Esperado: 403 Forbidden
```

## 🏗️ Configuración del Suite

### Datos de Test Automáticos
- ✅ **Organización de prueba**: `test-security-org`
- ✅ **3 tipos de usuarios**: Admin, Regular, Sin permisos
- ✅ **Permisos completos**: Todos los permisos del módulo gastos
- ✅ **Roles específicos**: Admin, User, Unauthorized
- ✅ **Gastos de prueba**: Para operaciones CRUD

### Base de Datos de Test
- ✅ **Migración automática**: Todas las tablas necesarias
- ✅ **Datos de prueba**: Creación automática
- ✅ **Limpieza automática**: Al finalizar los tests
- ✅ **Aislamiento**: No afecta datos de producción

## 📈 Métricas y Resultados

### Cobertura de Seguridad
- ✅ **Autenticación**: 100% (Sin token, token inválido, token válido)
- ✅ **Autorización**: 100% (Admin, Usuario, Sin permisos)
- ✅ **Endpoints**: 100% (Todos los endpoints del módulo)
- ✅ **Casos Edge**: Context injection, rate limiting

### Tipos de Ataques Prevenidos
- ✅ **Acceso no autorizado**: Sin token o token inválido
- ✅ **Escalación de privilegios**: Usuario regular intentando acceso admin
- ✅ **Context injection**: Acceso a organizaciones no autorizadas
- ✅ **Brute force**: Rate limiting en requests múltiples

## 🔧 Dependencias Necesarias

```go
// Librería principal para tests HTTP reales
github.com/gavv/httpexpect/v2

// Framework de testing
github.com/stretchr/testify/suite
github.com/stretchr/testify/require

// Dependencias del proyecto
practicev2/config
practicev2/database
practicev2/module/authentication/*
practicev2/module/gastos/*
```

## 📝 Ejemplo de Salida

```
=== RUN   TestSecuritySuite
=== RUN   TestSecuritySuite/TestAuthentication_Security
=== RUN   TestSecuritySuite/TestAuthentication_Security/🔐_Sin_Token_-_Debe_Denegar_Acceso
=== RUN   TestSecuritySuite/TestAuthentication_Security/🔐_Token_Inválido_-_Debe_Denegar_Acceso
=== RUN   TestSecuritySuite/TestAuthentication_Security/🔐_Token_Válido_-_Debe_Permitir_Acceso
=== RUN   TestSecuritySuite/TestRBAC_AdminPermissions
=== RUN   TestSecuritySuite/TestRBAC_AdminPermissions/👑_Admin_-_Debe_Acceder_a_Rutas_Administrativas
...
--- PASS: TestSecuritySuite (45.23s)
PASS
```

## 🚨 Importante

### Requisitos Previos
1. **Base de datos activa**: MySQL/PostgreSQL corriendo
2. **Configuración válida**: Archivo `.env` con credenciales de test
3. **Puerto libre**: Puerto 8080 disponible para el servidor de test
4. **Variable de entorno**: `RUN_SECURITY_TESTS=true`

### Seguridad
- ❗ **Solo para testing**: No ejecutar contra base de datos de producción
- ❗ **Datos temporales**: Todos los datos se limpian automáticamente
- ❗ **Tokens de test**: Los JWTs son solo para pruebas

### Troubleshooting
- **Puerto ocupado**: Cambiar puerto en `s.baseURL`
- **Base de datos**: Verificar conexión en archivo `.env`
- **Permisos**: Asegurar que el usuario de DB puede crear/eliminar tablas

---

¡Los tests están listos para ejecutar! 🎉
