# Tests de Integración HTTP - Sistema de Autenticación

## 📋 Descripción

Esta carpeta contiene tests de integración completos que validan el sistema de autenticación a través de peticiones HTTP reales. Los tests están diseñados para probar todos los flujos críticos del sistema contra las rutas HTTP implementadas.

## 🏗️ Arquitectura de Tests

### Estructura de Archivos

```
http_tests/
├── http_integration_test_suite.go     # Suite base con utilidades comunes
├── company_lifecycle_http_test.go     # Test completo de empresa
├── school_lifecycle_http_test.go      # Test completo de colegio
├── ecommerce_flow_http_test.go        # Test de flujo e-commerce
├── security_permissions_http_test.go  # Tests de seguridad y permisos
├── run_http_integration_tests.sh      # Script de ejecución (Linux/Mac)
├── run_http_integration_tests.ps1     # Script de ejecución (Windows)
└── README.md                          # Esta documentación
```

### Configuración del Entorno

Los tests utilizan la configuración del archivo `.env` ubicado en `cmd/.env` y crean un entorno de testing aislado:

- **Base de datos**: Se usa `{DB_NAME}_test` para aislamiento
- **Puerto**: Puerto de test diferente (3031 por defecto)
- **Modo**: `ENV=testing`

## 🚀 Cómo Ejecutar los Tests

### Prerequisitos

1. **Go instalado** (versión 1.19+)
2. **Base de datos MySQL/MariaDB** ejecutándose
3. **Archivo `.env`** configurado correctamente
4. **Dependencias de Go** instaladas (`go mod tidy`)

### Opción 1: Script Automatizado (Recomendado)

#### Windows (PowerShell)
```powershell
# Ejecutar todos los tests
.\module\authentication\integration\http_tests\run_http_integration_tests.ps1

# Ejecutar test específico
.\module\authentication\integration\http_tests\run_http_integration_tests.ps1 -TestPattern "Company"

# Continuar aunque fallen tests
.\module\authentication\integration\http_tests\run_http_integration_tests.ps1 -ContinueOnError
```

#### Linux/Mac (Bash)
```bash
# Hacer ejecutable el script
chmod +x module/authentication/integration/http_tests/run_http_integration_tests.sh

# Ejecutar todos los tests
./module/authentication/integration/http_tests/run_http_integration_tests.sh
```

### Opción 2: Ejecución Manual con Go

```bash
# Desde el directorio raíz del proyecto
cd "c:\Users\ASUS\OneDrive\FREDY ponceca\carta ponceca\Documentos GRUPO PONCECA\PROYECTO 2025-1\BaseLogin-clean"

# Configurar variables de entorno
set ENV=testing
set TEST_PORT=3031
set DB_NAME=gastos_ia_test

# Ejecutar todos los tests
go test -v -timeout 10m ./module/authentication/integration/http_tests

# Ejecutar test específico
go test -v -timeout 10m ./module/authentication/integration/http_tests -run TestCompleteCompanyLifecycleHTTP
```

## 📊 Tests Implementados

### 1. Test de Ciclo de Vida de Empresa (`TestCompleteCompanyLifecycleHTTP`)

**Objetivo**: Validar el flujo completo de una empresa desde su creación hasta operaciones diarias.

**Fases del Test**:
1. ✅ **Creación de Empresa y CEO** - `POST /api/v1/organizations`
2. ✅ **Verificación de Perfil del CEO** - `GET /api/v1/auth/me`
3. ✅ **Creación de Roles Corporativos** - `POST /api/v1/org/{slug}/roles`
4. ✅ **Listado de Roles** - `GET /api/v1/org/{slug}/roles`
5. ✅ **Invitación de Empleados** - `POST /api/v1/org/{slug}/invitations`
6. ✅ **Aceptación de Invitaciones** - `POST /api/v1/invitations/accept`
7. ✅ **Verificación de Permisos** - `GET /api/v1/org/{slug}/users`
8. ✅ **Lista de Invitaciones** - `GET /api/v1/org/{slug}/invitations`
9. ✅ **Refresh Token** - `POST /api/v1/auth/refresh`
10. ✅ **Logout** - `POST /api/v1/auth/logout`
11. ✅ **Login Directo de Empleado** - `POST /api/v1/auth/login`

**Roles Creados**:
- `department_manager` (Nivel 90)
- `senior_developer` (Nivel 70)
- `accountant` (Nivel 60)

### 2. Test de Ciclo de Vida de Colegio (`TestCompleteSchoolLifecycleHTTP`)

**Objetivo**: Validar el ecosistema educativo completo con todos los tipos de usuario.

