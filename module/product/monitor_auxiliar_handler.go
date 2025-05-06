package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// MonitorAuxiliarHandler maneja las solicitudes HTTP relacionadas con MonitorAuxiliar
type MonitorAuxiliarHandler struct {
	Service   *MonitorAuxiliarService
	validator *validator.Validate
}

// NewMonitorAuxiliarHandler crea una nueva instancia del handler
func NewMonitorAuxiliarHandler(s *MonitorAuxiliarService) *MonitorAuxiliarHandler {
	return &MonitorAuxiliarHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllMonitoresAuxiliar obtiene todos los monitores auxiliares para la empresa del usuario autenticado
func (h *MonitorAuxiliarHandler) GetAllMonitoresAuxiliar(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	monitores, err := h.Service.GetAllMonitoresAuxiliar(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar monitores auxiliares")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitores,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetMonitorAuxiliar obtiene un monitor auxiliar específico
func (h *MonitorAuxiliarHandler) GetMonitorAuxiliar(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	monitor, err := h.Service.GetMonitorAuxiliar(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Monitor auxiliar no encontrado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateMonitorAuxiliar crea un nuevo monitor auxiliar
func (h *MonitorAuxiliarHandler) CreateMonitorAuxiliar(c *fiber.Ctx) error {
	var monitor MonitorAuxiliar
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

	if err := h.Service.CreateMonitorAuxiliar(&monitor); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el monitor auxiliar")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"message":   "Monitor auxiliar creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateMonitorAuxiliar actualiza un monitor auxiliar existente
func (h *MonitorAuxiliarHandler) UpdateMonitorAuxiliar(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var monitor MonitorAuxiliar
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

	updated, err := h.Service.UpdateMonitorAuxiliar(id, empresaID, &monitor)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el monitor auxiliar")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Monitor auxiliar actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteMonitorAuxiliar elimina un monitor auxiliar
func (h *MonitorAuxiliarHandler) DeleteMonitorAuxiliar(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteMonitorAuxiliar(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el monitor auxiliar")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Monitor auxiliar eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedMonitoresAuxiliar inicializa múltiples monitores auxiliares (para administradores)
func (h *MonitorAuxiliarHandler) SeedMonitoresAuxiliar(c *fiber.Ctx) error {
	var monitores []MonitorAuxiliar
	if err := c.BodyParser(&monitores); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedMonitoresAuxiliar(monitores); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los monitores auxiliares")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Monitores auxiliares inicializados exitosamente",
		"count":     len(monitores),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchMonitorAuxiliar actualiza parcialmente un monitor auxiliar
func (h *MonitorAuxiliarHandler) PatchMonitorAuxiliar(c *fiber.Ctx) error {
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

	monitor, err := h.Service.PatchMonitorAuxiliar(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el monitor auxiliar")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      monitor,
		"message":   "Monitor auxiliar actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesMonitorAuxiliar registra las rutas para la API de MonitorAuxiliar
func RegisterRoutesMonitorAuxiliar(app *fiber.App) {
	repo := NewMonitorAuxiliarRepository()
	service := NewMonitorAuxiliarService(repo)
	handler := NewMonitorAuxiliarHandler(service)

	api := app.Group("/api/v1/monitores-auxiliar")

	api.Get("/", handler.GetAllMonitoresAuxiliar)
	api.Get("/:id", handler.GetMonitorAuxiliar)
	api.Post("/", handler.CreateMonitorAuxiliar)
	api.Put("/:id", handler.UpdateMonitorAuxiliar)
	api.Patch("/:id", handler.PatchMonitorAuxiliar)
	api.Delete("/:id", handler.DeleteMonitorAuxiliar)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedMonitoresAuxiliar)
}

func init() {
	registry.RegisterModule(RegisterRoutesMonitorAuxiliar)
}
