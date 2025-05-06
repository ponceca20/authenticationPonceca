package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ProductoRepository define las operaciones de persistencia para Producto
type ProductoRepository interface {
	GetAllProductos(empresaID uint64) ([]Producto, error)
	GetProducto(id uint64, empresaID uint64) (Producto, error)
	CreateProducto(p *Producto) error
	UpdateProducto(id uint64, empresaID uint64, p *Producto) (Producto, error)
	DeleteProducto(id uint64, empresaID uint64) error
	PatchProducto(id uint64, empresaID uint64, fields map[string]interface{}) (Producto, error)
}

type productoRepo struct {
	DB *gorm.DB
}

// NewProductoRepository crea una instancia del repositorio
func NewProductoRepository() ProductoRepository {
	return &productoRepo{DB: database.DBconn}
}

// GetAllProductos obtiene todos los productos de una empresa
func (r *productoRepo) GetAllProductos(empresaID uint64) ([]Producto, error) {
	var items []Producto
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar productos: %w", err)
	}
	return items, nil
}

// GetProducto obtiene un producto específico
func (r *productoRepo) GetProducto(id uint64, empresaID uint64) (Producto, error) {
	var item Producto
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("producto con id %d no encontrado", id)
		}
		return item, fmt.Errorf("error al recuperar producto %d: %w", id, err)
	}
	return item, nil
}

// CreateProducto crea un nuevo producto
func (r *productoRepo) CreateProducto(p *Producto) error {
	if err := r.DB.Create(p).Error; err != nil {
		return fmt.Errorf("error al crear producto: %w", err)
	}
	return nil
}

// UpdateProducto actualiza un producto existente
func (r *productoRepo) UpdateProducto(id uint64, empresaID uint64, updated *Producto) (Producto, error) {
	var item Producto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(updated).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar producto %d: %w", id, err)
	}
	return item, nil
}

// DeleteProducto elimina un producto
func (r *productoRepo) DeleteProducto(id uint64, empresaID uint64) error {
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Producto{}).Error
	if err != nil {
		return fmt.Errorf("error al eliminar producto: %w", err)
	}
	return nil
}

// PatchProducto actualiza parcialmente un producto
func (r *productoRepo) PatchProducto(id uint64, empresaID uint64, fields map[string]interface{}) (Producto, error) {
	var item Producto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar parcialmente producto %d: %w", id, err)
	}
	return item, nil
}
