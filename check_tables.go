package main

import (
	"fmt"
	"log"

	"practicev2/config"
	"practicev2/database"
)

func main() {
	fmt.Println("🔍 Verificando tablas de la base de datos...")

	// Inicializar configuración
	config.Init()

	// Conectar a la base de datos
	database.ConnectDatabase()
	db := database.DBconn

	// Obtener SQL DB para operaciones raw
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Error obteniendo SQL DB: %v", err)
	}
	defer sqlDB.Close()

	fmt.Printf("✅ Conectado a la base de datos: %s\n", config.GlobalConfig.Database.Name)

	// Listar todas las tablas
	fmt.Println("\n📋 Listando todas las tablas:")
	rows, err := sqlDB.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("❌ Error listando tablas: %v", err)
	}
	defer rows.Close()

	var tableCount int
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("⚠️ Error escaneando tabla: %v", err)
			continue
		}
		fmt.Printf("  📋 %s\n", tableName)
		tableCount++
	}

	if tableCount == 0 {
		fmt.Println("⚠️ No se encontraron tablas en la base de datos")
	} else {
		fmt.Printf("\n✅ Total de tablas encontradas: %d\n", tableCount)
	}

	// Verificar tablas específicas que deberían existir
	expectedTables := []string{
		"organizations",
		"roles",
		"users",
		"profiles",
		"departments",
		"invitations",
		"audit_logs",
		"user_organizations",
		"user_roles",
	}

	fmt.Println("\n🔍 Verificando tablas esperadas:")
	for _, table := range expectedTables {
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s')",
			config.GlobalConfig.Database.Name, table)

		err := sqlDB.QueryRow(query).Scan(&exists)
		if err != nil {
			fmt.Printf("  ❌ Error verificando tabla '%s': %v\n", table, err)
			continue
		}

		if exists {
			fmt.Printf("  ✅ %s\n", table)
		} else {
			fmt.Printf("  ❌ %s (NO EXISTE)\n", table)
		}
	}

	fmt.Println("\n🏁 Verificación completada")
}
