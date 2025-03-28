package auth

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ModuloRepository define la interfaz para operaciones de persistencia en Modulo
type ModuloRepository interface {
	GetAllModulos() ([]Modulo, error)
	GetModulo(id uint64) (Modulo, error)
	CreateModulo(modulo *Modulo) error
	UpdateModulo(id uint64, updatedModulo *Modulo) (Modulo, error)
	DeleteModulo(id uint64) error
	SeedModulos(modulos []Modulo) error
	PatchModulo(id uint64, fields map[string]interface{}) (Modulo, error)
}

type moduloRepo struct {
	DB *gorm.DB
}

// NewModuloRepository crea una nueva instancia del repositorio
func NewModuloRepository() ModuloRepository {
	return &moduloRepo{DB: database.DBconn}
}

func (r *moduloRepo) GetAllModulos() ([]Modulo, error) {
	var modulos []Modulo
	err := r.DB.Find(&modulos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar modulos: %w", err)
	}
	return modulos, nil
}

func (r *moduloRepo) GetModulo(id uint64) (Modulo, error) {
	var modulo Modulo
	err := r.DB.First(&modulo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modulo, fmt.Errorf("módulo con id %d no encontrado", id)
		}
		return modulo, fmt.Errorf("error al recuperar módulo con id %d: %w", id, err)
	}
	return modulo, nil
}

func (r *moduloRepo) CreateModulo(modulo *Modulo) error {
	err := r.DB.Create(modulo).Error
	if err != nil {
		return fmt.Errorf("error al crear módulo: %w", err)
	}
	return nil
}

func (r *moduloRepo) UpdateModulo(id uint64, updatedModulo *Modulo) (Modulo, error) {
	var modulo Modulo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&modulo, id).Error; err != nil {
			return err
		}
		return tx.Model(&modulo).Updates(updatedModulo).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modulo, fmt.Errorf("módulo con id %d no encontrado", id)
		}
		return modulo, fmt.Errorf("error al actualizar módulo con id %d: %w", id, err)
	}
	return modulo, nil
}

func (r *moduloRepo) DeleteModulo(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Modulo{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("módulo con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar módulo: %w", err)
	}
	return nil
}

func (r *moduloRepo) SeedModulos(modulos []Modulo) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&modulos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar modulos: %w", err)
	}
	return nil
}

func (r *moduloRepo) PatchModulo(id uint64, fields map[string]interface{}) (Modulo, error) {
	var modulo Modulo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&modulo, id).Error; err != nil {
			return err
		}
		return tx.Model(&modulo).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modulo, fmt.Errorf("módulo con id %d no encontrado", id)
		}
		return modulo, fmt.Errorf("error al actualizar parcialmente módulo con id %d: %w", id, err)
	}
	return modulo, nil
}
