package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// FamiliaHandler maneja las solicitudes HTTP relacionadas con Familia
type FamiliaHandler struct {
	Service   *FamiliaService
	validator *validator.Validate
}

// NewFamiliaHandler crea una nueva instancia del handler
func NewFamiliaHandler(s *FamiliaService) *FamiliaHandler {
	return &FamiliaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllFamilias obtiene todas las familias para la empresa del usuario autenticado
func (h *FamiliaHandler) GetAllFamilias(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	familias, err := h.Service.GetAllFamilias(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar familias")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      familias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetFamiliasByCartaID obtiene todas las familias para una carta específica
func (h *FamiliaHandler) GetFamiliasByCartaID(c *fiber.Ctx) error {
	cartaID, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de carta inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	familias, err := h.Service.GetFamiliasByCartaID(cartaID, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar familias para la carta especificada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      familias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetFamilia obtiene una familia específica
func (h *FamiliaHandler) GetFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	familia, err := h.Service.GetFamilia(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Familia no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      familia,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateFamilia crea una nueva familia
func (h *FamiliaHandler) CreateFamilia(c *fiber.Ctx) error {
	var familia Familia
	if err := c.BodyParser(&familia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&familia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	familia.EmpresaID = empresaID

	if err := h.Service.CreateFamilia(&familia); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la familia")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      familia,
		"message":   "Familia creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateFamilia actualiza una familia existente
func (h *FamiliaHandler) UpdateFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var familia Familia
	if err := c.BodyParser(&familia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&familia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	familia.EmpresaID = empresaID

	updated, err := h.Service.UpdateFamilia(id, empresaID, &familia)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Familia actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteFamilia elimina una familia
func (h *FamiliaHandler) DeleteFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteFamilia(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Familia eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedFamilias inicializa múltiples familias (para administradores)
func (h *FamiliaHandler) SeedFamilias(c *fiber.Ctx) error {
	var familias []Familia
	if err := c.BodyParser(&familias); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedFamilias(familias); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las familias")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Familias inicializadas exitosamente",
		"count":     len(familias),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchFamilia actualiza parcialmente una familia
func (h *FamiliaHandler) PatchFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Proteger el campo empresaID contra modificaciones
	delete(fields, "empresa_id")

	familia, err := h.Service.PatchFamilia(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      familia,
		"message":   "Familia actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesFamilia registra las rutas para la API de Familia
func RegisterRoutesFamilia(app *fiber.App) {
	repo := NewFamiliaRepository()
	service := NewFamiliaService(repo)
	handler := NewFamiliaHandler(service)

	api := app.Group("/api/v1/familias")

	api.Get("/", handler.GetAllFamilias)
	api.Get("/:id", handler.GetFamilia)
	api.Get("/carta/:id", handler.GetFamiliasByCartaID)
	api.Post("/", handler.CreateFamilia)
	api.Put("/:id", handler.UpdateFamilia)
	api.Patch("/:id", handler.PatchFamilia)
	api.Delete("/:id", handler.DeleteFamilia)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedFamilias)
}

func init() {
	registry.RegisterModule(RegisterRoutesFamilia)
}
