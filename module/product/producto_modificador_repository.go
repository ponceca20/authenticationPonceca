package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

type ProductoModificadorRepository interface {
	GetAllProductoModificadores(productoID uint64) ([]ProductoModificador, error)
	GetProductoModificador(id uint64) (ProductoModificador, error)
	CreateProductoModificador(pm *ProductoModificador) error
	UpdateProductoModificador(id uint64, pm *ProductoModificador) (ProductoModificador, error)
	DeleteProductoModificador(id uint64) error
	SeedProductoModificadores(pms []ProductoModificador) error
	PatchProductoModificador(id uint64, fields map[string]interface{}) (ProductoModificador, error)
}

type productoModificadorRepo struct {
	DB *gorm.DB
}

func NewProductoModificadorRepository() ProductoModificadorRepository {
	return &productoModificadorRepo{DB: database.DBconn}
}

func (r *productoModificadorRepo) GetAllProductoModificadores(productoID uint64) ([]ProductoModificador, error) {
	var pms []ProductoModificador
	var err error
	if productoID == 0 {
		err = r.DB.Find(&pms).Error // Todos los modificadores
	} else {
		err = r.DB.Where("producto_id = ?", productoID).Find(&pms).Error
	}
	if err != nil {
		return nil, fmt.Errorf("error al recuperar modificadores de producto: %w", err)
	}
	return pms, nil
}

func (r *productoModificadorRepo) GetProductoModificador(id uint64) (ProductoModificador, error) {
	var pm ProductoModificador
	err := r.DB.Where("id = ?", id).First(&pm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pm, fmt.Errorf("modificador de producto con id %d no encontrado", id)
		}
		return pm, fmt.Errorf("error al recuperar modificador de producto con id %d: %w", id, err)
	}
	return pm, nil
}

func (r *productoModificadorRepo) CreateProductoModificador(pm *ProductoModificador) error {
	// Verificar si el producto existe
	var count int64
	if err := r.DB.Model(&Producto{}).Where("id = ?", pm.ProductoID).Count(&count).Error; err != nil {
		return fmt.Errorf("error al verificar el producto: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("el producto con ID %d no existe", pm.ProductoID)
	}
	return r.DB.Create(pm).Error
}

func (r *productoModificadorRepo) UpdateProductoModificador(id uint64, updatedPM *ProductoModificador) (ProductoModificador, error) {
	var pm ProductoModificador
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&pm).Error; err != nil {
			return err
		}
		// Verificar si el producto existe si se está actualizando
		if updatedPM.ProductoID != pm.ProductoID {
			var count int64
			if err := tx.Model(&Producto{}).Where("id = ?", updatedPM.ProductoID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar el producto: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("el producto con ID %d no existe", updatedPM.ProductoID)
			}
		}
		return tx.Model(&pm).Updates(updatedPM).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pm, fmt.Errorf("modificador de producto con id %d no encontrado", id)
		}
		return pm, fmt.Errorf("error al actualizar modificador de producto con id %d: %w", id, err)
	}
	return pm, nil
}

func (r *productoModificadorRepo) DeleteProductoModificador(id uint64) error {
	result := r.DB.Delete(&ProductoModificador{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar modificador de producto: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("modificador de producto con id %d no encontrado", id)
	}
	return nil
}

func (r *productoModificadorRepo) SeedProductoModificadores(pms []ProductoModificador) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&pms).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar modificadores de producto: %w", err)
	}
	return nil
}

func (r *productoModificadorRepo) PatchProductoModificador(id uint64, fields map[string]interface{}) (ProductoModificador, error) {
	var pm ProductoModificador
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&pm).Error; err != nil {
			return err
		}
		if productoID, ok := fields["producto_id"].(uint64); ok && productoID != pm.ProductoID {
			var count int64
			if err := tx.Model(&Producto{}).Where("id = ?", productoID).Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar el producto: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("el producto con ID %d no existe", productoID)
			}
		}
		return tx.Model(&pm).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pm, fmt.Errorf("modificador de producto con id %d no encontrado", id)
		}
		return pm, fmt.Errorf("error al actualizar parcialmente modificador de producto con id %d: %w", id, err)
	}
	return pm, nil
}
