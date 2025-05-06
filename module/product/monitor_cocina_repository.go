package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// MonitorCocinaRepository define la interfaz para operaciones de persistencia de MonitorCocina
type MonitorCocinaRepository interface {
	GetAllMonitoresCocina(empresaID uint64) ([]MonitorCocina, error)
	GetMonitorCocina(id uint64, empresaID uint64) (MonitorCocina, error)
	CreateMonitorCocina(monitor *MonitorCocina) error
	UpdateMonitorCocina(id uint64, empresaID uint64, monitor *MonitorCocina) (MonitorCocina, error)
	DeleteMonitorCocina(id uint64, empresaID uint64) error
	SeedMonitoresCocina(monitores []MonitorCocina) error
	PatchMonitorCocina(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorCocina, error)
}

type monitorCocinaRepo struct {
	DB *gorm.DB
}

// NewMonitorCocinaRepository crea una instancia del repositorio
func NewMonitorCocinaRepository() MonitorCocinaRepository {
	return &monitorCocinaRepo{DB: database.DBconn}
}

// GetAllMonitoresCocina obtiene todos los monitores de cocina para una empresa
func (r *monitorCocinaRepo) GetAllMonitoresCocina(empresaID uint64) ([]MonitorCocina, error) {
	var monitores []MonitorCocina
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&monitores).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar monitores de cocina: %w", err)
	}
	return monitores, nil
}

// GetMonitorCocina obtiene un monitor de cocina específico
func (r *monitorCocinaRepo) GetMonitorCocina(id uint64, empresaID uint64) (MonitorCocina, error) {
	var monitor MonitorCocina
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor de cocina con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al recuperar monitor de cocina con id %d: %w", id, err)
	}
	return monitor, nil
}

// CreateMonitorCocina crea un nuevo monitor de cocina
func (r *monitorCocinaRepo) CreateMonitorCocina(monitor *MonitorCocina) error {
	err := r.DB.Create(monitor).Error
	if err != nil {
		return fmt.Errorf("error al crear monitor de cocina: %w", err)
	}
	return nil
}

// UpdateMonitorCocina actualiza un monitor de cocina existente
func (r *monitorCocinaRepo) UpdateMonitorCocina(id uint64, empresaID uint64, updatedMonitor *MonitorCocina) (MonitorCocina, error) {
	var monitor MonitorCocina
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error; err != nil {
			return err
		}
		return tx.Model(&monitor).Updates(updatedMonitor).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor de cocina con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al actualizar monitor de cocina con id %d: %w", id, err)
	}
	return monitor, nil
}

// DeleteMonitorCocina elimina un monitor de cocina
func (r *monitorCocinaRepo) DeleteMonitorCocina(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si existen productos asociados
		var count int64
		if err := tx.Model(&Producto{}).Where("monitor_cocina_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar el monitor de cocina porque tiene %d productos asociados", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&MonitorCocina{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("monitor de cocina con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar monitor de cocina: %w", err)
	}
	return nil
}

// SeedMonitoresCocina inicializa múltiples monitores de cocina
func (r *monitorCocinaRepo) SeedMonitoresCocina(monitores []MonitorCocina) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&monitores).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar monitores de cocina: %w", err)
	}
	return nil
}

// PatchMonitorCocina actualiza parcialmente un monitor de cocina
func (r *monitorCocinaRepo) PatchMonitorCocina(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorCocina, error) {
	var monitor MonitorCocina
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error; err != nil {
			return err
		}
		return tx.Model(&monitor).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor de cocina con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al actualizar parcialmente monitor de cocina con id %d: %w", id, err)
	}
	return monitor, nil
}
