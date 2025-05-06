package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// UnidadRepository define la interfaz para operaciones de persistencia de Unidad
type UnidadRepository interface {
	GetAllUnidades(empresaID uint64) ([]Unidad, error)
	GetUnidad(id uint64, empresaID uint64) (Unidad, error)
	CreateUnidad(unidad *Unidad) error
	UpdateUnidad(id uint64, empresaID uint64, unidad *Unidad) (Unidad, error)
	DeleteUnidad(id uint64, empresaID uint64) error
	SeedUnidades(unidades []Unidad) error
	PatchUnidad(id uint64, empresaID uint64, fields map[string]interface{}) (Unidad, error)
}

type unidadRepo struct {
	DB *gorm.DB
}

// NewUnidadRepository crea una instancia del repositorio
func NewUnidadRepository() UnidadRepository {
	return &unidadRepo{DB: database.DBconn}
}

// GetAllUnidades obtiene todas las unidades para una empresa
func (r *unidadRepo) GetAllUnidades(empresaID uint64) ([]Unidad, error) {
	var unidades []Unidad
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&unidades).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar unidades: %w", err)
	}
	return unidades, nil
}

// GetUnidad obtiene una unidad específica
func (r *unidadRepo) GetUnidad(id uint64, empresaID uint64) (Unidad, error) {
	var unidad Unidad
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&unidad).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unidad, fmt.Errorf("unidad con id %d no encontrada", id)
		}
		return unidad, fmt.Errorf("error al recuperar unidad con id %d: %w", id, err)
	}
	return unidad, nil
}

// CreateUnidad crea una nueva unidad
func (r *unidadRepo) CreateUnidad(unidad *Unidad) error {
	err := r.DB.Create(unidad).Error
	if err != nil {
		return fmt.Errorf("error al crear unidad: %w", err)
	}
	return nil
}

// UpdateUnidad actualiza una unidad existente
func (r *unidadRepo) UpdateUnidad(id uint64, empresaID uint64, updatedUnidad *Unidad) (Unidad, error) {
	var unidad Unidad
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&unidad).Error; err != nil {
			return err
		}
		return tx.Model(&unidad).Updates(updatedUnidad).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unidad, fmt.Errorf("unidad con id %d no encontrada", id)
		}
		return unidad, fmt.Errorf("error al actualizar unidad con id %d: %w", id, err)
	}
	return unidad, nil
}

// DeleteUnidad elimina una unidad
func (r *unidadRepo) DeleteUnidad(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si existen productos asociados
		var count int64
		if err := tx.Model(&Producto{}).Where("unidad_control_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar la unidad porque tiene %d productos asociados", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Unidad{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("unidad con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar unidad: %w", err)
	}
	return nil
}

// SeedUnidades inicializa múltiples unidades
func (r *unidadRepo) SeedUnidades(unidades []Unidad) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&unidades).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar unidades: %w", err)
	}
	return nil
}

// PatchUnidad actualiza parcialmente una unidad
func (r *unidadRepo) PatchUnidad(id uint64, empresaID uint64, fields map[string]interface{}) (Unidad, error) {
	var unidad Unidad
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&unidad).Error; err != nil {
			return err
		}
		return tx.Model(&unidad).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unidad, fmt.Errorf("unidad con id %d no encontrada", id)
		}
		return unidad, fmt.Errorf("error al actualizar parcialmente unidad con id %d: %w", id, err)
	}
	return unidad, nil
}
