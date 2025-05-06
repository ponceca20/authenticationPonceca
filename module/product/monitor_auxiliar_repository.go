package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// MonitorAuxiliarRepository define la interfaz para operaciones de persistencia de MonitorAuxiliar
type MonitorAuxiliarRepository interface {
	GetAllMonitoresAuxiliar(empresaID uint64) ([]MonitorAuxiliar, error)
	GetMonitorAuxiliar(id uint64, empresaID uint64) (MonitorAuxiliar, error)
	CreateMonitorAuxiliar(monitor *MonitorAuxiliar) error
	UpdateMonitorAuxiliar(id uint64, empresaID uint64, monitor *MonitorAuxiliar) (MonitorAuxiliar, error)
	DeleteMonitorAuxiliar(id uint64, empresaID uint64) error
	SeedMonitoresAuxiliar(monitores []MonitorAuxiliar) error
	PatchMonitorAuxiliar(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorAuxiliar, error)
}

type monitorAuxiliarRepo struct {
	DB *gorm.DB
}

// NewMonitorAuxiliarRepository crea una instancia del repositorio
func NewMonitorAuxiliarRepository() MonitorAuxiliarRepository {
	return &monitorAuxiliarRepo{DB: database.DBconn}
}

// GetAllMonitoresAuxiliar obtiene todos los monitores auxiliares para una empresa
func (r *monitorAuxiliarRepo) GetAllMonitoresAuxiliar(empresaID uint64) ([]MonitorAuxiliar, error) {
	var monitores []MonitorAuxiliar
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&monitores).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar monitores auxiliares: %w", err)
	}
	return monitores, nil
}

// GetMonitorAuxiliar obtiene un monitor auxiliar específico
func (r *monitorAuxiliarRepo) GetMonitorAuxiliar(id uint64, empresaID uint64) (MonitorAuxiliar, error) {
	var monitor MonitorAuxiliar
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor auxiliar con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al recuperar monitor auxiliar con id %d: %w", id, err)
	}
	return monitor, nil
}

// CreateMonitorAuxiliar crea un nuevo monitor auxiliar
func (r *monitorAuxiliarRepo) CreateMonitorAuxiliar(monitor *MonitorAuxiliar) error {
	err := r.DB.Create(monitor).Error
	if err != nil {
		return fmt.Errorf("error al crear monitor auxiliar: %w", err)
	}
	return nil
}

// UpdateMonitorAuxiliar actualiza un monitor auxiliar existente
func (r *monitorAuxiliarRepo) UpdateMonitorAuxiliar(id uint64, empresaID uint64, updatedMonitor *MonitorAuxiliar) (MonitorAuxiliar, error) {
	var monitor MonitorAuxiliar
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error; err != nil {
			return err
		}
		return tx.Model(&monitor).Updates(updatedMonitor).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor auxiliar con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al actualizar monitor auxiliar con id %d: %w", id, err)
	}
	return monitor, nil
}

// DeleteMonitorAuxiliar elimina un monitor auxiliar
func (r *monitorAuxiliarRepo) DeleteMonitorAuxiliar(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si existen productos asociados
		var count int64
		if err := tx.Model(&Producto{}).Where("monitor_auxiliar_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar el monitor auxiliar porque tiene %d productos asociados", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&MonitorAuxiliar{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("monitor auxiliar con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar monitor auxiliar: %w", err)
	}
	return nil
}

// SeedMonitoresAuxiliar inicializa múltiples monitores auxiliares
func (r *monitorAuxiliarRepo) SeedMonitoresAuxiliar(monitores []MonitorAuxiliar) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&monitores).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar monitores auxiliares: %w", err)
	}
	return nil
}

// PatchMonitorAuxiliar actualiza parcialmente un monitor auxiliar
func (r *monitorAuxiliarRepo) PatchMonitorAuxiliar(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorAuxiliar, error) {
	var monitor MonitorAuxiliar
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&monitor).Error; err != nil {
			return err
		}
		return tx.Model(&monitor).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return monitor, fmt.Errorf("monitor auxiliar con id %d no encontrado", id)
		}
		return monitor, fmt.Errorf("error al actualizar parcialmente monitor auxiliar con id %d: %w", id, err)
	}
	return monitor, nil
}
