// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/registry/migrations.go
package registry

import (
	"sort"

	"gorm.io/gorm"
)

type MigrationFunc func(db *gorm.DB) error

type prioritizedMigration struct {
	order   int
	migrate MigrationFunc
}

var migrations []prioritizedMigration

// RegisterMigration registra una función de migración con un orden específico.
func RegisterMigration(order int, migrate MigrationFunc) {
	migrations = append(migrations, prioritizedMigration{order: order, migrate: migrate})
}

// RunMigrations ordena y ejecuta todas las migraciones registradas.
func RunMigrations(db *gorm.DB) {
	sort.SliceStable(migrations, func(i, j int) bool {
		return migrations[i].order < migrations[j].order
	})
	for _, m := range migrations {
		if err := m.migrate(db); err != nil {
			panic("Migration failed: " + err.Error())
		}
	}
}
