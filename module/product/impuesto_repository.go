package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ImpuestoRepository define las operaciones de persistencia para Impuesto
type ImpuestoRepository interface {
	GetAllImpuestos(empresaID uint64) ([]Impuesto, error)
	GetImpuesto(id uint64, empresaID uint64) (Impuesto, error)
	CreateImpuesto(impuesto *Impuesto) error
	UpdateImpuesto(id uint64, empresaID uint64, impuesto *Impuesto) (Impuesto, error)
	DeleteImpuesto(id uint64, empresaID uint64) error
	SeedImpuestos(impuestos []Impuesto) error
	PatchImpuesto(id uint64, empresaID uint64, fields map[string]interface{}) (Impuesto, error)
}

type impuestoRepo struct {
	DB *gorm.DB
}

func NewImpuestoRepository() ImpuestoRepository {
	return &impuestoRepo{DB: database.DBconn}
}

// GetAllImpuestos obtiene todos los impuestos para una empresa
func (r *impuestoRepo) GetAllImpuestos(empresaID uint64) ([]Impuesto, error) {
	var impuestos []Impuesto
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&impuestos).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving impuestos: %w", err)
	}
	return impuestos, nil
}

// GetImpuesto obtiene un impuesto específico
func (r *impuestoRepo) GetImpuesto(id uint64, empresaID uint64) (Impuesto, error) {
	var impuesto Impuesto
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&impuesto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return impuesto, fmt.Errorf("impuesto with id %d not found", id)
		}
		return impuesto, fmt.Errorf("error retrieving impuesto with id %d: %w", id, err)
	}
	return impuesto, nil
}

// CreateImpuesto crea un nuevo impuesto
func (r *impuestoRepo) CreateImpuesto(impuesto *Impuesto) error {
	err := r.DB.Create(impuesto).Error
	if err != nil {
		return fmt.Errorf("error creating impuesto: %w", err)
	}
	return nil
}

// UpdateImpuesto actualiza un impuesto existente
func (r *impuestoRepo) UpdateImpuesto(id uint64, empresaID uint64, updatedImpuesto *Impuesto) (Impuesto, error) {
	var impuesto Impuesto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&impuesto).Error; err != nil {
			return err
		}
		return tx.Model(&impuesto).Updates(updatedImpuesto).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return impuesto, fmt.Errorf("impuesto with id %d not found", id)
		}
		return impuesto, fmt.Errorf("error updating impuesto with id %d: %w", id, err)
	}
	return impuesto, nil
}

// DeleteImpuesto elimina un impuesto
func (r *impuestoRepo) DeleteImpuesto(id uint64, empresaID uint64) error {
	result := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Impuesto{})
	if result.Error != nil {
		return fmt.Errorf("error deleting impuesto: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("impuesto with id %d not found", id)
	}
	return nil
}

// SeedImpuestos inicializa múltiples impuestos
func (r *impuestoRepo) SeedImpuestos(impuestos []Impuesto) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&impuestos).Error
	})
	if err != nil {
		return fmt.Errorf("error seeding impuestos: %w", err)
	}
	return nil
}

// PatchImpuesto actualiza parcialmente un impuesto
func (r *impuestoRepo) PatchImpuesto(id uint64, empresaID uint64, fields map[string]interface{}) (Impuesto, error) {
	var impuesto Impuesto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&impuesto).Error; err != nil {
			return err
		}
		return tx.Model(&impuesto).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return impuesto, fmt.Errorf("impuesto with id %d not found", id)
		}
		return impuesto, fmt.Errorf("error patching impuesto with id %d: %w", id, err)
	}
	return impuesto, nil
}
