package product

import (
	"errors"
	"fmt"

	"practicev2/database"

	"gorm.io/gorm"
)

// ListaPreciosPresentacionRepository define métodos para ListaPreciosPresentacion.
type ListaPreciosPresentacionRepository interface {
	GetAll() ([]ListaPreciosPresentacion, error)
	GetByID(id uint64) (ListaPreciosPresentacion, error)
	Create(item *ListaPreciosPresentacion) error
	Update(id uint64, item *ListaPreciosPresentacion) (ListaPreciosPresentacion, error)
	Delete(id uint64) error
	Patch(id uint64, fields map[string]interface{}) (ListaPreciosPresentacion, error)
}

type listaPreciosPresentacionRepo struct {
	DB *gorm.DB
}

// NewListaPreciosPresentacionRepository crea una nueva instancia del repositorio.
func NewListaPreciosPresentacionRepository() ListaPreciosPresentacionRepository {
	return &listaPreciosPresentacionRepo{DB: database.DBconn}
}

func (r *listaPreciosPresentacionRepo) GetAll() ([]ListaPreciosPresentacion, error) {
	var items []ListaPreciosPresentacion
	err := r.DB.Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar lista de presentaciones: %w", err)
	}
	return items, nil
}

func (r *listaPreciosPresentacionRepo) GetByID(id uint64) (ListaPreciosPresentacion, error) {
	var item ListaPreciosPresentacion
	err := r.DB.Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("presentación con id %d no encontrada", id)
		}
		return item, fmt.Errorf("error al recuperar presentación con id %d: %w", id, err)
	}
	return item, nil
}

func (r *listaPreciosPresentacionRepo) Create(item *ListaPreciosPresentacion) error {
	if err := r.DB.Create(item).Error; err != nil {
		return fmt.Errorf("error al crear presentación: %w", err)
	}
	return nil
}

func (r *listaPreciosPresentacionRepo) Update(id uint64, item *ListaPreciosPresentacion) (ListaPreciosPresentacion, error) {
	var existing ListaPreciosPresentacion
	err := r.DB.Where("id = ?", id).First(&existing).Error
	if err != nil {
		return existing, fmt.Errorf("presentación con id %d no encontrada: %w", id, err)
	}
	if err := r.DB.Model(&existing).Updates(item).Error; err != nil {
		return existing, fmt.Errorf("error al actualizar presentación con id %d: %w", id, err)
	}
	return existing, nil
}

func (r *listaPreciosPresentacionRepo) Delete(id uint64) error {
	result := r.DB.Where("id = ?", id).Delete(&ListaPreciosPresentacion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("presentación con id %d no encontrada", id)
	}
	return nil
}

func (r *listaPreciosPresentacionRepo) Patch(id uint64, fields map[string]interface{}) (ListaPreciosPresentacion, error) {
	var item ListaPreciosPresentacion
	delete(fields, "empresa_id")
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar parcialmente presentación con id %d: %w", id, err)
	}
	return item, nil
}
