package organization_config

import (
	"practicev2/module/authentication/models"

	"gorm.io/gorm"
)

//=============================================================================
// REPOSITORIO PARA CONFIGURACIÓN DE MÓDULOS ORGANIZACIONALES
//=============================================================================

// OrganizationConfigRepository interfaz del repositorio
type OrganizationConfigRepository interface {
	Create(config *models.OrganizationModuleConfig) error
	GetByID(id string) (*models.OrganizationModuleConfig, error)
	GetByOrganizationID(orgID string) ([]models.OrganizationModuleConfig, error)
	GetEnabledByOrganizationID(orgID string) ([]models.OrganizationModuleConfig, error)
	GetByModuleName(orgID, moduleName string) (*models.OrganizationModuleConfig, error)
	UpdateModuleStatus(orgID, moduleName string, isEnabled bool) error
	Delete(id string) error
	List(limit, offset int) ([]models.OrganizationModuleConfig, error)
}

// organizationConfigRepository implementación del repositorio
type organizationConfigRepository struct {
	db *gorm.DB
}

// NewOrganizationConfigRepository crea una nueva instancia del repositorio
func NewOrganizationConfigRepository(db *gorm.DB) OrganizationConfigRepository {
	return &organizationConfigRepository{db: db}
}

//=============================================================================
// IMPLEMENTACIÓN DE MÉTODOS DEL REPOSITORIO
//=============================================================================

// Create crea una nueva configuración de módulo
func (r *organizationConfigRepository) Create(config *models.OrganizationModuleConfig) error {
	return r.db.Create(config).Error
}

// GetByID obtiene una configuración por ID
func (r *organizationConfigRepository) GetByID(id string) (*models.OrganizationModuleConfig, error) {
	var config models.OrganizationModuleConfig
	err := r.db.Preload("Organization").Preload("Installer").
		Where("id = ?", id).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByOrganizationID obtiene todas las configuraciones de una organización
func (r *organizationConfigRepository) GetByOrganizationID(orgID string) ([]models.OrganizationModuleConfig, error) {
	var configs []models.OrganizationModuleConfig
	err := r.db.Preload("Organization").
		Where("organization_id = ?", orgID).
		Order("module_name ASC").
		Find(&configs).Error
	return configs, err
}

// GetEnabledByOrganizationID obtiene solo las configuraciones habilitadas de una organización
func (r *organizationConfigRepository) GetEnabledByOrganizationID(orgID string) ([]models.OrganizationModuleConfig, error) {
	var configs []models.OrganizationModuleConfig
	err := r.db.Preload("Organization").
		Where("organization_id = ? AND is_enabled = ?", orgID, true).
		Order("module_name ASC").
		Find(&configs).Error
	return configs, err
}

// GetByModuleName obtiene la configuración de un módulo específico de una organización
func (r *organizationConfigRepository) GetByModuleName(orgID, moduleName string) (*models.OrganizationModuleConfig, error) {
	var config models.OrganizationModuleConfig
	err := r.db.Preload("Organization").Preload("Installer").
		Where("organization_id = ? AND module_name = ?", orgID, moduleName).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateModuleStatus actualiza el estado habilitado/deshabilitado de un módulo
func (r *organizationConfigRepository) UpdateModuleStatus(orgID, moduleName string, isEnabled bool) error {
	return r.db.Model(&models.OrganizationModuleConfig{}).
		Where("organization_id = ? AND module_name = ?", orgID, moduleName).
		Update("is_enabled", isEnabled).Error
}

// Delete elimina una configuración de módulo
func (r *organizationConfigRepository) Delete(id string) error {
	return r.db.Delete(&models.OrganizationModuleConfig{}, "id = ?", id).Error
}

// List obtiene todas las configuraciones con paginación
func (r *organizationConfigRepository) List(limit, offset int) ([]models.OrganizationModuleConfig, error) {
	var configs []models.OrganizationModuleConfig
	query := r.db.Preload("Organization").Preload("Installer")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&configs).Error
	return configs, err
}
