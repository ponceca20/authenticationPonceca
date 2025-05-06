package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// OfertaRepository define la interfaz para operaciones de persistencia de Oferta
type OfertaRepository interface {
	GetAllOfertas(empresaID uint64) ([]Oferta, error)
	GetOferta(id uint64, empresaID uint64) (Oferta, error)
	CreateOferta(oferta *Oferta) error
	UpdateOferta(id uint64, empresaID uint64, oferta *Oferta) (Oferta, error)
	DeleteOferta(id uint64, empresaID uint64) error
	SeedOfertas(ofertas []Oferta) error
	PatchOferta(id uint64, empresaID uint64, fields map[string]interface{}) (Oferta, error)
}

type ofertaRepo struct {
	DB *gorm.DB
}

// NewOfertaRepository crea una instancia del repositorio de Oferta
func NewOfertaRepository() OfertaRepository {
	return &ofertaRepo{DB: database.DBconn}
}

func (r *ofertaRepo) GetAllOfertas(empresaID uint64) ([]Oferta, error) {
	var ofertas []Oferta
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&ofertas).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar ofertas: %w", err)
	}
	return ofertas, nil
}

func (r *ofertaRepo) GetOferta(id uint64, empresaID uint64) (Oferta, error) {
	var oferta Oferta
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&oferta).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return oferta, fmt.Errorf("oferta con id %d no encontrada", id)
		}
		return oferta, fmt.Errorf("error al recuperar oferta con id %d: %w", id, err)
	}
	return oferta, nil
}

func (r *ofertaRepo) CreateOferta(oferta *Oferta) error {
	err := r.DB.Create(oferta).Error
	if err != nil {
		return fmt.Errorf("error al crear oferta: %w", err)
	}
	return nil
}

func (r *ofertaRepo) UpdateOferta(id uint64, empresaID uint64, updatedOferta *Oferta) (Oferta, error) {
	var oferta Oferta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&oferta).Error; err != nil {
			return err
		}
		return tx.Model(&oferta).Updates(updatedOferta).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return oferta, fmt.Errorf("oferta con id %d no encontrada", id)
		}
		return oferta, fmt.Errorf("error al actualizar oferta con id %d: %w", id, err)
	}
	return oferta, nil
}

func (r *ofertaRepo) DeleteOferta(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Oferta{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("oferta con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar oferta: %w", err)
	}
	return nil
}

func (r *ofertaRepo) SeedOfertas(ofertas []Oferta) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ofertas).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar ofertas: %w", err)
	}
	return nil
}

func (r *ofertaRepo) PatchOferta(id uint64, empresaID uint64, fields map[string]interface{}) (Oferta, error) {
	var oferta Oferta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&oferta).Error; err != nil {
			return err
		}
		return tx.Model(&oferta).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return oferta, fmt.Errorf("oferta con id %d no encontrada", id)
		}
		return oferta, fmt.Errorf("error al actualizar parcialmente oferta con id %d: %w", id, err)
	}
	return oferta, nil
}
