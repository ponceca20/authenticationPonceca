package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// OfertaPresentacionRepository define la interfaz para operaciones de persistencia de OfertaPresentacion
type OfertaPresentacionRepository interface {
	GetAllOfertaPresentaciones(empresaID uint64) ([]OfertaPresentacion, error)
	GetOfertaPresentacion(id uint64, empresaID uint64) (OfertaPresentacion, error)
	CreateOfertaPresentacion(op *OfertaPresentacion) error
	UpdateOfertaPresentacion(id uint64, empresaID uint64, op *OfertaPresentacion) (OfertaPresentacion, error)
	DeleteOfertaPresentacion(id uint64, empresaID uint64) error
	SeedOfertaPresentaciones(ops []OfertaPresentacion) error
	PatchOfertaPresentacion(id uint64, empresaID uint64, fields map[string]interface{}) (OfertaPresentacion, error)
}

type ofertaPresentacionRepo struct {
	DB *gorm.DB
}

func NewOfertaPresentacionRepository() OfertaPresentacionRepository {
	return &ofertaPresentacionRepo{DB: database.DBconn}
}

func (r *ofertaPresentacionRepo) GetAllOfertaPresentaciones(empresaID uint64) ([]OfertaPresentacion, error) {
	var ops []OfertaPresentacion
	err := r.DB.Joins("JOIN ofertas ON ofertas.id = oferta_presentacion.oferta_id").
		Where("ofertas.empresa_id = ?", empresaID).
		Find(&ops).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar ofertas presentaciones: %w", err)
	}
	return ops, nil
}

func (r *ofertaPresentacionRepo) GetOfertaPresentacion(id uint64, empresaID uint64) (OfertaPresentacion, error) {
	var op OfertaPresentacion
	err := r.DB.Joins("JOIN ofertas ON ofertas.id = oferta_presentacion.oferta_id").
		Where("oferta_presentacion.id = ? AND ofertas.empresa_id = ?", id, empresaID).
		First(&op).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return op, fmt.Errorf("oferta presentación con id %d no encontrada", id)
		}
		return op, fmt.Errorf("error al recuperar oferta presentación con id %d: %w", id, err)
	}
	return op, nil
}

func (r *ofertaPresentacionRepo) CreateOfertaPresentacion(op *OfertaPresentacion) error {
	err := r.DB.Create(op).Error
	if err != nil {
		return fmt.Errorf("error al crear oferta presentación: %w", err)
	}
	return nil
}

func (r *ofertaPresentacionRepo) UpdateOfertaPresentacion(id uint64, empresaID uint64, updatedOP *OfertaPresentacion) (OfertaPresentacion, error) {
	var op OfertaPresentacion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN ofertas ON ofertas.id = oferta_presentacion.oferta_id").
			Where("oferta_presentacion.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			First(&op).Error; err != nil {
			return err
		}
		return tx.Model(&op).Updates(updatedOP).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return op, fmt.Errorf("oferta presentación con id %d no encontrada", id)
		}
		return op, fmt.Errorf("error al actualizar oferta presentación con id %d: %w", id, err)
	}
	return op, nil
}

func (r *ofertaPresentacionRepo) DeleteOfertaPresentacion(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Joins("JOIN ofertas ON ofertas.id = oferta_presentacion.oferta_id").
			Where("oferta_presentacion.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			Delete(&OfertaPresentacion{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("oferta presentación con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar oferta presentación: %w", err)
	}
	return nil
}

func (r *ofertaPresentacionRepo) SeedOfertaPresentaciones(ops []OfertaPresentacion) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ops).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar ofertas presentaciones: %w", err)
	}
	return nil
}

func (r *ofertaPresentacionRepo) PatchOfertaPresentacion(id uint64, empresaID uint64, fields map[string]interface{}) (OfertaPresentacion, error) {
	var op OfertaPresentacion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN ofertas ON ofertas.id = oferta_presentacion.oferta_id").
			Where("oferta_presentacion.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			First(&op).Error; err != nil {
			return err
		}
		return tx.Model(&op).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return op, fmt.Errorf("oferta presentación con id %d no encontrada", id)
		}
		return op, fmt.Errorf("error al actualizar parcialmente oferta presentación con id %d: %w", id, err)
	}
	return op, nil
}
