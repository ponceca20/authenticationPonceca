package auth

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// TipoDocumentoSunatRepository define la interfaz para operaciones de persistencia de TipoDocumentoSunat.
type TipoDocumentoSunatRepository interface {
	GetAllTiposDocumento() ([]TipoDocumentoSunat, error)
	GetTipoDocumento(id uint) (TipoDocumentoSunat, error)
	CreateTipoDocumento(tipoDocumento *TipoDocumentoSunat) error
	UpdateTipoDocumento(id uint, updatedTipoDocumento *TipoDocumentoSunat) (TipoDocumentoSunat, error)
	DeleteTipoDocumento(id uint) error
	SeedTiposDocumento(tiposDocumento []TipoDocumentoSunat) error
	PatchTipoDocumento(id uint, fields map[string]interface{}) (TipoDocumentoSunat, error)
	GetTipoDocumentoByCodigo(codigo string) (TipoDocumentoSunat, error)
}

type tipoDocumentoRepo struct {
	DB *gorm.DB
}

func NewTipoDocumentoSunatRepository() TipoDocumentoSunatRepository {
	return &tipoDocumentoRepo{DB: database.DBconn}
}

func (r *tipoDocumentoRepo) GetAllTiposDocumento() ([]TipoDocumentoSunat, error) {
	var tiposDocumento []TipoDocumentoSunat
	err := r.DB.Find(&tiposDocumento).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar tipos de documento: %w", err)
	}
	return tiposDocumento, nil
}

func (r *tipoDocumentoRepo) GetTipoDocumento(id uint) (TipoDocumentoSunat, error) {
	var tipoDocumento TipoDocumentoSunat
	err := r.DB.First(&tipoDocumento, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tipoDocumento, fmt.Errorf("tipo de documento con id %d no encontrado", id)
		}
		return tipoDocumento, fmt.Errorf("error al recuperar tipo de documento con id %d: %w", id, err)
	}
	return tipoDocumento, nil
}

func (r *tipoDocumentoRepo) GetTipoDocumentoByCodigo(codigo string) (TipoDocumentoSunat, error) {
	var tipoDocumento TipoDocumentoSunat
	err := r.DB.Where("codigo = ?", codigo).First(&tipoDocumento).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tipoDocumento, fmt.Errorf("tipo de documento con código %s no encontrado", codigo)
		}
		return tipoDocumento, fmt.Errorf("error al recuperar tipo de documento con código %s: %w", codigo, err)
	}
	return tipoDocumento, nil
}

func (r *tipoDocumentoRepo) CreateTipoDocumento(tipoDocumento *TipoDocumentoSunat) error {
	err := r.DB.Create(tipoDocumento).Error
	if err != nil {
		return fmt.Errorf("error al crear tipo de documento: %w", err)
	}
	return nil
}

func (r *tipoDocumentoRepo) UpdateTipoDocumento(id uint, updatedTipoDocumento *TipoDocumentoSunat) (TipoDocumentoSunat, error) {
	var tipoDocumento TipoDocumentoSunat
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&tipoDocumento, id).Error; err != nil {
			return err
		}
		return tx.Model(&tipoDocumento).Updates(updatedTipoDocumento).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tipoDocumento, fmt.Errorf("tipo de documento con id %d no encontrado", id)
		}
		return tipoDocumento, fmt.Errorf("error al actualizar tipo de documento con id %d: %w", id, err)
	}
	// Reload the document type to get updated data
	r.DB.First(&tipoDocumento, id)
	return tipoDocumento, nil
}

func (r *tipoDocumentoRepo) DeleteTipoDocumento(id uint) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&TipoDocumentoSunat{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("tipo de documento con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar tipo de documento: %w", err)
	}
	return nil
}

func (r *tipoDocumentoRepo) SeedTiposDocumento(tiposDocumento []TipoDocumentoSunat) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&tiposDocumento).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar tipos de documento: %w", err)
	}
	return nil
}

func (r *tipoDocumentoRepo) PatchTipoDocumento(id uint, fields map[string]interface{}) (TipoDocumentoSunat, error) {
	var tipoDocumento TipoDocumentoSunat
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&tipoDocumento, id).Error; err != nil {
			return err
		}
		return tx.Model(&tipoDocumento).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tipoDocumento, fmt.Errorf("tipo de documento con id %d no encontrado", id)
		}
		return tipoDocumento, fmt.Errorf("error al actualizar parcialmente tipo de documento con id %d: %w", id, err)
	}
	// Reload the document type to get updated data
	r.DB.First(&tipoDocumento, id)
	return tipoDocumento, nil
}
