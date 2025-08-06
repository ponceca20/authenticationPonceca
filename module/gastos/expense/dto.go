package expense

import (
	"practicev2/module/gastos/model"
	"time"
)

// =============================================================================
// DTOs PARA EL MÓDULO DE GASTOS CON CONTEXTO ORGANIZACIONAL
// =============================================================================

// CreateExpenseRequest representa los datos necesarios para crear un nuevo gasto
type CreateExpenseRequest struct {
	Title        string  `json:"title" validate:"required,min=3,max=255"`
	Description  string  `json:"description" validate:"max=1000"`
	Amount       float64 `json:"amount" validate:"required,min=0.01"`
	DepartmentID *uint   `json:"department_id,omitempty"`
}

// UpdateExpenseRequest representa los datos que se pueden actualizar de un gasto
type UpdateExpenseRequest struct {
	Title        *string  `json:"title" validate:"omitempty,min=3,max=255"`
	Description  *string  `json:"description" validate:"omitempty,max=1000"`
	Amount       *float64 `json:"amount" validate:"omitempty,min=0.01"`
	DepartmentID *uint    `json:"department_id,omitempty"`
	Status       *string  `json:"status" validate:"omitempty,oneof=pending approved rejected"`
}

// ExpenseResponse representa un gasto en las respuestas de la API
type ExpenseResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Amount         float64   `json:"amount"`
	CreatedByID    uint      `json:"created_by_id"`
	OrganizationID uint      `json:"organization_id"`
	DepartmentID   *uint     `json:"department_id,omitempty"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ExpenseListResponse representa una lista de gastos con paginación (para futuras mejoras)
type ExpenseListResponse struct {
	Expenses []ExpenseResponse `json:"expenses"`
	Total    int               `json:"total"`
	Page     int               `json:"page,omitempty"`
	Limit    int               `json:"limit,omitempty"`
}

// ExpenseStatsResponse representa estadísticas de gastos
type ExpenseStatsResponse struct {
	TotalExpenses int            `json:"total_expenses"`
	TotalAmount   float64        `json:"total_amount"`
	AvgAmount     float64        `json:"avg_amount"`
	Status        map[string]int `json:"status_breakdown"`
}

// =============================================================================
// FUNCIONES DE CONVERSIÓN
// =============================================================================

// ToExpenseResponse convierte un modelo Expense a ExpenseResponse
func ToExpenseResponse(expense *model.Expense) *ExpenseResponse {
	return &ExpenseResponse{
		ID:             expense.ID,
		Title:          expense.Title,
		Description:    expense.Description,
		Amount:         expense.Amount,
		CreatedByID:    expense.CreatedByID,
		OrganizationID: expense.OrganizationID,
		DepartmentID:   expense.DepartmentID,
		Status:         expense.Status,
		CreatedAt:      expense.CreatedAt,
		UpdatedAt:      expense.UpdatedAt,
	}
}

// ToExpenseResponseList convierte una lista de modelos Expense a ExpenseListResponse
func ToExpenseResponseList(expenses []model.Expense) *ExpenseListResponse {
	expenseResponses := make([]ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = *ToExpenseResponse(&expense)
	}

	return &ExpenseListResponse{
		Expenses: expenseResponses,
		Total:    len(expenses),
	}
}

// CalculateStats calcula estadísticas básicas de una lista de gastos
func CalculateStats(expenses []model.Expense) *ExpenseStatsResponse {
	stats := &ExpenseStatsResponse{
		TotalExpenses: len(expenses),
		Status:        make(map[string]int),
	}

	if len(expenses) == 0 {
		return stats
	}

	var totalAmount float64
	for _, expense := range expenses {
		totalAmount += expense.Amount
		stats.Status[expense.Status]++
	}

	stats.TotalAmount = totalAmount
	stats.AvgAmount = totalAmount / float64(len(expenses))

	return stats
}
