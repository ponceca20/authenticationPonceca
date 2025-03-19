package users

import (
	"practicev2/registry"

	"gorm.io/gorm"
)

// Migrate ejecuta la migración de las tablas de usuarios en el orden correcto.
func Migrate(db *gorm.DB) error {
	// Migrar LoginAttempt (independiente)
	if err := db.AutoMigrate(&LoginAttempt{}); err != nil {
		return err
	}
	// Primero se migra Persona (no tiene dependencias)
	if err := db.AutoMigrate(&Persona{}); err != nil {
		return err
	}
	// Migrar Usuario, que depende de Persona
	if err := db.AutoMigrate(&Usuario{}); err != nil {
		return err
	}
	// Migrar Rol (se recomienda que exista previamente la tabla Empresa si aplica)
	if err := db.AutoMigrate(&Rol{}); err != nil {
		return err
	}
	// Migrar Modulo (independiente)
	if err := db.AutoMigrate(&Modulo{}); err != nil {
		return err
	}
	// Migrar Sesion, que depende de Usuario
	if err := db.AutoMigrate(&Sesion{}); err != nil {
		return err
	}
	// Migrar UsuarioEmpresa, que depende de Usuario y Rol
	if err := db.AutoMigrate(&UsuarioEmpresa{}); err != nil {
		return err
	}
	// Migrar RolModulo, que depende de Rol y Modulo
	if err := db.AutoMigrate(&RolModulo{}); err != nil {
		return err
	}

	return nil
}

func init() {
	// Se registra la migración con orden 2 (posterior a categorías).
	registry.RegisterMigration(2, Migrate)
}
