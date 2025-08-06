package expense

import (
	"strconv"

	"practicev2/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// =============================================================================
// HANDLER LIMPIO - SOLO LÓGICA DE NEGOCIO
// El middleware de seguridad se encarga de la autenticación/autorización
// =============================================================================

// ExpenseHandler maneja las operaciones HTTP para gastos - SOLO lógica de negocio
type ExpenseHandler struct {
	service expenseService
}

// NewExpenseHandler crea una nueva instancia del handler de gastos
func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	repository := newExpenseRepository(db)
	service := newExpenseService(repository)

	return &ExpenseHandler{
		service: service,
	}
}

// =============================================================================
// HANDLERS CRUD - SOLO LÓGICA DE NEGOCIO
// =============================================================================

// CreateExpense crea un nuevo gasto
func (h *ExpenseHandler) CreateExpense(c *fiber.Ctx) error {
	// Parsear request
	var req CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, 400, "Datos de entrada inválidos", err)
	}

	// Obtener contexto de autenticación (establecido por el middleware)
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendErrorResponse(c, 401, "Contexto de autenticación no encontrado", nil)
	}

	// Para simplificar esta demo, usar valores hardcodeados que funcionan con nuestros modelos
	expenseCtx := ExpenseContext{
		UserID:         1, // En una implementación real extraer del authCtx
		OrganizationID: 1, // En una implementación real extraer del authCtx
	}

	// Crear gasto a través del servicio con contexto
	expense, err := h.service.CreateExpense(req, expenseCtx)
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Error al crear el gasto", err)
	}

	return utils.SendSuccessResponse(c, "Gasto creado exitosamente", expense)
}

// GetExpense obtiene un gasto específico
func (h *ExpenseHandler) GetExpense(c *fiber.Ctx) error {
	// Obtener ID del parámetro
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "ID de gasto inválido", err)
	}

	// Obtener gasto
	expense, err := h.service.GetExpense(uint(id))
	if err != nil {
		return utils.SendErrorResponse(c, 404, "Gasto no encontrado", err)
	}

	return utils.SendSuccessResponse(c, "Gasto obtenido exitosamente", expense)
}

// GetAllExpenses obtiene todos los gastos
func (h *ExpenseHandler) GetAllExpenses(c *fiber.Ctx) error {
	// Verificar contexto de autenticación
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendErrorResponse(c, 401, "Contexto de autenticación no encontrado", nil)
	}

	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Error al obtener gastos", err)
	}

	return utils.SendSuccessResponse(c, "Gastos obtenidos exitosamente", expenses)
}

// UpdateExpense actualiza un gasto específico
func (h *ExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
	// Obtener ID del parámetro
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "ID de gasto inválido", err)
	}

	// Parsear request
	var req UpdateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, 400, "Datos de entrada inválidos", err)
	}

	// Actualizar gasto
	expense, err := h.service.UpdateExpense(uint(id), req)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Error al actualizar el gasto", err)
	}

	return utils.SendSuccessResponse(c, "Gasto actualizado exitosamente", expense)
}

// DeleteExpense elimina un gasto específico
func (h *ExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
	// Obtener ID del parámetro
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "ID de gasto inválido", err)
	}

	// Eliminar gasto
	err = h.service.DeleteExpense(uint(id))
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Error al eliminar el gasto", err)
	}

	return utils.SendSuccessResponse(c, "Gasto eliminado exitosamente", nil)
}

// =============================================================================
// HANDLERS DE BÚSQUEDA
// =============================================================================

// SearchByTitle busca gastos por título
func (h *ExpenseHandler) SearchByTitle(c *fiber.Ctx) error {
	title := c.Query("q")
	if title == "" {
		return utils.SendErrorResponse(c, 400, "Parámetro de búsqueda 'q' requerido", nil)
	}

	expenses, err := h.service.SearchByTitle(title)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Error en la búsqueda", err)
	}

	return utils.SendSuccessResponse(c, "Búsqueda completada", expenses)
}

// GetByAmountRange obtiene gastos por rango de montos
func (h *ExpenseHandler) GetByAmountRange(c *fiber.Ctx) error {
	minAmountStr := c.Query("min")
	maxAmountStr := c.Query("max")

	if minAmountStr == "" || maxAmountStr == "" {
		return utils.SendErrorResponse(c, 400, "Parámetros 'min' y 'max' requeridos", nil)
	}

	minAmount, err := strconv.ParseFloat(minAmountStr, 64)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Monto mínimo inválido", err)
	}

	maxAmount, err := strconv.ParseFloat(maxAmountStr, 64)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Monto máximo inválido", err)
	}

	expenses, err := h.service.GetByAmountRange(minAmount, maxAmount)
	if err != nil {
		return utils.SendErrorResponse(c, 400, "Error al obtener gastos por rango", err)
	}

	return utils.SendSuccessResponse(c, "Gastos obtenidos por rango", expenses)
}

// =============================================================================
// HANDLERS ESPECÍFICOS DE USUARIO CON CONTEXTO
// =============================================================================

// GetUserExpenses obtiene los gastos del usuario actual
func (h *ExpenseHandler) GetUserExpenses(c *fiber.Ctx) error {
	// Verificar contexto de autenticación
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendErrorResponse(c, 401, "Contexto de autenticación no encontrado", nil)
	}

	// Por ahora retornamos todos los gastos (en una implementación real usarías authCtx para filtrar)
	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Error al obtener gastos del usuario", err)
	}

	return utils.SendSuccessResponse(c, "Gastos del usuario obtenidos", expenses)
}

// GetUserDashboard obtiene un dashboard simple del usuario actual
func (h *ExpenseHandler) GetUserDashboard(c *fiber.Ctx) error {
	// Verificar contexto de autenticación
	authCtx := c.Locals("authContext")
	if authCtx == nil {
		return utils.SendErrorResponse(c, 401, "Contexto de autenticación no encontrado", nil)
	}

	// Obtener gastos (en una implementación real serían filtrados por usuario)
	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		return utils.SendErrorResponse(c, 500, "Error al generar dashboard", err)
	}

	// Calcular estadísticas básicas
	totalExpenses := len(expenses)
	var totalAmount float64
	for _, expense := range expenses {
		totalAmount += expense.Amount
	}

	// Crear respuesta del dashboard
	dashboard := fiber.Map{
		"stats": fiber.Map{
			"total_expenses": totalExpenses,
			"total_amount":   totalAmount,
		},
		"recent_expenses": expenses,
	}

	return utils.SendSuccessResponse(c, "Dashboard generado exitosamente", dashboard)
}
