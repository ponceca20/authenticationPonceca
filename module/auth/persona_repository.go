package auth

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PersonaRepository define la interfaz para operaciones de persistencia de Persona.
type PersonaRepository interface {
	GetAllPersonas() ([]Persona, error)
	GetPersona(id uint) (Persona, error) // Changed from uint64 to uint
	CreatePersona(persona *Persona) error
	UpdatePersona(id uint, updatedPersona *Persona) (Persona, error) // Changed from uint64 to uint
	DeletePersona(id uint) error                                     // Changed from uint64 to uint
	SeedPersonas(personas []Persona) error
	PatchPersona(id uint, fields map[string]interface{}) (Persona, error) // Changed from uint64 to uint
}

type personaRepo struct {
	DB *gorm.DB
}

func NewPersonaRepository() PersonaRepository {
	return &personaRepo{DB: database.DBconn}
}

func (r *personaRepo) GetAllPersonas() ([]Persona, error) {
	var personas []Persona
	//err := r.DB.Preload("DocumentoTipo").Find(&personas).Error
	err := r.DB.Find(&personas).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar personas: %w", err)
	}
	return personas, nil
}

func (r *personaRepo) GetPersona(id uint) (Persona, error) {
	var persona Persona
	// Cambiado Preload("DocumentoTipo") a Preload("TipoDocumento")
	err := r.DB.First(&persona, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return persona, fmt.Errorf("persona con id %d no encontrada", id)
		}
		return persona, fmt.Errorf("error al recuperar persona con id %d: %w", id, err)
	}
	return persona, nil
}

func (r *personaRepo) CreatePersona(persona *Persona) error {
	err := r.DB.Create(persona).Error
	if err != nil {
		return fmt.Errorf("error al crear persona: %w", err)
	}
	return nil
}

func (r *personaRepo) UpdatePersona(id uint, updatedPersona *Persona) (Persona, error) {
	var persona Persona
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Cambiado Preload("DocumentoTipo") a Preload("TipoDocumento")
		if err := tx.Preload("TipoDocumento").First(&persona, id).Error; err != nil {
			return err
		}
		return tx.Model(&persona).Updates(updatedPersona).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return persona, fmt.Errorf("persona con id %d no encontrada", id)
		}
		return persona, fmt.Errorf("error al actualizar persona con id %d: %w", id, err)
	}
	// Reload the persona to get updated data including TipoDocumento
	r.DB.Preload("TipoDocumento").First(&persona, id)
	return persona, nil
}

func (r *personaRepo) DeletePersona(id uint) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Persona{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("persona con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar persona: %w", err)
	}
	return nil
}

func (r *personaRepo) SeedPersonas(personas []Persona) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&personas).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar personas: %w", err)
	}
	return nil
}

func (r *personaRepo) PatchPersona(id uint, fields map[string]interface{}) (Persona, error) {
	var persona Persona
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&persona, id).Error; err != nil {
			return err
		}
		return tx.Model(&persona).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return persona, fmt.Errorf("persona con id %d no encontrada", id)
		}
		return persona, fmt.Errorf("error al actualizar parcialmente persona con id %d: %w", id, err)
	}
	// Reload the persona to get updated data including TipoDocumento only
	r.DB.Preload("TipoDocumento").First(&persona, id)
	return persona, nil
}
