package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ModificadorOpcionRepository define métodos para operar con ModificadorOpcion
type ModificadorOpcionRepository interface {
	GetAllModificadorOpciones(modificadorID uint64) ([]ModificadorOpcion, error)
	GetModificadorOpcion(id uint64) (ModificadorOpcion, error)
	CreateModificadorOpcion(mo *ModificadorOpcion) error
	UpdateModificadorOpcion(id uint64, mo *ModificadorOpcion) (ModificadorOpcion, error)
	DeleteModificadorOpcion(id uint64) error
	SeedModificadorOpciones(mos []ModificadorOpcion) error
	PatchModificadorOpcion(id uint64, fields map[string]interface{}) (ModificadorOpcion, error)
}

type modificadorOpcionRepo struct {
	DB *gorm.DB
}

func NewModificadorOpcionRepository() ModificadorOpcionRepository {
	return &modificadorOpcionRepo{DB: database.DBconn}
}

func (r *modificadorOpcionRepo) GetAllModificadorOpciones(modificadorID uint64) ([]ModificadorOpcion, error) {
	var mos []ModificadorOpcion
	err := r.DB.Where("modificador_id = ?", modificadorID).Find(&mos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar opciones de modificador: %w", err)
	}
	return mos, nil
}

func (r *modificadorOpcionRepo) GetModificadorOpcion(id uint64) (ModificadorOpcion, error) {
	var mo ModificadorOpcion
	err := r.DB.Where("id = ?", id).First(&mo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mo, fmt.Errorf("opción de modificador con id %d no encontrada", id)
		}
		return mo, fmt.Errorf("error al recuperar opción de modificador con id %d: %w", id, err)
	}
	return mo, nil
}

func (r *modificadorOpcionRepo) CreateModificadorOpcion(mo *ModificadorOpcion) error {
	// Verificar si el modificador padre existe
	var count int64
	if err := r.DB.Model(&Modificador{}).Where("id = ?", mo.ModificadorID).Count(&count).Error; err != nil {
		return fmt.Errorf("error al verificar modificador: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("el modificador con ID %d no existe", mo.ModificadorID)
	}
	if err := r.DB.Create(mo).Error; err != nil {
		return fmt.Errorf("error al crear opción de modificador: %w", err)
	}
	return nil
}

func (r *modificadorOpcionRepo) UpdateModificadorOpcion(id uint64, updatedMO *ModificadorOpcion) (ModificadorOpcion, error) {
	var mo ModificadorOpcion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&mo).Error; err != nil {
			return err
		}
		// Validar modificador padre si se cambia
		if updatedMO.ModificadorID != mo.ModificadorID {
			var count int64
			if err := tx.Model(&Modificador{}).Where("id = ?", updatedMO.ModificadorID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar modificador: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("el modificador con ID %d no existe", updatedMO.ModificadorID)
			}
		}
		return tx.Model(&mo).Updates(updatedMO).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mo, fmt.Errorf("opción de modificador con id %d no encontrada", id)
		}
		return mo, fmt.Errorf("error al actualizar opción de modificador con id %d: %w", id, err)
	}
	return mo, nil
}

func (r *modificadorOpcionRepo) DeleteModificadorOpcion(id uint64) error {
	result := r.DB.Delete(&ModificadorOpcion{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar opción de modificador: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("opción de modificador con id %d no encontrada", id)
	}
	return nil
}

func (r *modificadorOpcionRepo) SeedModificadorOpciones(mos []ModificadorOpcion) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&mos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar opciones de modificador: %w", err)
	}
	return nil
}

func (r *modificadorOpcionRepo) PatchModificadorOpcion(id uint64, fields map[string]interface{}) (ModificadorOpcion, error) {
	var mo ModificadorOpcion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&mo).Error; err != nil {
			return err
		}
		// Validar modificador padre si se cambia
		if modificadorID, ok := fields["modificador_id"].(uint64); ok && modificadorID != mo.ModificadorID {
			var count int64
			if err := tx.Model(&Modificador{}).Where("id = ?", modificadorID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar modificador: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("el modificador con ID %d no existe", modificadorID)
			}
		}
		return tx.Model(&mo).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mo, fmt.Errorf("opción de modificador con id %d no encontrada", id)
		}
		return mo, fmt.Errorf("error al actualizar parcialmente opción de modificador con id %d: %w", id, err)
	}
	return mo, nil
}
