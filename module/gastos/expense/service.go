package expense

import (
	"fmt"
	"practicev2/module/gastos/model"
)

// =============================================================================
// INTERFACES Y ESTRUCTURAS DEL SERVICIO CON CONTEXTO
// =============================================================================

// ExpenseContext representa el contexto de operación para gastos
type ExpenseContext struct {
	UserID         uint
	OrganizationID uint
	DepartmentID   *uint
	Role           string
}

// expenseService define las operaciones de negocio para gastos con contexto organizacional
type expenseService interface {
	// CRUD básico con contexto
	CreateExpense(req CreateExpenseRequest, ctx ExpenseContext) (*ExpenseResponse, error)
	GetExpense(id uint) (*ExpenseResponse, error)
	GetAllExpenses() ([]*ExpenseResponse, error)
	UpdateExpense(id uint, req UpdateExpenseRequest) (*ExpenseResponse, error)
	DeleteExpense(id uint) error

	// Consultas con contexto
	GetUserExpenses(userID uint) ([]*ExpenseResponse, error)
	GetOrganizationExpenses(orgID uint) ([]*ExpenseResponse, error)
	GetDepartmentExpenses(deptID uint) ([]*ExpenseResponse, error)

	// Búsquedas
	SearchByTitle(title string) ([]*ExpenseResponse, error)
	GetByAmountRange(minAmount, maxAmount float64) ([]*ExpenseResponse, error)

	// Estadísticas
	GetExpenseStats() (*ExpenseStatsResponse, error)
	GetUserStats(userID uint) (*ExpenseStatsResponse, error)
}

// expenseServiceImpl implementa expenseService
type expenseServiceImpl struct {
	repository expenseRepository
}

// newExpenseService crea una nueva instancia del servicio de gastos
func newExpenseService(repository expenseRepository) expenseService {
	return &expenseServiceImpl{
		repository: repository,
	}
}

// =============================================================================
// IMPLEMENTACIÓN DE OPERACIONES CRUD CON CONTEXTO
// =============================================================================

// CreateExpense crea un nuevo gasto con contexto organizacional
func (s *expenseServiceImpl) CreateExpense(req CreateExpenseRequest, ctx ExpenseContext) (*ExpenseResponse, error) {
	// Validar datos básicos
	if req.Amount <= 0 {
		return nil, fmt.Errorf("el monto debe ser mayor a 0")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("el título es requerido")
	}

	// Convertir DTO a modelo con contexto organizacional
	expense := &model.Expense{
		Title:          req.Title,
		Description:    req.Description,
		Amount:         req.Amount,
		CreatedByID:    ctx.UserID,
		OrganizationID: ctx.OrganizationID,
		DepartmentID:   req.DepartmentID, // Puede ser nulo
		Status:         "pending",        // Estado inicial
	}

	// Si se proporciona un departamento en el contexto, usarlo como fallback
	if expense.DepartmentID == nil && ctx.DepartmentID != nil {
		expense.DepartmentID = ctx.DepartmentID
	}

	// Guardar en repositorio
	err := s.repository.Create(expense)
	if err != nil {
		return nil, fmt.Errorf("error al crear el gasto: %w", err)
	}

	// Convertir a DTO de respuesta
	return ToExpenseResponse(expense), nil
}

// GetExpense obtiene un gasto por ID
func (s *expenseServiceImpl) GetExpense(id uint) (*ExpenseResponse, error) {
	expense, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el gasto: %w", err)
	}

	return ToExpenseResponse(expense), nil
}

// GetAllExpenses obtiene todos los gastos
func (s *expenseServiceImpl) GetAllExpenses() ([]*ExpenseResponse, error) {
	expenses, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener gastos: %w", err)
	}

	// Convertir a DTOs de respuesta
	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// UpdateExpense actualiza un gasto existente
