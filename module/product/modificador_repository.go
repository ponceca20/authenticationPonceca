package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ModificadorRepository define métodos para operar con Modificador
type ModificadorRepository interface {
	GetAllModificadores(empresaID uint64) ([]Modificador, error)
	GetModificador(id uint64, empresaID uint64) (Modificador, error)
	CreateModificador(m *Modificador) error
	UpdateModificador(id uint64, empresaID uint64, m *Modificador) (Modificador, error)
	DeleteModificador(id uint64, empresaID uint64) error
	SeedModificadores(mods []Modificador) error
	PatchModificador(id uint64, empresaID uint64, fields map[string]interface{}) (Modificador, error)
}

type modificadorRepo struct {
	DB *gorm.DB
}

func NewModificadorRepository() ModificadorRepository {
	return &modificadorRepo{DB: database.DBconn}
}

func (r *modificadorRepo) GetAllModificadores(empresaID uint64) ([]Modificador, error) {
	var mods []Modificador
	err := r.DB.Where("empresa_id = ?", empresaID).Preload("Opciones").Find(&mods).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar modificadores: %w", err)
	}
	return mods, nil
}

func (r *modificadorRepo) GetModificador(id uint64, empresaID uint64) (Modificador, error) {
	var mod Modificador
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Preload("Opciones").First(&mod).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mod, fmt.Errorf("modificador con id %d no encontrado", id)
		}
		return mod, fmt.Errorf("error al recuperar modificador con id %d: %w", id, err)
	}
	return mod, nil
}

func (r *modificadorRepo) CreateModificador(m *Modificador) error {
	err := r.DB.Create(m).Error
	if err != nil {
		return fmt.Errorf("error al crear modificador: %w", err)
	}
	return nil
}

func (r *modificadorRepo) UpdateModificador(id uint64, empresaID uint64, updatedM *Modificador) (Modificador, error) {
	var mod Modificador
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&mod).Error; err != nil {
			return err
		}
		return tx.Model(&mod).Updates(updatedM).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mod, fmt.Errorf("modificador con id %d no encontrado", id)
		}
		return mod, fmt.Errorf("error al actualizar modificador con id %d: %w", id, err)
	}
	
	// Recargar el modificador con sus opciones
	err = r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Preload("Opciones").First(&mod).Error
	if err != nil {
		return mod, fmt.Errorf("error al recargar modificador con id %d: %w", id, err)
	}
	return mod, nil
}

func (r *modificadorRepo) DeleteModificador(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Modificador{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("modificador con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar modificador: %w", err)
	}
	return nil
}

func (r *modificadorRepo) SeedModificadores(mods []Modificador) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&mods).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar modificadores: %w", err)
	}
	return nil
}

func (r *modificadorRepo) PatchModificador(id uint64, empresaID uint64, fields map[string]interface{}) (Modificador, error) {
	var mod Modificador
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&mod).Error; err != nil {
			return err
		}
		return tx.Model(&mod).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mod, fmt.Errorf("modificador con id %d no encontrado", id)
		}
		return mod, fmt.Errorf("error al actualizar parcialmente modificador con id %d: %w", id, err)
	}
	
	// Recargar el modificador con sus opciones
	err = r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Preload("Opciones").First(&mod).Error
	if err != nil {
		return mod, fmt.Errorf("error al recargar modificador con id %d: %w", id, err)
	}
	return mod, nil
}