**Fases del Test**:
1. ✅ **Creación de Colegio y Director** - `POST /api/v1/organizations`
2. ✅ **Creación de Roles Educativos** - Múltiples `POST /api/v1/org/{slug}/roles`
3. ✅ **Invitación de Personal** - Coordinador, Profesores, Secretaria
4. ✅ **Aceptación de Personal** - `POST /api/v1/invitations/accept`
5. ✅ **Creación de Estudiantes y Padres** - Flujo completo de registro
6. ✅ **Login de Profesor** - `POST /api/v1/auth/login`
7. ✅ **Login de Estudiante** - Verificación de contexto académico
8. ✅ **Lista de Usuarios por Rol** - `GET /api/v1/org/{slug}/roles/{id}/users`
9. ✅ **Perfil Completo del Director** - `GET /api/v1/auth/me`
10. ✅ **Verificación de Auditoría** - `GET /api/v1/org/{slug}/audits`

**Roles Educativos**:
- `academic_coordinator` (Nivel 90)
- `teacher` (Nivel 60)
- `student` (Nivel 20)
- `parent` (Nivel 30)
- `secretary` (Nivel 50)

### 3. Test de Flujo E-Commerce (`TestCompleteECommerceFlowHTTP`)

**Objetivo**: Validar todo el flujo de comercio electrónico desde invitado hasta cliente.

**Fases del Test**:
1. ✅ **Sesión de Invitado** - `POST /api/v1/guests/session`
2. ✅ **Verificación de Sesión** - `GET /api/v1/guests/session`
3. ✅ **Carrito de Invitado** - `PUT /api/v1/guests/session/cart`
4. ✅ **Registro Directo de Cliente** - `POST /api/v1/customers/register`
5. ✅ **Perfil de Cliente** - `GET /api/v1/customers/profile`
6. ✅ **Direcciones de Envío** - `POST|GET|PUT /api/v1/customers/addresses`
7. ✅ **Listado de Direcciones** - Verificación de direcciones múltiples
8. ✅ **Conversión de Invitado** - `POST /api/v1/customers/register`
9. ✅ **Login de Cliente** - `POST /api/v1/auth/login`
10. ✅ **Gestión de Preferencias** - `GET|PUT /api/v1/customers/preferences`
11. ✅ **Actualización de Dirección** - `PUT /api/v1/customers/addresses/{id}`
12. ✅ **Refresh Token Cliente** - `POST /api/v1/auth/refresh`
13. ✅ **Logout Cliente** - `POST /api/v1/auth/logout`

### 4. Test de Seguridad y Permisos (`TestSecurityAndPermissionsHTTP`)

**Objetivo**: Validar todos los aspectos de seguridad del sistema.

**Fases del Test**:
1. ✅ **Acceso No Autorizado** - Intentos sin token
2. ✅ **Tokens Inválidos** - Diferentes tipos de tokens malformados
3. ✅ **Credenciales Incorrectas** - Múltiples escenarios de login fallido
4. ✅ **Rate Limiting** - Intentos masivos de login
5. ✅ **Permisos Organizacionales** - Acceso cross-organizacional
6. ✅ **Headers de Seguridad** - Verificación de headers HTTP
7. ✅ **Refresh Token Security** - Tokens de refresh inválidos
8. ✅ **Sesiones Múltiples** - Gestión de múltiples sesiones
9. ✅ **Logout e Invalidación** - Invalidación correcta de tokens
10. ✅ **Expiración de Token** - Estructura y validación de tokens
11. ✅ **Invitaciones Expiradas** - Sistema de expiración

## 🔧 Configuración Avanzada

### Variables de Entorno para Testing

```bash
# Configuración básica
ENV=testing
TEST_PORT=3031

# Base de datos de testing
DB_NAME=gastos_ia_test
DB_HOST=127.0.0.1
DB_PORT=3308
DB_USER=root
DB_PASSWORD=toor

# JWT para testing
JWT_KEY=test-jwt-key-for-testing-only
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_HOURS=168

# Configuración de rate limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW_MINUTES=15
```

### Personalización de Tests

#### Modificar Timeout
```powershell
# PowerShell - Timeout de 20 minutos
.\run_http_integration_tests.ps1 -Timeout 1200
```

#### Ejecutar Tests Específicos
```powershell
# Solo tests de empresa
.\run_http_integration_tests.ps1 -TestPattern "Company"

# Solo tests de seguridad
.\run_http_integration_tests.ps1 -TestPattern "Security"
```

## 📝 Interpretación de Resultados

