package rutascompuestas

import (
	"gorm.io/gorm"

	"practicev2/database"
	"practicev2/module/auth"
	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
)

// ----------------------------------------------------------------------------
// Helper
// ----------------------------------------------------------------------------

// getAuthenticatedUser extrae y valida el usuario autenticado del contexto.
func getAuthenticatedUser(c *fiber.Ctx) (uint64, error) {
	uid, ok := c.Locals("usuario_id").(uint64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Usuario no autenticado")
	}
	return uid, nil
}

// ----------------------------------------------------------------------------
// Repositorio
// ----------------------------------------------------------------------------

func updateActiveSessionCompany(usuarioID, selectedEmpresaID uint64) error {
	tx := database.DBconn.Begin()

	// Verificar que el usuario tiene la empresa autorizada
	var count int64
	if err := tx.Model(&auth.UsuarioEmpresa{}).
		Where("usuario_id = ? AND empresa_id = ?", usuarioID, selectedEmpresaID).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		tx.Rollback()
		return fiber.NewError(fiber.StatusUnauthorized, "Empresa no autorizada")
	}

	// Actualizar: marcar la empresa seleccionada como activa
	if err := tx.Model(&auth.UsuarioEmpresa{}).
		Where("usuario_id = ?", usuarioID).
		Update("active_sesion", gorm.Expr("CASE WHEN empresa_id = ? THEN ? ELSE ? END", selectedEmpresaID, true, false)).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func getFirstActiveCompanyService(usuarioID uint64) (uint64, error) {
	var usuarioEmpresa auth.UsuarioEmpresa
	err := database.DBconn.
		Where("usuario_id = ? AND active_sesion = ?", usuarioID, true).
		Order("created_at").
		First(&usuarioEmpresa).Error
	if err != nil {
		return 0, err
	}
	return usuarioEmpresa.EmpresaID, nil
}

// ----------------------------------------------------------------------------
// Service
// ----------------------------------------------------------------------------

func updateActiveCompanyService(selectedEmpresaID, usuarioID uint64) error {
	return updateActiveSessionCompany(usuarioID, selectedEmpresaID)
}

// ----------------------------------------------------------------------------
// Handler
// ----------------------------------------------------------------------------

func updateActiveCompanyHandler(c *fiber.Ctx) error {
	type Request struct {
		EmpresaID uint64 `json:"empresa_id"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Solicitud inválida"})
	}
	uid, err := getAuthenticatedUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	if err := updateActiveCompanyService(req.EmpresaID, uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Empresa activa actualizada"})
}

func getFirstActiveCompanyHandler(c *fiber.Ctx) error {
	uid, err := getAuthenticatedUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	empresaID, err := getFirstActiveCompanyService(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"empresa_id": empresaID})
}

// ----------------------------------------------------------------------------
// Registro de rutas
// ----------------------------------------------------------------------------

func RegisterActiveCompanyRoute(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/empresa/active", updateActiveCompanyHandler)
	api.Get("/empresa/active-company", getFirstActiveCompanyHandler)
}

func init() {
	registry.RegisterModule(RegisterActiveCompanyRoute)
}
