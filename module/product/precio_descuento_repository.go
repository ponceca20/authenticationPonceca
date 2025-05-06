package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PrecioDescuentoRepository define la interfaz para operaciones de persistencia.
type PrecioDescuentoRepository interface {
	GetAllPrecioDescuento(empresaID uint64) ([]PrecioDescuentoPorCantidad, error)
	GetPrecioDescuento(id uint64, empresaID uint64) (PrecioDescuentoPorCantidad, error)
	CreatePrecioDescuento(pd *PrecioDescuentoPorCantidad) error
	UpdatePrecioDescuento(id uint64, empresaID uint64, pd *PrecioDescuentoPorCantidad) (PrecioDescuentoPorCantidad, error)
	DeletePrecioDescuento(id uint64, empresaID uint64) error
	PatchPrecioDescuento(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioDescuentoPorCantidad, error)
}

type precioDescuentoRepo struct {
	DB *gorm.DB
}

func NewPrecioDescuentoRepository() PrecioDescuentoRepository {
	return &precioDescuentoRepo{DB: database.DBconn}
}

func (r *precioDescuentoRepo) GetAllPrecioDescuento(empresaID uint64) ([]PrecioDescuentoPorCantidad, error) {
	var items []PrecioDescuentoPorCantidad
	// Se filtra utilizando join con la tabla Producto.
	err := r.DB.Joins("JOIN productos ON productos.id = precio_descuento_por_cantidades.producto_id").
		Where("productos.empresa_id = ?", empresaID).
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving precio descuento: %w", err)
	}
	return items, nil
}

func (r *precioDescuentoRepo) GetPrecioDescuento(id uint64, empresaID uint64) (PrecioDescuentoPorCantidad, error) {
	var item PrecioDescuentoPorCantidad
	err := r.DB.Joins("JOIN productos ON productos.id = precio_descuento_por_cantidades.producto_id").
		Where("precio_descuento_por_cantidades.id = ? AND productos.empresa_id = ?", id, empresaID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio descuento with id %d not found", id)
		}
		return item, fmt.Errorf("error retrieving precio descuento with id %d: %w", id, err)
	}
	return item, nil
}

func (r *precioDescuentoRepo) CreatePrecioDescuento(pd *PrecioDescuentoPorCantidad) error {
	err := r.DB.Create(pd).Error
	if err != nil {
		return fmt.Errorf("error creating precio descuento: %w", err)
	}
	return nil
}

func (r *precioDescuentoRepo) UpdatePrecioDescuento(id uint64, empresaID uint64, pd *PrecioDescuentoPorCantidad) (PrecioDescuentoPorCantidad, error) {
	var item PrecioDescuentoPorCantidad
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN productos ON productos.id = precio_descuento_por_cantidades.producto_id").
			Where("precio_descuento_por_cantidades.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(pd).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio descuento with id %d not found", id)
		}
		return item, fmt.Errorf("error updating precio descuento with id %d: %w", id, err)
	}
	return item, nil
}

func (r *precioDescuentoRepo) DeletePrecioDescuento(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var item PrecioDescuentoPorCantidad
		if err := tx.Joins("JOIN productos ON productos.id = precio_descuento_por_cantidades.producto_id").
			Where("precio_descuento_por_cantidades.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		result := tx.Delete(&item)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("precio descuento with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting precio descuento: %w", err)
	}
	return nil
}

func (r *precioDescuentoRepo) PatchPrecioDescuento(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioDescuentoPorCantidad, error) {
	var item PrecioDescuentoPorCantidad
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN productos ON productos.id = precio_descuento_por_cantidades.producto_id").
			Where("precio_descuento_por_cantidades.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("precio descuento with id %d not found", id)
		}
		return item, fmt.Errorf("error patching precio descuento with id %d: %w", id, err)
	}
	return item, nil
}
