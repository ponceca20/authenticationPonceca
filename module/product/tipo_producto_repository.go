package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// TipoProductoRepository define las operaciones de persistencia para TipoProducto
type TipoProductoRepository interface {
	GetAllTipoProductos(empresaID uint64) ([]TipoProducto, error)
	GetTipoProducto(id uint64, empresaID uint64) (TipoProducto, error)
	CreateTipoProducto(tp *TipoProducto) error
	UpdateTipoProducto(id uint64, empresaID uint64, tp *TipoProducto) (TipoProducto, error)
	DeleteTipoProducto(id uint64, empresaID uint64) error
	PatchTipoProducto(id uint64, empresaID uint64, fields map[string]interface{}) (TipoProducto, error)
}

type tipoProductoRepo struct {
	DB *gorm.DB
}

// NewTipoProductoRepository crea una instancia del repositorio
func NewTipoProductoRepository() TipoProductoRepository {
	return &tipoProductoRepo{DB: database.DBconn}
}

// GetAllTipoProductos obtiene todos los tipos de producto de una empresa
func (r *tipoProductoRepo) GetAllTipoProductos(empresaID uint64) ([]TipoProducto, error) {
	var items []TipoProducto
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar tipos de producto: %w", err)
	}
	return items, nil
}

// GetTipoProducto obtiene un tipo de producto específico
func (r *tipoProductoRepo) GetTipoProducto(id uint64, empresaID uint64) (TipoProducto, error) {
	var item TipoProducto
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("tipo de producto con id %d no encontrado", id)
		}
		return item, fmt.Errorf("error al recuperar tipo de producto %d: %w", id, err)
	}
	return item, nil
}

// CreateTipoProducto crea un nuevo tipo de producto
func (r *tipoProductoRepo) CreateTipoProducto(tp *TipoProducto) error {
	if err := r.DB.Create(tp).Error; err != nil {
		return fmt.Errorf("error al crear tipo de producto: %w", err)
	}
	return nil
}

// UpdateTipoProducto actualiza un tipo de producto existente
func (r *tipoProductoRepo) UpdateTipoProducto(id uint64, empresaID uint64, updated *TipoProducto) (TipoProducto, error) {
	var item TipoProducto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(updated).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar tipo de producto %d: %w", id, err)
	}
	return item, nil
}

// DeleteTipoProducto elimina un tipo de producto
func (r *tipoProductoRepo) DeleteTipoProducto(id uint64, empresaID uint64) error {
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&TipoProducto{}).Error
	if err != nil {
		return fmt.Errorf("error al eliminar tipo de producto: %w", err)
	}
	return nil
}

// PatchTipoProducto actualiza parcialmente un tipo de producto
func (r *tipoProductoRepo) PatchTipoProducto(id uint64, empresaID uint64, fields map[string]interface{}) (TipoProducto, error) {
	var item TipoProducto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar parcialmente tipo de producto %d: %w", id, err)
	}
	return item, nil
}