func (s *expenseServiceImpl) UpdateExpense(id uint, req UpdateExpenseRequest) (*ExpenseResponse, error) {
	// Obtener gasto existente
	existingExpense, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el gasto: %w", err)
	}

	// Actualizar campos si están presentes
	if req.Title != nil && *req.Title != "" {
		existingExpense.Title = *req.Title
	}
	if req.Description != nil {
		existingExpense.Description = *req.Description
	}
	if req.Amount != nil && *req.Amount > 0 {
		existingExpense.Amount = *req.Amount
	}
	if req.DepartmentID != nil {
		existingExpense.DepartmentID = req.DepartmentID
	}
	if req.Status != nil {
		existingExpense.Status = *req.Status
	}

	// Guardar cambios
	err = s.repository.Update(existingExpense)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar el gasto: %w", err)
	}

	return ToExpenseResponse(existingExpense), nil
}

// DeleteExpense elimina un gasto
func (s *expenseServiceImpl) DeleteExpense(id uint) error {
	// Verificar que el gasto existe
	_, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("error al obtener el gasto: %w", err)
	}

	// Eliminar
	err = s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("error al eliminar el gasto: %w", err)
	}

	return nil
}

// =============================================================================
// CONSULTAS CON CONTEXTO ORGANIZACIONAL
// =============================================================================

// GetUserExpenses obtiene todos los gastos de un usuario específico
func (s *expenseServiceImpl) GetUserExpenses(userID uint) ([]*ExpenseResponse, error) {
	expenses, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener gastos del usuario: %w", err)
	}

	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// GetOrganizationExpenses obtiene todos los gastos de una organización
func (s *expenseServiceImpl) GetOrganizationExpenses(orgID uint) ([]*ExpenseResponse, error) {
	expenses, err := s.repository.GetByOrganizationID(orgID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener gastos de la organización: %w", err)
	}

	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// GetDepartmentExpenses obtiene todos los gastos de un departamento
func (s *expenseServiceImpl) GetDepartmentExpenses(deptID uint) ([]*ExpenseResponse, error) {
	expenses, err := s.repository.GetByDepartmentID(deptID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener gastos del departamento: %w", err)
	}

	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// =============================================================================
// BÚSQUEDAS Y CONSULTAS ESPECÍFICAS
// =============================================================================

// SearchByTitle busca gastos por título
func (s *expenseServiceImpl) SearchByTitle(title string) ([]*ExpenseResponse, error) {
	if title == "" {
		return nil, fmt.Errorf("el título de búsqueda es requerido")
	}

	expenses, err := s.repository.GetByTitle(title)
	if err != nil {
		return nil, fmt.Errorf("error al buscar gastos por título: %w", err)
	}

	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// GetByAmountRange obtiene gastos dentro de un rango de montos
func (s *expenseServiceImpl) GetByAmountRange(minAmount, maxAmount float64) ([]*ExpenseResponse, error) {
	if minAmount < 0 || maxAmount < 0 {
		return nil, fmt.Errorf("los montos no pueden ser negativos")
	}
	if minAmount > maxAmount {
		return nil, fmt.Errorf("el monto mínimo no puede ser mayor al máximo")
	}

	expenses, err := s.repository.GetByAmountRange(minAmount, maxAmount)
	if err != nil {
		return nil, fmt.Errorf("error al obtener gastos por rango de monto: %w", err)
	}

	expenseResponses := make([]*ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ToExpenseResponse(&expense)
	}

	return expenseResponses, nil
}

// =============================================================================
// ESTADÍSTICAS
// =============================================================================

// GetExpenseStats obtiene estadísticas generales de gastos
func (s *expenseServiceImpl) GetExpenseStats() (*ExpenseStatsResponse, error) {
	expenses, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas de gastos: %w", err)
	}

	return CalculateStats(expenses), nil
}

// GetUserStats obtiene estadísticas de gastos de un usuario específico
func (s *expenseServiceImpl) GetUserStats(userID uint) (*ExpenseStatsResponse, error) {
	expenses, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas del usuario: %w", err)
	}

	return CalculateStats(expenses), nil
}
