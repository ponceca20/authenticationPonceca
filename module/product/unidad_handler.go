package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UnidadHandler maneja las solicitudes HTTP relacionadas con Unidad
type UnidadHandler struct {
	Service   *UnidadService
	validator *validator.Validate
}

// NewUnidadHandler crea una nueva instancia del handler
func NewUnidadHandler(s *UnidadService) *UnidadHandler {
	return &UnidadHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllUnidades obtiene todas las unidades para la empresa del usuario autenticado
func (h *UnidadHandler) GetAllUnidades(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	unidades, err := h.Service.GetAllUnidades(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar unidades")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      unidades,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetUnidad obtiene una unidad específica
func (h *UnidadHandler) GetUnidad(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	unidad, err := h.Service.GetUnidad(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Unidad no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      unidad,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateUnidad crea una nueva unidad
func (h *UnidadHandler) CreateUnidad(c *fiber.Ctx) error {
	var unidad Unidad
	if err := c.BodyParser(&unidad); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&unidad); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	unidad.EmpresaID = empresaID

	if err := h.Service.CreateUnidad(&unidad); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la unidad")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      unidad,
		"message":   "Unidad creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateUnidad actualiza una unidad existente
func (h *UnidadHandler) UpdateUnidad(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var unidad Unidad
	if err := c.BodyParser(&unidad); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&unidad); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	unidad.EmpresaID = empresaID

	updated, err := h.Service.UpdateUnidad(id, empresaID, &unidad)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la unidad")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Unidad actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteUnidad elimina una unidad
func (h *UnidadHandler) DeleteUnidad(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteUnidad(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la unidad")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Unidad eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedUnidades inicializa múltiples unidades (para administradores)
func (h *UnidadHandler) SeedUnidades(c *fiber.Ctx) error {
	var unidades []Unidad
	if err := c.BodyParser(&unidades); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedUnidades(unidades); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las unidades")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Unidades inicializadas exitosamente",
		"count":     len(unidades),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchUnidad actualiza parcialmente una unidad
func (h *UnidadHandler) PatchUnidad(c *fiber.Ctx) error {
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

	unidad, err := h.Service.PatchUnidad(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la unidad")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      unidad,
		"message":   "Unidad actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesUnidad registra las rutas para la API de Unidad
func RegisterRoutesUnidad(app *fiber.App) {
	repo := NewUnidadRepository()
	service := NewUnidadService(repo)
	handler := NewUnidadHandler(service)

	api := app.Group("/api/v1/unidades")

	api.Get("/", handler.GetAllUnidades)
	api.Get("/:id", handler.GetUnidad)
	api.Post("/", handler.CreateUnidad)
	api.Put("/:id", handler.UpdateUnidad)
	api.Patch("/:id", handler.PatchUnidad)
	api.Delete("/:id", handler.DeleteUnidad)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedUnidades)
}

func init() {
	registry.RegisterModule(RegisterRoutesUnidad)
}
