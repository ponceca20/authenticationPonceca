package department

import "practicev2/module/authentication/models"

// DepartmentDTO is used for creating or updating a department.
type DepartmentDTO struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

// DepartmentResponseDTO is a public representation of a department.
type DepartmentResponseDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
}

// ToDepartmentResponseDTO converts a Department model to a public DTO.
func ToDepartmentResponseDTO(dept *models.Department) DepartmentResponseDTO {
	return DepartmentResponseDTO{
		ID:          dept.ID,
		Name:        dept.Name,
		Description: dept.Description,
		ParentID:    dept.ParentID,
	}
}
