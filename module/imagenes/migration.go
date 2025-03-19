package imagenes

import (
	"practicev2/registry"

	"gorm.io/gorm"
)

// Migration order
const (
	ImagesMigrationOrder = 20
)

// MigrateImages creates the images table in the database
func MigrateImages(db *gorm.DB) error {
	return db.AutoMigrate(&Image{})
}

func init() {
	// Register migration with the registry
	registry.RegisterMigration(ImagesMigrationOrder, MigrateImages)
}
