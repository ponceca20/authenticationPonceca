package organization_config

import (
	"encoding/json"
	"practicev2/module/authentication/models"
)

//=============================================================================
// DTOs PARA CONFIGURACIÓN DE MÓDULOS ORGANIZACIONALES
//=============================================================================

// CreateModuleConfigRequest DTO para crear configuración de módulo
type CreateModuleConfigRequest struct {
	ModuleName   string                 `json:"module_name" validate:"required,min=2,max=100"`
	IsEnabled    bool                   `json:"is_enabled"`
	Resources    map[string][]string    `json:"resources" validate:"required"`
	DefaultRoles map[string]RoleConfig  `json:"default_roles"`
	Settings     map[string]interface{} `json:"settings"`
	Version      string                 `json:"version"`
}

// UpdateModuleConfigRequest DTO para actualizar configuración de módulo
type UpdateModuleConfigRequest struct {
	IsEnabled    *bool                  `json:"is_enabled,omitempty"`
	Resources    map[string][]string    `json:"resources,omitempty"`
	DefaultRoles map[string]RoleConfig  `json:"default_roles,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
	Version      string                 `json:"version,omitempty"`
}

// ModuleConfigResponse DTO para respuesta de configuración de módulo
type ModuleConfigResponse struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organization_id"`
	ModuleName     string                 `json:"module_name"`
	IsEnabled      bool                   `json:"is_enabled"`
	Resources      map[string][]string    `json:"resources"`
	DefaultRoles   map[string]RoleConfig  `json:"default_roles"`
	Settings       map[string]interface{} `json:"settings"`
	Version        string                 `json:"version"`
	InstalledBy    string                 `json:"installed_by"`
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`

	// Relaciones expandidas
	Organization *OrganizationInfo `json:"organization,omitempty"`
	Installer    *UserInfo         `json:"installer,omitempty"`
}

// OrganizationInfo información básica de organización
type OrganizationInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Type string `json:"type"`
}

// UserInfo información básica de usuario
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AvailableModuleResponse DTO para módulos disponibles
type AvailableModuleResponse struct {
	ModuleName  string              `json:"module_name"`
	Description string              `json:"description"`
	Version     string              `json:"version"`
	Resources   map[string][]string `json:"resources"`
	Category    string              `json:"category"`
}

// ModuleInstallationRequest DTO para instalar módulo predefinido
type ModuleInstallationRequest struct {
	ModuleName      string                 `json:"module_name" validate:"required"`
	CustomSettings  map[string]interface{} `json:"custom_settings"`
	EnabledFeatures []string               `json:"enabled_features"`
}

//=============================================================================
// MÉTODOS DE CONVERSIÓN
//=============================================================================

// ToModuleConfigResponse convierte modelo a DTO de respuesta
func ToModuleConfigResponse(config *models.OrganizationModuleConfig) *ModuleConfigResponse {
	response := &ModuleConfigResponse{
		ID:             config.ID,
		OrganizationID: config.OrganizationID,
		ModuleName:     config.ModuleName,
		IsEnabled:      config.IsEnabled,
		Version:        config.Version,
		InstalledBy:    config.InstalledBy,
		CreatedAt:      config.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      config.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// Parsear JSON fields
	if config.Resources != "" {
		json.Unmarshal([]byte(config.Resources), &response.Resources)
	}
	if config.DefaultRoles != "" {
		json.Unmarshal([]byte(config.DefaultRoles), &response.DefaultRoles)
	}
	if config.Settings != "" {
		json.Unmarshal([]byte(config.Settings), &response.Settings)
	}

	// Agregar información de relaciones si está disponible
	if config.Organization.ID != "" {
		response.Organization = &OrganizationInfo{
			ID:   config.Organization.ID,
			Name: config.Organization.Name,
			Slug: config.Organization.Slug,
			Type: config.Organization.Type,
		}
	}

	if config.Installer.ID != "" {
		response.Installer = &UserInfo{
			ID:        config.Installer.ID,
			Email:     config.Installer.Email,
			FirstName: config.Installer.FirstName,
			LastName:  config.Installer.LastName,
		}
	}

	return response
}

// ToModuleConfigResponseList convierte lista de modelos a DTOs
func ToModuleConfigResponseList(configs []models.OrganizationModuleConfig) []*ModuleConfigResponse {
	responses := make([]*ModuleConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = ToModuleConfigResponse(&config)
	}
	return responses
}

// ToCreateModuleConfiguration convierte DTO a configuración de módulo
func (r *CreateModuleConfigRequest) ToCreateModuleConfiguration() ModuleConfiguration {
	return ModuleConfiguration{
		ModuleName:   r.ModuleName,
		IsEnabled:    r.IsEnabled,
		Resources:    r.Resources,
		DefaultRoles: r.DefaultRoles,
		Settings:     r.Settings,
		Version:      r.Version,
	}
}

// GetAvailableModulesResponse retorna módulos disponibles para instalar
func GetAvailableModulesResponse() []AvailableModuleResponse {
	return []AvailableModuleResponse{
		{
			ModuleName:  "expenses",
			Description: "Sistema completo de gestión de gastos y presupuestos",
			Version:     "1.0.0",
			Resources:   models.ModuleResources["expenses"],
			Category:    "financial",
		},
		{
			ModuleName:  "inventory",
			Description: "Gestión avanzada de inventarios y almacenes",
			Version:     "1.0.0",
			Resources:   models.ModuleResources["inventory"],
			Category:    "warehouse",
		},
		{
			ModuleName:  "ecommerce",
			Description: "Plataforma completa de comercio electrónico",
			Version:     "1.0.0",
			Resources:   models.ModuleResources["ecommerce"],
			Category:    "sales",
		},
		{
			ModuleName:  "education",
			Description: "Sistema de gestión educativa y académica",
			Version:     "1.0.0",
			Resources:   models.ModuleResources["education"],
			Category:    "education",
		},
		{
			ModuleName:  "reports",
			Description: "Sistema avanzado de reportes y analytics",
			Version:     "1.0.0",
			Resources:   models.ModuleResources["reports"],
			Category:    "analytics",
		},
	}
}
