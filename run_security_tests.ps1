# =============================================================================
# SCRIPT PARA EJECUTAR TESTS DE SEGURIDAD - WINDOWS POWERSHELL
# =============================================================================

Write-Host "🛡️  INICIANDO TESTS DE SEGURIDAD COMPLETOS" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan

# Configurar variables de entorno para tests
$env:ENV = "test"
$env:RUN_SECURITY_TESTS = "true"

# Configurar base de datos de test
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "3308"
$env:DB_USER = "root"
$env:DB_PASSWORD = ""
$env:DB_NAME = "gastos_ia_test"

Write-Host ""
Write-Host "📋 Configuración de Test:" -ForegroundColor Yellow
Write-Host "   - Entorno: $($env:ENV)" -ForegroundColor White
Write-Host "   - Base de datos: $($env:DB_NAME)" -ForegroundColor White
Write-Host "   - Host: $($env:DB_HOST):$($env:DB_PORT)" -ForegroundColor White
Write-Host ""

Write-Host "🧹 Preparando entorno de test..." -ForegroundColor Yellow

Write-Host "🚀 Ejecutando tests de seguridad..." -ForegroundColor Green
Write-Host "   - Tests de Autenticación" -ForegroundColor White
Write-Host "   - Tests de RBAC" -ForegroundColor White
Write-Host "   - Tests de Endpoints" -ForegroundColor White
Write-Host "   - Tests de Inyección de Contexto" -ForegroundColor White
Write-Host "   - Tests de Rate Limiting" -ForegroundColor White
Write-Host ""

# Ejecutar tests con verbose output
try {
    go test -v ./module/gastos/test/ -run TestSecuritySuite -timeout 300s
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✅ Tests de seguridad completados exitosamente" -ForegroundColor Green
        Write-Host "===============================================" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "❌ Algunos tests de seguridad fallaron" -ForegroundColor Red
        Write-Host "Revisa los logs arriba para más detalles" -ForegroundColor Red
        Write-Host "===============================================" -ForegroundColor Red
    }
} catch {
    Write-Host ""
    Write-Host "❌ Error ejecutando tests de seguridad: $_" -ForegroundColor Red
    Write-Host "===============================================" -ForegroundColor Red
}
