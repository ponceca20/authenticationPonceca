package gastos

import (
	"fmt"
	"practicev2/module/gastos/model"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// migrateExpenseModels executes the database migrations for the gastos module.
func migrateExpenseModels(db *gorm.DB) error {
	// Configurar logger silencioso para las migraciones
	originalLogger := db.Logger
	defer func() {
		db.Logger = originalLogger // Restaurar logger original
	}()

	// Usar logger silencioso solo para migraciones
	db.Logger = db.Logger.LogMode(logger.Silent)

	// GORM's AutoMigrate will create tables, missing foreign keys, constraints, columns, and indexes.
	// It will NOT delete unused columns to protect your data.
	err := db.AutoMigrate(
		&model.Expense{},
	)

	if err != nil {
		return fmt.Errorf("gastos module migration failed: %w", err)
	}

	// Create additional composite indexes for better performance and data integrity
	err = createAdditionalIndexes(db)
	if err != nil {
		return fmt.Errorf("failed to create additional indexes: %w", err)
	}

	// Sembrar datos iniciales si es necesario
	err = seedExpenseData(db)
	if err != nil {
		return fmt.Errorf("failed to seed expense data: %w", err)
	}

	fmt.Println("✅ Gastos module migrated successfully.")
	return nil
}

// createAdditionalIndexes creates composite indexes for better performance and data integrity
func createAdditionalIndexes(db *gorm.DB) error {
	// Detectar el dialecto y configurar índices apropiados
	var indexes []indexDefinition

	switch db.Dialector.Name() {
	case "mysql":
		indexes = []indexDefinition{
			{
				Name:        "idx_expense_title",
				Table:       "expense",
				Query:       "CREATE INDEX idx_expense_title ON expense (title)",
				Description: "Expense title index for search optimization",
			},
			{
				Name:        "idx_expense_amount",
				Table:       "expense",
				Query:       "CREATE INDEX idx_expense_amount ON expense (amount)",
				Description: "Expense amount index for range queries",
			},
			{
				Name:        "idx_expense_created_at",
				Table:       "expense",
				Query:       "CREATE INDEX idx_expense_created_at ON expense (created_at)",
				Description: "Expense creation date index",
			},
		}
	case "sqlite":
		indexes = []indexDefinition{
			{
				Name:        "idx_expense_title",
				Table:       "expense",
				Query:       "CREATE INDEX IF NOT EXISTS idx_expense_title ON expense (title)",
				Description: "Expense title index for search optimization",
			},
			{
				Name:        "idx_expense_amount",
				Table:       "expense",
				Query:       "CREATE INDEX IF NOT EXISTS idx_expense_amount ON expense (amount)",
				Description: "Expense amount index for range queries",
			},
			{
				Name:        "idx_expense_created_at",
				Table:       "expense",
				Query:       "CREATE INDEX IF NOT EXISTS idx_expense_created_at ON expense (created_at)",
				Description: "Expense creation date index",
			},
		}
	default:
		// Para otros dialectos, usar sintaxis básica
		indexes = []indexDefinition{
			{
				Name:        "idx_expense_title",
				Table:       "expense",
				Query:       "CREATE INDEX idx_expense_title ON expense (title)",
				Description: "Expense title index for search optimization",
			},
		}
	}

	// Crear índices de forma silenciosa
	createdCount := 0
	skippedCount := 0

	for _, idx := range indexes {
		if err := createIndexIfNotExists(db, idx, &createdCount, &skippedCount); err != nil {
			return fmt.Errorf("failed to create %s: %w", idx.Description, err)
		}
	}

	// Solo mostrar un resumen
	if createdCount > 0 || skippedCount > 0 {
		fmt.Printf("📊 Gastos Índices: %d creados, %d existentes\n", createdCount, skippedCount)
	}

	return nil
}

// seedExpenseData siembra datos iniciales si es necesario
func seedExpenseData(db *gorm.DB) error {
	// Verificar si ya existen datos de ejemplo
	var count int64

	// Ejecutar con logger silencioso
	originalLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Silent)
	db.Model(&model.Expense{}).Count(&count)
	db.Logger = originalLogger

	// Si no hay datos, crear algunos ejemplos
	if count == 0 {
		sampleExpenses := []model.Expense{
			{
				Title:       "Compra de suministros de oficina",
				Description: "Papel, bolígrafos y carpetas para el mes",
				Amount:      150.00,
			},
			{
				Title:       "Almuerzo de trabajo",
				Description: "Reunión con cliente potencial",
				Amount:      75.50,
			},
			{
				Title:       "Software de productividad",
				Description: "Licencia anual de software",
				Amount:      299.99,
			},
		}

		// Crear datos de ejemplo con logger silencioso
		db.Logger = db.Logger.LogMode(logger.Silent)
		for _, expense := range sampleExpenses {
			if err := db.Create(&expense).Error; err != nil {
				db.Logger = originalLogger
				return fmt.Errorf("failed to create sample expense: %w", err)
			}
		}
		db.Logger = originalLogger

		fmt.Printf("🌱 Gastos: %d registros de ejemplo creados\n", len(sampleExpenses))
	}

	return nil
}

