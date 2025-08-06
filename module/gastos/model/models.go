package model

import (
	"time"
)

// =============================================================================
// MODELO DE GASTOS CON CONTEXTO ORGANIZACIONAL
// =============================================================================

// Expense representa un gasto en el sistema con información de contexto completa
type Expense struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string  `json:"title" gorm:"not null;size:255"`
	Description string  `json:"description" gorm:"size:1000"`
	Amount      float64 `json:"amount" gorm:"not null;type:decimal(10,2)"`

	// Información de contexto organizacional
	CreatedByID    uint  `json:"created_by_id" gorm:"not null;index"`
	OrganizationID uint  `json:"organization_id" gorm:"not null;index"`
	DepartmentID   *uint `json:"department_id" gorm:"index"`

	// Campos de auditoría
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Estado del gasto (para futuras implementaciones de workflow)
	Status string `json:"status" gorm:"default:'pending';size:50"`
}

// TableName ensures the table uses singular naming convention
func (Expense) TableName() string {
	return "expense"
}

// =============================================================================
// MÉTODOS DEL MODELO
// =============================================================================

// IsOwnedBy verifica si el gasto pertenece al usuario especificado
func (e *Expense) IsOwnedBy(userID uint) bool {
	return e.CreatedByID == userID
}

// BelongsToOrganization verifica si el gasto pertenece a la organización especificada
func (e *Expense) BelongsToOrganization(orgID uint) bool {
	return e.OrganizationID == orgID
}

// BelongsToDepartment verifica si el gasto pertenece al departamento especificado
func (e *Expense) BelongsToDepartment(deptID uint) bool {
	return e.DepartmentID != nil && *e.DepartmentID == deptID
}
