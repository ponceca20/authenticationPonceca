package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// GrupoOfertaRepository define la interfaz para operaciones de persistencia de GrupoOferta
type GrupoOfertaRepository interface {
	GetAllGrupoOfertas(empresaID uint64) ([]GrupoOferta, error)
	GetGrupoOferta(id uint64, empresaID uint64) (GrupoOferta, error)
	CreateGrupoOferta(grupo *GrupoOferta) error
	UpdateGrupoOferta(id uint64, empresaID uint64, grupo *GrupoOferta) (GrupoOferta, error)
	DeleteGrupoOferta(id uint64, empresaID uint64) error
	SeedGrupoOfertas(grupos []GrupoOferta) error
	PatchGrupoOferta(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoOferta, error)
}

type grupoOfertaRepo struct {
	DB *gorm.DB
}

func NewGrupoOfertaRepository() GrupoOfertaRepository {
	return &grupoOfertaRepo{DB: database.DBconn}
}

func (r *grupoOfertaRepo) GetAllGrupoOfertas(empresaID uint64) ([]GrupoOferta, error) {
	var grupos []GrupoOferta
	err := r.DB.Joins("JOIN ofertas ON ofertas.id = grupo_oferta.oferta_id").
		Where("ofertas.empresa_id = ?", empresaID).
		Find(&grupos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar grupos de oferta: %w", err)
	}
	return grupos, nil
}

func (r *grupoOfertaRepo) GetGrupoOferta(id uint64, empresaID uint64) (GrupoOferta, error) {
	var grupo GrupoOferta
	err := r.DB.Joins("JOIN ofertas ON ofertas.id = grupo_oferta.oferta_id").
		Where("grupo_oferta.id = ? AND ofertas.empresa_id = ?", id, empresaID).
		First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grupo, fmt.Errorf("grupo de oferta con id %d no encontrado", id)
		}
		return grupo, fmt.Errorf("error al recuperar grupo de oferta con id %d: %w", id, err)
	}
	return grupo, nil
}

func (r *grupoOfertaRepo) CreateGrupoOferta(grupo *GrupoOferta) error {
	err := r.DB.Create(grupo).Error
	if err != nil {
		return fmt.Errorf("error al crear grupo de oferta: %w", err)
	}
	return nil
}

func (r *grupoOfertaRepo) UpdateGrupoOferta(id uint64, empresaID uint64, updatedGrupo *GrupoOferta) (GrupoOferta, error) {
	var grupo GrupoOferta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN ofertas ON ofertas.id = grupo_oferta.oferta_id").
			Where("grupo_oferta.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			First(&grupo).Error; err != nil {
			return err
		}
		return tx.Model(&grupo).Updates(updatedGrupo).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grupo, fmt.Errorf("grupo de oferta con id %d no encontrado", id)
		}
		return grupo, fmt.Errorf("error al actualizar grupo de oferta con id %d: %w", id, err)
	}
	return grupo, nil
}

func (r *grupoOfertaRepo) DeleteGrupoOferta(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Joins("JOIN ofertas ON ofertas.id = grupo_oferta.oferta_id").
			Where("grupo_oferta.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			Delete(&GrupoOferta{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("grupo de oferta con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar grupo de oferta: %w", err)
	}
	return nil
}

func (r *grupoOfertaRepo) SeedGrupoOfertas(grupos []GrupoOferta) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&grupos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar grupos de oferta: %w", err)
	}
	return nil
}

func (r *grupoOfertaRepo) PatchGrupoOferta(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoOferta, error) {
	var grupo GrupoOferta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN ofertas ON ofertas.id = grupo_oferta.oferta_id").
			Where("grupo_oferta.id = ? AND ofertas.empresa_id = ?", id, empresaID).
			First(&grupo).Error; err != nil {
			return err
		}
		return tx.Model(&grupo).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grupo, fmt.Errorf("grupo de oferta con id %d no encontrado", id)
		}
		return grupo, fmt.Errorf("error al actualizar parcialmente grupo de oferta con id %d: %w", id, err)
	}
	return grupo, nil
}