// indexDefinition represents an index to be created
type indexDefinition struct {
	Name        string
	Table       string
	Query       string
	Description string
}

// createIndexIfNotExists creates an index only if it doesn't already exist
func createIndexIfNotExists(db *gorm.DB, idx indexDefinition, createdCount, skippedCount *int) error {
	// Check if index already exists
	exists, err := indexExists(db, idx.Table, idx.Name)
	if err != nil {
		return fmt.Errorf("error checking if index %s exists: %w", idx.Name, err)
	}

	if exists {
		*skippedCount++
		return nil
	}

	// Create the index with silent logger
	originalLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Silent)
	err = db.Exec(idx.Query).Error
	db.Logger = originalLogger

	if err != nil {
		// Check if the error is due to duplicate key name (index already exists)
		if isDuplicateIndexError(err) {
			*skippedCount++
			return nil
		}
		return fmt.Errorf("error creating index %s: %w", idx.Name, err)
	}

	*createdCount++
	return nil
}

// indexExists checks if an index exists on a table
func indexExists(db *gorm.DB, tableName, indexName string) (bool, error) {
	var count int64

	// Detectar el dialecto de la base de datos
	switch db.Dialector.Name() {
	case "mysql":
		// Query para MySQL
		query := `
			SELECT COUNT(*) 
			FROM information_schema.statistics 
			WHERE table_schema = DATABASE() 
			AND table_name = ? 
			AND index_name = ?
		`
		originalLogger := db.Logger
		db.Logger = db.Logger.LogMode(logger.Silent)
		err := db.Raw(query, tableName, indexName).Scan(&count).Error
		db.Logger = originalLogger
		return count > 0, err

	case "sqlite":
		// Para SQLite, usar PRAGMA index_list
		originalLogger := db.Logger
		db.Logger = db.Logger.LogMode(logger.Silent)

		// Ejecutar consulta PRAGMA sin parámetros
		rows, err := db.Raw("PRAGMA index_list(" + tableName + ")").Rows()
		db.Logger = originalLogger

		if err != nil {
			return false, err
		}
		defer rows.Close()

		// Buscar el índice en los resultados
		for rows.Next() {
			var seq int
			var name string
			var unique int
			var origin string
			var partial int

			if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
				continue
			}

			if name == indexName {
				return true, nil
			}
		}
		return false, nil

	default:
		// Para otros dialectos, asumir que no existe para forzar creación
		return false, nil
	}
}

// isDuplicateIndexError checks if the error is due to a duplicate index
func isDuplicateIndexError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// MySQL error codes for duplicate index/key
	duplicateErrors := []string{
		"Error 1061", // Duplicate key name
		"Error 1062", // Duplicate entry
		"Duplicate key name",
		"duplicate key name",
	}

	for _, dupErr := range duplicateErrors {
		if strings.Contains(errStr, dupErr) {
			return true
		}
	}

	return false
}
