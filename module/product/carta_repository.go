package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// CartaRepository define la interfaz para operaciones de persistencia de Carta
type CartaRepository interface {
	GetAllCartas(empresaID uint64) ([]Carta, error)
	GetCarta(id uint64, empresaID uint64) (Carta, error)
	CreateCarta(carta *Carta) error
	UpdateCarta(id uint64, empresaID uint64, carta *Carta) (Carta, error)
	DeleteCarta(id uint64, empresaID uint64) error
	SeedCartas(cartas []Carta) error
	PatchCarta(id uint64, empresaID uint64, fields map[string]interface{}) (Carta, error)
}

type cartaRepo struct {
	DB *gorm.DB
}

// NewCartaRepository crea una instancia del repositorio
func NewCartaRepository() CartaRepository {
	return &cartaRepo{DB: database.DBconn}
}

// GetAllCartas obtiene todas las cartas para una empresa
func (r *cartaRepo) GetAllCartas(empresaID uint64) ([]Carta, error) {
	var cartas []Carta
	err := r.DB.Preload("Familias").Where("empresa_id = ?", empresaID).Find(&cartas).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar cartas: %w", err)
	}
	return cartas, nil
}

// GetCarta obtiene una carta específica
func (r *cartaRepo) GetCarta(id uint64, empresaID uint64) (Carta, error) {
	var carta Carta
	err := r.DB.Preload("Familias").Where("id = ? AND empresa_id = ?", id, empresaID).First(&carta).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return carta, fmt.Errorf("carta con id %d no encontrada", id)
		}
		return carta, fmt.Errorf("error al recuperar carta con id %d: %w", id, err)
	}
	return carta, nil
}

// CreateCarta crea una nueva carta
func (r *cartaRepo) CreateCarta(carta *Carta) error {
	err := r.DB.Create(carta).Error
	if err != nil {
		return fmt.Errorf("error al crear carta: %w", err)
	}
	return nil
}

// UpdateCarta actualiza una carta existente
func (r *cartaRepo) UpdateCarta(id uint64, empresaID uint64, updatedCarta *Carta) (Carta, error) {
	var carta Carta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&carta).Error; err != nil {
			return err
		}
		return tx.Model(&carta).Updates(updatedCarta).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return carta, fmt.Errorf("carta con id %d no encontrada", id)
		}
		return carta, fmt.Errorf("error al actualizar carta con id %d: %w", id, err)
	}
	return carta, nil
}

// DeleteCarta elimina una carta
func (r *cartaRepo) DeleteCarta(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar si hay familias asociadas
		var count int64
		if err := tx.Model(&Familia{}).Where("carta_id = ?", id).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("no se puede eliminar la carta porque tiene %d familias asociadas", count)
		}

		result := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Carta{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("carta con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar carta: %w", err)
	}
	return nil
}

// SeedCartas inicializa múltiples cartas
func (r *cartaRepo) SeedCartas(cartas []Carta) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&cartas).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar cartas: %w", err)
	}
	return nil
}

// PatchCarta actualiza parcialmente una carta
func (r *cartaRepo) PatchCarta(id uint64, empresaID uint64, fields map[string]interface{}) (Carta, error) {
	var carta Carta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&carta).Error; err != nil {
			return err
		}
		return tx.Model(&carta).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return carta, fmt.Errorf("carta con id %d no encontrada", id)
		}
		return carta, fmt.Errorf("error al actualizar parcialmente carta con id %d: %w", id, err)
	}
	return carta, nil
}
