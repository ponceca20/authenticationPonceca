package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PrecioHistoricoRepository define la interfaz para operaciones de persistencia.
type PrecioHistoricoRepository interface {
	GetAllPrecioHistorico(empresaID uint64) ([]PrecioHistorico, error)
	GetPrecioHistorico(id uint64, empresaID uint64) (PrecioHistorico, error)
	CreatePrecioHistorico(ph *PrecioHistorico) error
	UpdatePrecioHistorico(id uint64, empresaID uint64, ph *PrecioHistorico) (PrecioHistorico, error)
	DeletePrecioHistorico(id uint64, empresaID uint64) error
	PatchPrecioHistorico(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioHistorico, error)
}

type precioHistoricoRepo struct {
	DB *gorm.DB
}

func NewPrecioHistoricoRepository() PrecioHistoricoRepository {
	return &precioHistoricoRepo{DB: database.DBconn}
}

func (r *precioHistoricoRepo) GetAllPrecioHistorico(empresaID uint64) ([]PrecioHistorico, error) {
	var items []PrecioHistorico
	// Se filtra mediante join: presentacion -> producto para obtener empresa_id.
	err := r.DB.Joins("JOIN presentacions ON presentacions.id = precio_historicos.presentacion_id").
		Joins("JOIN productos ON productos.id = presentacions.producto_id").
		Where("productos.empresa_id = ?", empresaID).
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving precio historico: %w", err)
	}
	return items, nil
}

func (r *precioHistoricoRepo) GetPrecioHistorico(id uint64, empresaID uint64) (PrecioHistorico, error) {
	var item PrecioHistorico
	err := r.DB.Joins("JOIN presentacions ON presentacions.id = precio_historicos.presentacion_id").
		Joins("JOIN productos ON productos.id = presentacions.producto_id").
		Where("precio_historicos.id = ? AND productos.empresa_id = ?", id, empresaID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio historico with id %d not found", id)
		}
		return item, fmt.Errorf("error retrieving precio historico with id %d: %w", id, err)
	}
	return item, nil
}

func (r *precioHistoricoRepo) CreatePrecioHistorico(ph *PrecioHistorico) error {
	err := r.DB.Create(ph).Error
	if err != nil {
		return fmt.Errorf("error creating precio historico: %w", err)
	}
	return nil
}

func (r *precioHistoricoRepo) UpdatePrecioHistorico(id uint64, empresaID uint64, ph *PrecioHistorico) (PrecioHistorico, error) {
	var item PrecioHistorico
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN presentacions ON presentacions.id = precio_historicos.presentacion_id").
			Joins("JOIN productos ON productos.id = presentacions.producto_id").
			Where("precio_historicos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(ph).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio historico with id %d not found", id)
		}
		return item, fmt.Errorf("error updating precio historico with id %d: %w", id, err)
	}
	return item, nil
}

func (r *precioHistoricoRepo) DeletePrecioHistorico(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var item PrecioHistorico
		if err := tx.Joins("JOIN presentacions ON presentacions.id = precio_historicos.presentacion_id").
			Joins("JOIN productos ON productos.id = presentacions.producto_id").
			Where("precio_historicos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		result := tx.Delete(&item)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("precio historico with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting precio historico: %w", err)
	}
	return nil
}

func (r *precioHistoricoRepo) PatchPrecioHistorico(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioHistorico, error) {
	var item PrecioHistorico
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN presentacions ON presentacions.id = precio_historicos.presentacion_id").
			Joins("JOIN productos ON productos.id = presentacions.producto_id").
			Where("precio_historicos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio historico with id %d not found", id)
		}
		return item, fmt.Errorf("error patching precio historico with id %d: %w", id, err)
	}
	return item, nil
}
