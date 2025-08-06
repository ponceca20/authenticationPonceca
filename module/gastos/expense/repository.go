package expense

import (
	"errors"
	"fmt"
	"practicev2/module/gastos/model"

	"gorm.io/gorm"
)

// =============================================================================
// INTERFACE DEL REPOSITORIO CON CONTEXTO ORGANIZACIONAL
// =============================================================================

// expenseRepository define la interfaz para el acceso a datos de gastos con contexto
type expenseRepository interface {
	// CRUD básico
	Create(expense *model.Expense) error
	GetByID(id uint) (*model.Expense, error)
	GetAll() ([]model.Expense, error)
	Update(expense *model.Expense) error
	Delete(id uint) error

	// Consultas con contexto organizacional
	GetByUserID(userID uint) ([]model.Expense, error)
	GetByOrganizationID(orgID uint) ([]model.Expense, error)
	GetByDepartmentID(deptID uint) ([]model.Expense, error)

	// Consultas específicas
	GetByTitle(title string) ([]model.Expense, error)
	GetByAmountRange(minAmount, maxAmount float64) ([]model.Expense, error)
	GetByStatus(status string) ([]model.Expense, error)

	// Consultas combinadas con contexto
	GetByUserAndOrganization(userID, orgID uint) ([]model.Expense, error)
	GetByTitleAndOrganization(title string, orgID uint) ([]model.Expense, error)
	GetByAmountRangeAndOrganization(minAmount, maxAmount float64, orgID uint) ([]model.Expense, error)
}

// =============================================================================
// IMPLEMENTACIÓN DEL REPOSITORIO
// =============================================================================

// expenseRepositoryImpl implementa expenseRepository usando GORM
type expenseRepositoryImpl struct {
	db *gorm.DB
}

// newExpenseRepository crea una nueva instancia del repositorio de gastos
func newExpenseRepository(db *gorm.DB) expenseRepository {
	return &expenseRepositoryImpl{db: db}
}

// =============================================================================
// MÉTODOS CRUD BÁSICOS
// =============================================================================

// Create crea un nuevo gasto en la base de datos
func (r *expenseRepositoryImpl) Create(expense *model.Expense) error {
	return r.db.Create(expense).Error
}

// GetByID obtiene un gasto por su ID
func (r *expenseRepositoryImpl) GetByID(id uint) (*model.Expense, error) {
	var expense model.Expense
	err := r.db.Where("id = ?", id).First(&expense).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("expense with ID %d not found", id)
		}
		return nil, err
	}

	return &expense, nil
}

// GetAll obtiene todos los gastos
func (r *expenseRepositoryImpl) GetAll() ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// Update actualiza un gasto existente
func (r *expenseRepositoryImpl) Update(expense *model.Expense) error {
	return r.db.Save(expense).Error
}

// Delete elimina un gasto por su ID
func (r *expenseRepositoryImpl) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Expense{}).Error
}

// =============================================================================
// CONSULTAS CON CONTEXTO ORGANIZACIONAL
// =============================================================================

// GetByUserID obtiene todos los gastos de un usuario específico
func (r *expenseRepositoryImpl) GetByUserID(userID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("created_by_id = ?", userID).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByOrganizationID obtiene todos los gastos de una organización específica
func (r *expenseRepositoryImpl) GetByOrganizationID(orgID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("organization_id = ?", orgID).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByDepartmentID obtiene todos los gastos de un departamento específico
func (r *expenseRepositoryImpl) GetByDepartmentID(deptID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("department_id = ?", deptID).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// =============================================================================
// CONSULTAS ESPECÍFICAS
// =============================================================================

// GetByTitle busca gastos por título (búsqueda parcial)
func (r *expenseRepositoryImpl) GetByTitle(title string) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("title LIKE ?", "%"+title+"%").Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByAmountRange obtiene gastos dentro de un rango de montos
func (r *expenseRepositoryImpl) GetByAmountRange(minAmount, maxAmount float64) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("amount BETWEEN ? AND ?", minAmount, maxAmount).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByStatus obtiene gastos por estado
func (r *expenseRepositoryImpl) GetByStatus(status string) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// =============================================================================
// CONSULTAS COMBINADAS CON CONTEXTO
// =============================================================================

// GetByUserAndOrganization obtiene gastos de un usuario en una organización específica
func (r *expenseRepositoryImpl) GetByUserAndOrganization(userID, orgID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("created_by_id = ? AND organization_id = ?", userID, orgID).
		Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByTitleAndOrganization busca gastos por título dentro de una organización
func (r *expenseRepositoryImpl) GetByTitleAndOrganization(title string, orgID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("title LIKE ? AND organization_id = ?", "%"+title+"%", orgID).
		Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}

// GetByAmountRangeAndOrganization obtiene gastos por rango de montos dentro de una organización
func (r *expenseRepositoryImpl) GetByAmountRangeAndOrganization(minAmount, maxAmount float64, orgID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Where("amount BETWEEN ? AND ? AND organization_id = ?", minAmount, maxAmount, orgID).
		Order("created_at DESC").Find(&expenses).Error
	return expenses, err
}
