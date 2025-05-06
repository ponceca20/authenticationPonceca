package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ListaPreciosRepository define la interfaz para operaciones de persistencia de ListaPrecios
type ListaPreciosRepository interface {
	GetAllListaPrecios(empresaID uint64) ([]ListaPrecios, error)
	GetListaPreciosWithPresentaciones(id uint64, empresaID uint64) (ListaPrecios, error)
	GetListaPrecios(id uint64, empresaID uint64) (ListaPrecios, error)
	CreateListaPrecios(listaPrecios *ListaPrecios) error
	UpdateListaPrecios(id uint64, empresaID uint64, listaPrecios *ListaPrecios) (ListaPrecios, error)
	DeleteListaPrecios(id uint64, empresaID uint64) error
	SeedListaPrecios(listaPrecios []ListaPrecios) error
	PatchListaPrecios(id uint64, empresaID uint64, fields map[string]interface{}) (ListaPrecios, error)
}

type listaPreciosRepo struct {
	DB *gorm.DB
}

// NewListaPreciosRepository crea una instancia del repositorio
func NewListaPreciosRepository() ListaPreciosRepository {
	return &listaPreciosRepo{DB: database.DBconn}
}

// GetAllListaPrecios obtiene todas las listas de precios para una empresa
func (r *listaPreciosRepo) GetAllListaPrecios(empresaID uint64) ([]ListaPrecios, error) {
	var listaPrecios []ListaPrecios
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&listaPrecios).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar listas de precios: %w", err)
	}
	return listaPrecios, nil
}

// GetListaPreciosWithPresentaciones obtiene una lista de precios con sus presentaciones asociadas
func (r *listaPreciosRepo) GetListaPreciosWithPresentaciones(id uint64, empresaID uint64) (ListaPrecios, error) {
	var listaPrecios ListaPrecios
	err := r.DB.
		Preload("ListaPresentacion").
		Preload("ListaPresentacion.Presentacion").
		Preload("ListaPresentacion.Presentacion.Producto").
		Where("id = ? AND empresa_id = ?", id, empresaID).
		First(&listaPrecios).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return listaPrecios, fmt.Errorf("lista de precios con id %d no encontrada", id)
		}
		return listaPrecios, fmt.Errorf("error al recuperar lista de precios con id %d: %w", id, err)
	}
	return listaPrecios, nil
}

// GetListaPrecios obtiene una lista de precios específica
func (r *listaPreciosRepo) GetListaPrecios(id uint64, empresaID uint64) (ListaPrecios, error) {
	var listaPrecios ListaPrecios
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&listaPrecios).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return listaPrecios, fmt.Errorf("lista de precios con id %d no encontrada", id)
		}
		return listaPrecios, fmt.Errorf("error al recuperar lista de precios con id %d: %w", id, err)
	}
	return listaPrecios, nil
}

// CreateListaPrecios crea una nueva lista de precios
func (r *listaPreciosRepo) CreateListaPrecios(listaPrecios *ListaPrecios) error {
	err := r.DB.Create(listaPrecios).Error
	if err != nil {
		return fmt.Errorf("error al crear lista de precios: %w", err)
	}
	return nil
}

// UpdateListaPrecios actualiza una lista de precios existente
func (r *listaPreciosRepo) UpdateListaPrecios(id uint64, empresaID uint64, updatedListaPrecios *ListaPrecios) (ListaPrecios, error) {
	var listaPrecios ListaPrecios
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&listaPrecios).Error; err != nil {
			return err
		}
		return tx.Model(&listaPrecios).Updates(updatedListaPrecios).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return listaPrecios, fmt.Errorf("lista de precios con id %d no encontrada", id)
		}
		return listaPrecios, fmt.Errorf("error al actualizar lista de precios con id %d: %w", id, err)
	}
	return listaPrecios, nil
}

// DeleteListaPrecios elimina una lista de precios
func (r *listaPreciosRepo) DeleteListaPrecios(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si hay presentaciones asociadas a esta lista de precios
		var count int64
		if err := tx.Model(&ListaPreciosPresentacion{}).Where("lista_precios_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar la lista de precios porque tiene %d presentaciones asociadas", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&ListaPrecios{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("lista de precios con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar lista de precios: %w", err)
	}
	return nil
}

// SeedListaPrecios inicializa múltiples listas de precios
func (r *listaPreciosRepo) SeedListaPrecios(listaPrecios []ListaPrecios) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&listaPrecios).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar listas de precios: %w", err)
	}
	return nil
}

// PatchListaPrecios actualiza parcialmente una lista de precios
func (r *listaPreciosRepo) PatchListaPrecios(id uint64, empresaID uint64, fields map[string]interface{}) (ListaPrecios, error) {
	var listaPrecios ListaPrecios
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&listaPrecios).Error; err != nil {
			return err
		}
		return tx.Model(&listaPrecios).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return listaPrecios, fmt.Errorf("lista de precios con id %d no encontrada", id)
		}
		return listaPrecios, fmt.Errorf("error al actualizar parcialmente lista de precios con id %d: %w", id, err)
	}
	return listaPrecios, nil
}
