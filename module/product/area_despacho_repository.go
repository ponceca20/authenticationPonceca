package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// AreaDespachoRepository define la interfaz para operaciones de persistencia de AreaDespacho
type AreaDespachoRepository interface {
	GetAllAreasDespacho(empresaID uint64) ([]AreaDespacho, error)
	GetAreaDespacho(id uint64, empresaID uint64) (AreaDespacho, error)
	CreateAreaDespacho(areaDespacho *AreaDespacho) error
	UpdateAreaDespacho(id uint64, empresaID uint64, areaDespacho *AreaDespacho) (AreaDespacho, error)
	DeleteAreaDespacho(id uint64, empresaID uint64) error
	SeedAreasDespacho(areasDespacho []AreaDespacho) error
	PatchAreaDespacho(id uint64, empresaID uint64, fields map[string]interface{}) (AreaDespacho, error)
}

type areaDespachoRepo struct {
	DB *gorm.DB
}

// NewAreaDespachoRepository crea una instancia del repositorio
func NewAreaDespachoRepository() AreaDespachoRepository {
	return &areaDespachoRepo{DB: database.DBconn}
}

// GetAllAreasDespacho obtiene todas las áreas de despacho para una empresa
func (r *areaDespachoRepo) GetAllAreasDespacho(empresaID uint64) ([]AreaDespacho, error) {
	var areasDespacho []AreaDespacho
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&areasDespacho).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar áreas de despacho: %w", err)
	}
	return areasDespacho, nil
}

// GetAreaDespacho obtiene un área de despacho específica
func (r *areaDespachoRepo) GetAreaDespacho(id uint64, empresaID uint64) (AreaDespacho, error) {
	var areaDespacho AreaDespacho
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&areaDespacho).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return areaDespacho, fmt.Errorf("área de despacho con id %d no encontrada", id)
		}
		return areaDespacho, fmt.Errorf("error al recuperar área de despacho con id %d: %w", id, err)
	}
	return areaDespacho, nil
}

// CreateAreaDespacho crea una nueva área de despacho
func (r *areaDespachoRepo) CreateAreaDespacho(areaDespacho *AreaDespacho) error {
	err := r.DB.Create(areaDespacho).Error
	if err != nil {
		return fmt.Errorf("error al crear área de despacho: %w", err)
	}
	return nil
}

// UpdateAreaDespacho actualiza un área de despacho existente
func (r *areaDespachoRepo) UpdateAreaDespacho(id uint64, empresaID uint64, updatedAreaDespacho *AreaDespacho) (AreaDespacho, error) {
	var areaDespacho AreaDespacho
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&areaDespacho).Error; err != nil {
			return err
		}
		return tx.Model(&areaDespacho).Updates(updatedAreaDespacho).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return areaDespacho, fmt.Errorf("área de despacho con id %d no encontrada", id)
		}
		return areaDespacho, fmt.Errorf("error al actualizar área de despacho con id %d: %w", id, err)
	}
	return areaDespacho, nil
}

// DeleteAreaDespacho elimina un área de despacho
func (r *areaDespachoRepo) DeleteAreaDespacho(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si existen productos asociados
		var count int64
		if err := tx.Model(&Producto{}).Where("area_despacho_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar el área de despacho porque tiene %d productos asociados", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&AreaDespacho{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("área de despacho con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar área de despacho: %w", err)
	}
	return nil
}

// SeedAreasDespacho inicializa múltiples áreas de despacho
func (r *areaDespachoRepo) SeedAreasDespacho(areasDespacho []AreaDespacho) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&areasDespacho).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar áreas de despacho: %w", err)
	}
	return nil
}

// PatchAreaDespacho actualiza parcialmente un área de despacho
func (r *areaDespachoRepo) PatchAreaDespacho(id uint64, empresaID uint64, fields map[string]interface{}) (AreaDespacho, error) {
	var areaDespacho AreaDespacho
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&areaDespacho).Error; err != nil {
			return err
		}
		return tx.Model(&areaDespacho).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return areaDespacho, fmt.Errorf("área de despacho con id %d no encontrada", id)
		}
		return areaDespacho, fmt.Errorf("error al actualizar parcialmente área de despacho con id %d: %w", id, err)
	}
	return areaDespacho, nil
}
