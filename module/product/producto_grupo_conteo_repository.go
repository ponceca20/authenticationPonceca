package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ProductoGrupoConteoRepository define la interfaz para operaciones de persistencia.
type ProductoGrupoConteoRepository interface {
	GetAllProductoGrupoConteo(empresaID uint64) ([]ProductoGrupoConteo, error)
	GetProductoGrupoConteo(id uint64, empresaID uint64) (ProductoGrupoConteo, error)
	CreateProductoGrupoConteo(pg *ProductoGrupoConteo) error
	UpdateProductoGrupoConteo(id uint64, empresaID uint64, pg *ProductoGrupoConteo) (ProductoGrupoConteo, error)
	DeleteProductoGrupoConteo(id uint64, empresaID uint64) error
	PatchProductoGrupoConteo(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoGrupoConteo, error)
}

type productoGrupoConteoRepo struct {
	DB *gorm.DB
}

func NewProductoGrupoConteoRepository() ProductoGrupoConteoRepository {
	return &productoGrupoConteoRepo{DB: database.DBconn}
}

func (r *productoGrupoConteoRepo) GetAllProductoGrupoConteo(empresaID uint64) ([]ProductoGrupoConteo, error) {
	var items []ProductoGrupoConteo
	// Se filtra vía join con la tabla Producto para verificar empresa_id.
	err := r.DB.Joins("JOIN productos ON productos.id = producto_grupo_conteos.producto_id").
		Where("productos.empresa_id = ?", empresaID).
		Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error retrieving producto grupo conteo: %w", err)
	}
	return items, nil
}

func (r *productoGrupoConteoRepo) GetProductoGrupoConteo(id uint64, empresaID uint64) (ProductoGrupoConteo, error) {
	var item ProductoGrupoConteo
	err := r.DB.Joins("JOIN productos ON productos.id = producto_grupo_conteos.producto_id").
		Where("producto_grupo_conteos.id = ? AND productos.empresa_id = ?", id, empresaID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("producto grupo conteo with id %d not found", id)
		}
		return item, fmt.Errorf("error retrieving producto grupo conteo with id %d: %w", id, err)
	}
	return item, nil
}

func (r *productoGrupoConteoRepo) CreateProductoGrupoConteo(pg *ProductoGrupoConteo) error {
	err := r.DB.Create(pg).Error
	if err != nil {
		return fmt.Errorf("error creating producto grupo conteo: %w", err)
	}
	return nil
}

func (r *productoGrupoConteoRepo) UpdateProductoGrupoConteo(id uint64, empresaID uint64, pg *ProductoGrupoConteo) (ProductoGrupoConteo, error) {
	var item ProductoGrupoConteo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN productos ON productos.id = producto_grupo_conteos.producto_id").
			Where("producto_grupo_conteos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(pg).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("producto grupo conteo with id %d not found", id)
		}
		return item, fmt.Errorf("error updating producto grupo conteo with id %d: %w", id, err)
	}
	return item, nil
}

func (r *productoGrupoConteoRepo) DeleteProductoGrupoConteo(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var item ProductoGrupoConteo
		if err := tx.Joins("JOIN productos ON productos.id = producto_grupo_conteos.producto_id").
			Where("producto_grupo_conteos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		result := tx.Delete(&item)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("producto grupo conteo with id %d not found", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting producto grupo conteo: %w", err)
	}
	return nil
}

func (r *productoGrupoConteoRepo) PatchProductoGrupoConteo(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoGrupoConteo, error) {
	var item ProductoGrupoConteo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN productos ON productos.id = producto_grupo_conteos.producto_id").
			Where("producto_grupo_conteos.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("producto grupo conteo with id %d not found", id)
		}
		return item, fmt.Errorf("error patching producto grupo conteo with id %d: %w", id, err)
	}
	return item, nil
}
