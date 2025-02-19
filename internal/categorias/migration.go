// // Language: go
// filepath: /c:/Users/ASUS/OneDrive/FREDY ponceca/carta ponceca/Documentos GRUPO PONCECA/PROYECTO 2025-1/febrero-practice-go/internal/categorias/migration.go
package categorias

import (
	"practicev2/registry"

	"gorm.io/gorm"
)

// Migrate ejecuta la migración del modelo Categoria.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Categoria{})
}

func init() {
	// Se registra la migración con orden 1 (mayor prioridad).
	registry.RegisterMigration(1, Migrate)
}
