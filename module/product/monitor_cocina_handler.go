package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// MonitorCocinaHandler maneja las solicitudes HTTP relacionadas con MonitorCocina
type MonitorCocinaHandler struct {
	Service   *MonitorCocinaService
	validator *validator.Validate
}

// NewMonitorCocinaHandler crea una nueva instancia del handler
func NewMonitorCocinaHandler(s *MonitorCocinaService) *MonitorCocinaHandler {
	return &MonitorCocinaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllMonitoresCocina obtiene todos los monitores de cocina para la empresa del usuario autenticado
func (h *MonitorCocinaHandler) GetAllMonitoresCocina(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	monitores, err := h.Service.GetAllMonitoresCocina(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar monitores de cocina")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitores,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetMonitorCocina obtiene un monitor de cocina específico
func (h *MonitorCocinaHandler) GetMonitorCocina(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	monitor, err := h.Service.GetMonitorCocina(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Monitor de cocina no encontrado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateMonitorCocina crea un nuevo monitor de cocina
func (h *MonitorCocinaHandler) CreateMonitorCocina(c *fiber.Ctx) error {
	var monitor MonitorCocina
	if err := c.BodyParser(&monitor); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&monitor); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	monitor.EmpresaID = empresaID

	if err := h.Service.CreateMonitorCocina(&monitor); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el monitor de cocina")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"message":   "Monitor de cocina creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateMonitorCocina actualiza un monitor de cocina existente
func (h *MonitorCocinaHandler) UpdateMonitorCocina(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var monitor MonitorCocina
	if err := c.BodyParser(&monitor); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&monitor); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	monitor.EmpresaID = empresaID

	updated, err := h.Service.UpdateMonitorCocina(id, empresaID, &monitor)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el monitor de cocina")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Monitor de cocina actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteMonitorCocina elimina un monitor de cocina
func (h *MonitorCocinaHandler) DeleteMonitorCocina(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteMonitorCocina(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el monitor de cocina")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Monitor de cocina eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedMonitoresCocina inicializa múltiples monitores de cocina (para administradores)
func (h *MonitorCocinaHandler) SeedMonitoresCocina(c *fiber.Ctx) error {
	var monitores []MonitorCocina
	if err := c.BodyParser(&monitores); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedMonitoresCocina(monitores); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los monitores de cocina")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Monitores de cocina inicializados exitosamente",
		"count":     len(monitores),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchMonitorCocina actualiza parcialmente un monitor de cocina
func (h *MonitorCocinaHandler) PatchMonitorCocina(c *fiber.Ctx) error {
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

	monitor, err := h.Service.PatchMonitorCocina(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el monitor de cocina")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"message":   "Monitor de cocina actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesMonitorCocina registra las rutas para la API de MonitorCocina
func RegisterRoutesMonitorCocina(app *fiber.App) {
	repo := NewMonitorCocinaRepository()
	service := NewMonitorCocinaService(repo)
	handler := NewMonitorCocinaHandler(service)

	api := app.Group("/api/v1/monitores-cocina")

	api.Get("/", handler.GetAllMonitoresCocina)
	api.Get("/:id", handler.GetMonitorCocina)
	api.Post("/", handler.CreateMonitorCocina)
	api.Put("/:id", handler.UpdateMonitorCocina)
	api.Patch("/:id", handler.PatchMonitorCocina)
	api.Delete("/:id", handler.DeleteMonitorCocina)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedMonitoresCocina)
}

func init() {
	registry.RegisterModule(RegisterRoutesMonitorCocina)
}
