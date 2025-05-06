package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// FamiliaRepository define la interfaz para operaciones de persistencia de Familia
type FamiliaRepository interface {
	GetAllFamilias(empresaID uint64) ([]Familia, error)
	GetFamiliasByCartaID(cartaID uint64, empresaID uint64) ([]Familia, error)
	GetFamilia(id uint64, empresaID uint64) (Familia, error)
	CreateFamilia(familia *Familia) error
	UpdateFamilia(id uint64, empresaID uint64, familia *Familia) (Familia, error)
	DeleteFamilia(id uint64, empresaID uint64) error
	SeedFamilias(familias []Familia) error
	PatchFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (Familia, error)
}

type familiaRepo struct {
	DB *gorm.DB
}

// NewFamiliaRepository crea una instancia del repositorio
func NewFamiliaRepository() FamiliaRepository {
	return &familiaRepo{DB: database.DBconn}
}

// GetAllFamilias obtiene todas las familias para una empresa
func (r *familiaRepo) GetAllFamilias(empresaID uint64) ([]Familia, error) {
	var familias []Familia
	err := r.DB.Preload("ProductoFamilias").Where("empresa_id = ?", empresaID).Find(&familias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar familias: %w", err)
	}
	return familias, nil
}

// GetFamiliasByCartaID obtiene todas las familias para una carta específica
func (r *familiaRepo) GetFamiliasByCartaID(cartaID uint64, empresaID uint64) ([]Familia, error) {
	var familias []Familia
	err := r.DB.Preload("ProductoFamilias").Where("carta_id = ? AND empresa_id = ?", cartaID, empresaID).Find(&familias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar familias para carta ID %d: %w", cartaID, err)
	}
	return familias, nil
}

// GetFamilia obtiene una familia específica
func (r *familiaRepo) GetFamilia(id uint64, empresaID uint64) (Familia, error) {
	var familia Familia
	err := r.DB.Preload("ProductoFamilias").Where("id = ? AND empresa_id = ?", id, empresaID).First(&familia).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return familia, fmt.Errorf("familia con id %d no encontrada", id)
		}
		return familia, fmt.Errorf("error al recuperar familia con id %d: %w", id, err)
	}
	return familia, nil
}

// CreateFamilia crea una nueva familia
func (r *familiaRepo) CreateFamilia(familia *Familia) error {
	// Verificar que la carta existe y pertenece a la misma empresa
	var carta Carta
	if err := r.DB.Where("id = ? AND empresa_id = ?", familia.CartaID, familia.EmpresaID).First(&carta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("carta con id %d no encontrada o no pertenece a la empresa", familia.CartaID)
		}
		return fmt.Errorf("error al verificar carta: %w", err)
	}

	err := r.DB.Create(familia).Error
	if err != nil {
		return fmt.Errorf("error al crear familia: %w", err)
	}
	return nil
}

// UpdateFamilia actualiza una familia existente
func (r *familiaRepo) UpdateFamilia(id uint64, empresaID uint64, updatedFamilia *Familia) (Familia, error) {
	var familia Familia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&familia).Error; err != nil {
			return err
		}

		// Si se está cambiando la carta, verificar que la nueva carta existe y pertenece a la misma empresa
		if updatedFamilia.CartaID != 0 && updatedFamilia.CartaID != familia.CartaID {
			var carta Carta
			if err := tx.Where("id = ? AND empresa_id = ?", updatedFamilia.CartaID, empresaID).First(&carta).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("carta con id %d no encontrada o no pertenece a la empresa", updatedFamilia.CartaID)
				}
				return fmt.Errorf("error al verificar carta: %w", err)
			}
		}

		return tx.Model(&familia).Updates(updatedFamilia).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return familia, fmt.Errorf("familia con id %d no encontrada", id)
		}
		return familia, fmt.Errorf("error al actualizar familia con id %d: %w", id, err)
	}
	return familia, nil
}

// DeleteFamilia elimina una familia
func (r *familiaRepo) DeleteFamilia(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si hay productos asociados
		var count int64
		if err := tx.Model(&ProductoFamilia{}).Where("familia_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar la familia porque tiene %d productos asociados", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Familia{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("familia con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar familia: %w", err)
	}
	return nil
}

// SeedFamilias inicializa múltiples familias
func (r *familiaRepo) SeedFamilias(familias []Familia) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&familias).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar familias: %w", err)
	}
	return nil
}

// PatchFamilia actualiza parcialmente una familia
func (r *familiaRepo) PatchFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (Familia, error) {
	var familia Familia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&familia).Error; err != nil {
			return err
		}

		// Si se está cambiando la carta, verificar que la nueva carta existe y pertenece a la misma empresa
		if cartaID, ok := fields["carta_id"].(float64); ok {
			var carta Carta
			if err := tx.Where("id = ? AND empresa_id = ?", uint64(cartaID), empresaID).First(&carta).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("carta con id %d no encontrada o no pertenece a la empresa", uint64(cartaID))
				}
				return fmt.Errorf("error al verificar carta: %w", err)
			}
		}

		return tx.Model(&familia).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return familia, fmt.Errorf("familia con id %d no encontrada", id)
		}
		return familia, fmt.Errorf("error al actualizar parcialmente familia con id %d: %w", id, err)
	}
	return familia, nil
}
