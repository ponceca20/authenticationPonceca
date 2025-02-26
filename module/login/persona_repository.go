package users

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PersonaRepository define la interfaz para operaciones de persistencia de Persona.
type PersonaRepository interface {
	GetAllPersonas() ([]Persona, error)
	GetPersona(id uint64) (Persona, error)
	CreatePersona(persona *Persona) error
	UpdatePersona(id uint64, updatedPersona *Persona) (Persona, error)
	DeletePersona(id uint64) error
	SeedPersonas(personas []Persona) error
	PatchPersona(id uint64, fields map[string]interface{}) (Persona, error)
}

type personaRepo struct {
	DB *gorm.DB
}

func NewPersonaRepository() PersonaRepository {
	return &personaRepo{DB: database.DBconn}
}

func (r *personaRepo) GetAllPersonas() ([]Persona, error) {
	var personas []Persona
	err := r.DB.Find(&personas).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar personas: %w", err)
	}
	return personas, nil
}

func (r *personaRepo) GetPersona(id uint64) (Persona, error) {
	var persona Persona
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

func (r *personaRepo) UpdatePersona(id uint64, updatedPersona *Persona) (Persona, error) {
	var persona Persona
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&persona, id).Error; err != nil {
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
	return persona, nil
}

func (r *personaRepo) DeletePersona(id uint64) error {
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

func (r *personaRepo) PatchPersona(id uint64, fields map[string]interface{}) (Persona, error) {
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
	return persona, nil
}
