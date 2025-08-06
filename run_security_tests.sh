#!/bin/bash

# =============================================================================
# SCRIPT PARA EJECUTAR TESTS DE SEGURIDAD
# =============================================================================

echo "🛡️  INICIANDO TESTS DE SEGURIDAD COMPLETOS"
echo "=============================================="

# Configurar variables de entorno para tests
export ENV=test
export RUN_SECURITY_TESTS=true

# Configurar base de datos de test
export DB_HOST=127.0.0.1
export DB_PORT=3308
export DB_USER=root
export DB_PASSWORD=
export DB_NAME=gastos_ia_test

echo "📋 Configuración de Test:"
echo "   - Entorno: $ENV"
echo "   - Base de datos: $DB_NAME"
echo "   - Host: $DB_HOST:$DB_PORT"
echo ""

echo "🧹 Limpiando base de datos de test..."
# Opcional: limpiar la base de datos antes de empezar
# mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -e "DROP DATABASE IF EXISTS $DB_NAME; CREATE DATABASE $DB_NAME;"

echo "🚀 Ejecutando tests de seguridad..."
echo "   - Tests de Autenticación"
echo "   - Tests de RBAC"
echo "   - Tests de Endpoints"
echo "   - Tests de Inyección de Contexto"
echo "   - Tests de Rate Limiting"
echo ""

# Ejecutar tests con verbose output
go test -v ./module/gastos/test/ -run TestSecuritySuite -timeout 300s

echo ""
echo "✅ Tests de seguridad completados"
echo "=============================================="
