package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PresentacionTiempoRepository define la interfaz para operaciones de persistencia de PresentacionTiempo
type PresentacionTiempoRepository interface {
	GetAllPresentacionesTiempo(presentacionID uint64) ([]PresentacionTiempo, error)
	GetPresentacionTiempo(id uint64) (PresentacionTiempo, error)
	CreatePresentacionTiempo(tiempo *PresentacionTiempo) error
	UpdatePresentacionTiempo(id uint64, tiempo *PresentacionTiempo) (PresentacionTiempo, error)
	DeletePresentacionTiempo(id uint64) error
	SeedPresentacionesTiempo(tiempos []PresentacionTiempo) error
	PatchPresentacionTiempo(id uint64, fields map[string]interface{}) (PresentacionTiempo, error)
}

type presentacionTiempoRepo struct {
	DB *gorm.DB
}

// NewPresentacionTiempoRepository crea una instancia del repositorio
func NewPresentacionTiempoRepository() PresentacionTiempoRepository {
	return &presentacionTiempoRepo{DB: database.DBconn}
}

// GetAllPresentacionesTiempo obtiene todas las configuraciones de tiempo para una presentación
func (r *presentacionTiempoRepo) GetAllPresentacionesTiempo(presentacionID uint64) ([]PresentacionTiempo, error) {
	var tiempos []PresentacionTiempo
	err := r.DB.Where("presentacion_id = ?", presentacionID).Find(&tiempos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar configuraciones de tiempo: %w", err)
	}
	return tiempos, nil
}

// GetPresentacionTiempo obtiene una configuración de tiempo específica
func (r *presentacionTiempoRepo) GetPresentacionTiempo(id uint64) (PresentacionTiempo, error) {
	var tiempo PresentacionTiempo
	err := r.DB.Where("id = ?", id).First(&tiempo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tiempo, fmt.Errorf("configuración de tiempo con id %d no encontrada", id)
		}
		return tiempo, fmt.Errorf("error al recuperar configuración de tiempo con id %d: %w", id, err)
	}
	return tiempo, nil
}

// CreatePresentacionTiempo crea una nueva configuración de tiempo
func (r *presentacionTiempoRepo) CreatePresentacionTiempo(tiempo *PresentacionTiempo) error {
	// Verificar si la presentación existe
	var count int64
	if err := r.DB.Model(&Presentacion{}).Where("id = ?", tiempo.PresentacionID).Count(&count).Error; err != nil {
		return fmt.Errorf("error al verificar la presentación: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("la presentación con ID %d no existe", tiempo.PresentacionID)
	}

	err := r.DB.Create(tiempo).Error
	if err != nil {
		return fmt.Errorf("error al crear configuración de tiempo: %w", err)
	}
	return nil
}

// UpdatePresentacionTiempo actualiza una configuración de tiempo existente
func (r *presentacionTiempoRepo) UpdatePresentacionTiempo(id uint64, updatedTiempo *PresentacionTiempo) (PresentacionTiempo, error) {
	var tiempo PresentacionTiempo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&tiempo).Error; err != nil {
			return err
		}

		// Verificar si la presentación existe si se está actualizando
		if updatedTiempo.PresentacionID != tiempo.PresentacionID {
			var count int64
			if err := tx.Model(&Presentacion{}).Where("id = ?", updatedTiempo.PresentacionID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar la presentación: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("la presentación con ID %d no existe", updatedTiempo.PresentacionID)
			}
		}

		return tx.Model(&tiempo).Updates(updatedTiempo).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tiempo, fmt.Errorf("configuración de tiempo con id %d no encontrada", id)
		}
		return tiempo, fmt.Errorf("error al actualizar configuración de tiempo con id %d: %w", id, err)
	}
	return tiempo, nil
}

// DeletePresentacionTiempo elimina una configuración de tiempo
func (r *presentacionTiempoRepo) DeletePresentacionTiempo(id uint64) error {
	result := r.DB.Delete(&PresentacionTiempo{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar configuración de tiempo: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("configuración de tiempo con id %d no encontrada", id)
	}
	return nil
}

// SeedPresentacionesTiempo inicializa múltiples configuraciones de tiempo
func (r *presentacionTiempoRepo) SeedPresentacionesTiempo(tiempos []PresentacionTiempo) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&tiempos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar configuraciones de tiempo: %w", err)
	}
	return nil
}

// PatchPresentacionTiempo actualiza parcialmente una configuración de tiempo
func (r *presentacionTiempoRepo) PatchPresentacionTiempo(id uint64, fields map[string]interface{}) (PresentacionTiempo, error) {
	var tiempo PresentacionTiempo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&tiempo).Error; err != nil {
			return err
		}

		// Verificar si la presentación existe si se está actualizando
		if presentacionID, ok := fields["presentacion_id"].(uint64); ok && presentacionID != tiempo.PresentacionID {
			var count int64
			if err := tx.Model(&Presentacion{}).Where("id = ?", presentacionID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar la presentación: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("la presentación con ID %d no existe", presentacionID)
			}
		}

		return tx.Model(&tiempo).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tiempo, fmt.Errorf("configuración de tiempo con id %d no encontrada", id)
		}
		return tiempo, fmt.Errorf("error al actualizar parcialmente configuración de tiempo con id %d: %w", id, err)
	}
	return tiempo, nil
}
