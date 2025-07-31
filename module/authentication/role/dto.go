package role

import "practicev2/module/authentication/models"

// PermissionDTO is used within the RoleDTO to define permissions.
type PermissionDTO struct {
	Resource string   `json:"resource" validate:"required"`
	Actions  []string `json:"actions" validate:"required,min=1"`
	Scope    string   `json:"scope" validate:"required,oneof=own department organization all"`
}

// RoleDTO is used for creating or updating a role.
type RoleDTO struct {
	Name           string          `json:"name" validate:"required"`
	DisplayName    string          `json:"display_name" validate:"required"`
	Description    string          `json:"description,omitempty"`
	HierarchyLevel int             `json:"hierarchy_level"`
	Permissions    []PermissionDTO `json:"permissions" validate:"required,dive"`
}

// RoleResponseDTO is a public representation of a role.
type RoleResponseDTO struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	DisplayName    string          `json:"display_name"`
	Description    string          `json:"description,omitempty"`
	HierarchyLevel int             `json:"hierarchy_level"`
	IsSystemRole   bool            `json:"is_system_role"`
	Permissions    []PermissionDTO `json:"permissions"`
}

// ToRoleResponseDTO converts a Role model to a public DTO.
func ToRoleResponseDTO(role *models.Role) RoleResponseDTO {
	dto := RoleResponseDTO{
		ID:             role.ID,
		Name:           role.Name,
		DisplayName:    role.DisplayName,
		Description:    role.Description,
		HierarchyLevel: role.HierarchyLevel,
		IsSystemRole:   role.IsSystemRole,
		Permissions:    []PermissionDTO{},
	}

	// Group permissions by resource and scope for the DTO
	permMap := make(map[string]*PermissionDTO)
	for _, p := range role.Permissions {
		key := p.Resource + ":" + p.Scope
		if existing, ok := permMap[key]; ok {
			existing.Actions = append(existing.Actions, p.Action)
		} else {
			permMap[key] = &PermissionDTO{
				Resource: p.Resource,
				Actions:  []string{p.Action},
				Scope:    p.Scope,
			}
		}
	}

	for _, p := range permMap {
		dto.Permissions = append(dto.Permissions, *p)
	}

	return dto
}