### Salida Exitosa
```
[2025-01-01 10:00:00] ✅ Ciclo de Vida de Empresa - EXITOSO
[2025-01-01 10:05:00] ✅ Ciclo de Vida de Colegio - EXITOSO
[2025-01-01 10:10:00] ✅ Flujo de E-Commerce - EXITOSO
[2025-01-01 10:15:00] ✅ Seguridad y Permisos - EXITOSO

================================================================
📊 RESUMEN DE RESULTADOS
================================================================
Total de tests ejecutados: 4
Tests exitosos: 4
Tests fallidos: 0
Tasa de éxito: 100%

🎉 TODOS LOS TESTS DE INTEGRACIÓN HTTP EXITOSOS
✅ El sistema de autenticación está funcionando correctamente
```

### Interpretación de Errores

#### Error de Conexión a BD
```
[ERROR] No se puede conectar a la base de datos en 127.0.0.1:3308
```
**Solución**: Verificar que MySQL/MariaDB esté ejecutándose en el puerto correcto.

#### Error de Token Inválido
```
[ERROR] Token inválido o expirado
```
**Solución**: Verificar configuración JWT en `.env`.

#### Error de Permisos
```
[ERROR] Acceso denegado - Sin permisos suficientes
```
**Solución**: Verificar implementación de middleware de permisos.

## 🛠️ Debugging y Troubleshooting

### Logs Detallados
```bash
# Ejecutar con logs verbose de Go
go test -v -timeout 10m ./module/authentication/integration/http_tests -run TestCompleteCompanyLifecycleHTTP 2>&1 | tee test_output.log
```

### Verificación de Rutas
Los tests validan que estas rutas estén implementadas correctamente:

#### Rutas Públicas
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/organizations`
- `POST /api/v1/customers/register`
- `POST /api/v1/guests/session`
- `POST /api/v1/invitations/accept`

#### Rutas Protegidas
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`
- `GET /api/v1/customers/profile`
- `GET|POST /api/v1/customers/addresses`
- `GET|POST /api/v1/org/{slug}/roles`
- `GET|POST /api/v1/org/{slug}/invitations`
- `GET /api/v1/org/{slug}/users`
- `GET /api/v1/org/{slug}/audits`

### Verificación de Base de Datos

```sql
-- Verificar que las tablas necesarias existen
SHOW TABLES LIKE '%identities%';
SHOW TABLES LIKE '%organizations%';
SHOW TABLES LIKE '%roles%';
SHOW TABLES LIKE '%invitations%';
SHOW TABLES LIKE '%customer_profiles%';
SHOW TABLES LIKE '%guest_sessions%';
```

## 🔄 Integración Continua

### GitHub Actions (Ejemplo)
```yaml
name: HTTP Integration Tests
on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: toor
          MYSQL_DATABASE: gastos_ia_test
        ports:
          - 3308:3306

    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.19

    - name: Run HTTP Integration Tests
      run: |
        chmod +x module/authentication/integration/http_tests/run_http_integration_tests.sh
        ./module/authentication/integration/http_tests/run_http_integration_tests.sh
```

## 📈 Métricas y Cobertura

### Cobertura de Funcionalidades
- ✅ **Autenticación Básica**: Login, Logout, Refresh
- ✅ **Organizaciones Multi-tenant**: Empresas y Colegios
- ✅ **Roles y Permisos**: RBAC completo
- ✅ **Invitaciones**: Flujo completo de invitar usuarios
- ✅ **E-commerce**: Clientes, direcciones, preferencias
- ✅ **Sesiones de Invitado**: Carrito temporal
- ✅ **Seguridad**: Validación de tokens, permisos, rate limiting
- ✅ **Auditoría**: Logs de todas las acciones

### Cobertura de Rutas HTTP
Los tests cubren **100%** de las rutas críticas implementadas en `module/authentication/routes.go`.

## 🎯 Objetivos y Beneficios

### ✅ Validación Completa
- Verifica que todas las rutas funcionan correctamente
- Prueba flujos de usuario reales end-to-end
- Valida integración entre todos los módulos

### ✅ Detección Temprana de Errores
- Encuentra problemas de integración antes de producción
- Valida cambios en la API automáticamente
- Detecta regresiones en funcionalidad existente

### ✅ Documentación Viva
- Los tests sirven como documentación de la API
- Muestran ejemplos reales de uso de cada endpoint
- Mantienen la documentación actualizada automáticamente

### ✅ Confianza en Despliegues
- Garantiza que el sistema funciona correctamente
- Reduce riesgo de errores en producción
- Facilita refactoring seguro

---

## 📞 Soporte

Si encuentras problemas con los tests:

1. **Verificar configuración**: Revisar archivo `.env`
2. **Verificar base de datos**: Confirmar que está ejecutándose
3. **Verificar dependencias**: Ejecutar `go mod tidy`
4. **Revisar logs**: Analizar mensajes de error detallados
5. **Ejecutar tests individuales**: Aislar el problema

Los tests están diseñados para ser robustos y proporcionar mensajes de error claros para facilitar el debugging.
