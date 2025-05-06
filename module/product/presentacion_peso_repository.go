package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PresentacionPesoRepository define la interfaz para operaciones de persistencia de PresentacionPeso
type PresentacionPesoRepository interface {
	GetAllPresentacionesPeso(presentacionID uint64) ([]PresentacionPeso, error)
	GetPresentacionPeso(id uint64) (PresentacionPeso, error)
	CreatePresentacionPeso(peso *PresentacionPeso) error
	UpdatePresentacionPeso(id uint64, peso *PresentacionPeso) (PresentacionPeso, error)
	DeletePresentacionPeso(id uint64) error
	SeedPresentacionesPeso(pesos []PresentacionPeso) error
	PatchPresentacionPeso(id uint64, fields map[string]interface{}) (PresentacionPeso, error)
}

type presentacionPesoRepo struct {
	DB *gorm.DB
}

// NewPresentacionPesoRepository crea una instancia del repositorio
func NewPresentacionPesoRepository() PresentacionPesoRepository {
	return &presentacionPesoRepo{DB: database.DBconn}
}

// GetAllPresentacionesPeso obtiene todas las configuraciones de peso para una presentación
func (r *presentacionPesoRepo) GetAllPresentacionesPeso(presentacionID uint64) ([]PresentacionPeso, error) {
	var pesos []PresentacionPeso
	err := r.DB.Where("presentacion_id = ?", presentacionID).Find(&pesos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar configuraciones de peso: %w", err)
	}
	return pesos, nil
}

// GetPresentacionPeso obtiene una configuración de peso específica
func (r *presentacionPesoRepo) GetPresentacionPeso(id uint64) (PresentacionPeso, error) {
	var peso PresentacionPeso
	err := r.DB.Where("id = ?", id).First(&peso).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return peso, fmt.Errorf("configuración de peso con id %d no encontrada", id)
		}
		return peso, fmt.Errorf("error al recuperar configuración de peso con id %d: %w", id, err)
	}
	return peso, nil
}

// CreatePresentacionPeso crea una nueva configuración de peso
func (r *presentacionPesoRepo) CreatePresentacionPeso(peso *PresentacionPeso) error {
	// Verificar si la presentación existe
	var count int64
	if err := r.DB.Model(&Presentacion{}).Where("id = ?", peso.PresentacionID).Count(&count).Error; err != nil {
		return fmt.Errorf("error al verificar la presentación: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("la presentación con ID %d no existe", peso.PresentacionID)
	}

	err := r.DB.Create(peso).Error
	if err != nil {
		return fmt.Errorf("error al crear configuración de peso: %w", err)
	}
	return nil
}

// UpdatePresentacionPeso actualiza una configuración de peso existente
func (r *presentacionPesoRepo) UpdatePresentacionPeso(id uint64, updatedPeso *PresentacionPeso) (PresentacionPeso, error) {
	var peso PresentacionPeso
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&peso).Error; err != nil {
			return err
		}

		// Verificar si la presentación existe si se está actualizando
		if updatedPeso.PresentacionID != peso.PresentacionID {
			var count int64
			if err := tx.Model(&Presentacion{}).Where("id = ?", updatedPeso.PresentacionID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar la presentación: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("la presentación con ID %d no existe", updatedPeso.PresentacionID)
			}
		}

		return tx.Model(&peso).Updates(updatedPeso).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return peso, fmt.Errorf("configuración de peso con id %d no encontrada", id)
		}
		return peso, fmt.Errorf("error al actualizar configuración de peso con id %d: %w", id, err)
	}
	return peso, nil
}

// DeletePresentacionPeso elimina una configuración de peso
func (r *presentacionPesoRepo) DeletePresentacionPeso(id uint64) error {
	result := r.DB.Delete(&PresentacionPeso{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar configuración de peso: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("configuración de peso con id %d no encontrada", id)
	}
	return nil
}

// SeedPresentacionesPeso inicializa múltiples configuraciones de peso
func (r *presentacionPesoRepo) SeedPresentacionesPeso(pesos []PresentacionPeso) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&pesos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar configuraciones de peso: %w", err)
	}
	return nil
}

// PatchPresentacionPeso actualiza parcialmente una configuración de peso
func (r *presentacionPesoRepo) PatchPresentacionPeso(id uint64, fields map[string]interface{}) (PresentacionPeso, error) {
	var peso PresentacionPeso
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&peso).Error; err != nil {
			return err
		}

		// Verificar si la presentación existe si se está actualizando
		if presentacionID, ok := fields["presentacion_id"].(uint64); ok && presentacionID != peso.PresentacionID {
			var count int64
			if err := tx.Model(&Presentacion{}).Where("id = ?", presentacionID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar la presentación: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("la presentación con ID %d no existe", presentacionID)
			}
		}

		return tx.Model(&peso).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return peso, fmt.Errorf("configuración de peso con id %d no encontrada", id)
		}
		return peso, fmt.Errorf("error al actualizar parcialmente configuración de peso con id %d: %w", id, err)
	}
	return peso, nil
}
