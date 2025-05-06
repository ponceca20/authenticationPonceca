package product

import (
	"errors"
	"fmt"

	"practicev2/database"

	"gorm.io/gorm"
)

// GrupoConteoRepository define operaciones de persistencia para GrupoConteo.
type GrupoConteoRepository interface {
	GetAll(empresaID uint64) ([]GrupoConteo, error)
	GetByID(id uint64, empresaID uint64) (GrupoConteo, error)
	Create(item *GrupoConteo) error
	Update(id uint64, empresaID uint64, item *GrupoConteo) (GrupoConteo, error)
	Delete(id uint64, empresaID uint64) error
	Patch(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoConteo, error)
}

type grupoConteoRepo struct {
	DB *gorm.DB
}

// NewGrupoConteoRepository crea una instancia del repositorio.
func NewGrupoConteoRepository() GrupoConteoRepository {
	return &grupoConteoRepo{DB: database.DBconn}
}

func (r *grupoConteoRepo) GetAll(empresaID uint64) ([]GrupoConteo, error) {
	var items []GrupoConteo
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar grupo de conteo: %w", err)
	}
	return items, nil
}

func (r *grupoConteoRepo) GetByID(id uint64, empresaID uint64) (GrupoConteo, error) {
	var item GrupoConteo
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("grupo de conteo con id %d no encontrada", id)
		}
		return item, fmt.Errorf("error al recuperar grupo de conteo con id %d: %w", id, err)
	}
	return item, nil
}

func (r *grupoConteoRepo) Create(item *GrupoConteo) error {
	if err := r.DB.Create(item).Error; err != nil {
		return fmt.Errorf("error al crear grupo de conteo: %w", err)
	}
	return nil
}

func (r *grupoConteoRepo) Update(id uint64, empresaID uint64, item *GrupoConteo) (GrupoConteo, error) {
	var existing GrupoConteo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&existing).Error; err != nil {
			return err
		}
		return tx.Model(&existing).Updates(item).Error
	})
	if err != nil {
		return existing, fmt.Errorf("error al actualizar grupo de conteo con id %d: %w", id, err)
	}
	return existing, nil
}

func (r *grupoConteoRepo) Delete(id uint64, empresaID uint64) error {
	result := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&GrupoConteo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("grupo de conteo con id %d no encontrado", id)
	}
	return nil
}

func (r *grupoConteoRepo) Patch(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoConteo, error) {
	var item GrupoConteo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&item).Error; err != nil {
			return err
		}
		return tx.Model(&item).Updates(fields).Error
	})
	if err != nil {
		return item, fmt.Errorf("error al actualizar parcialmente grupo de conteo con id %d: %w", id, err)
	}
	return item, nil
}
