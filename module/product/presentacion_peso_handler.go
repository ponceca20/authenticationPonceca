package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PresentacionPesoHandler maneja las solicitudes HTTP relacionadas con PresentacionPeso
type PresentacionPesoHandler struct {
	Service   *PresentacionPesoService
	validator *validator.Validate
}

// NewPresentacionPesoHandler crea una nueva instancia del handler
func NewPresentacionPesoHandler(s *PresentacionPesoService) *PresentacionPesoHandler {
	return &PresentacionPesoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllPresentacionesPeso obtiene todas las configuraciones de peso para una presentación
func (h *PresentacionPesoHandler) GetAllPresentacionesPeso(c *fiber.Ctx) error {
	presentacionIDStr := c.Query("presentacion_id")
	if presentacionIDStr == "" {
		return RespondWithError(c, fiber.StatusBadRequest, nil, "Se requiere el ID de la presentación")
	}

	presentacionID, err := strconv.ParseUint(presentacionIDStr, 10, 64)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de presentación inválido")
	}

	pesos, err := h.Service.GetAllPresentacionesPeso(presentacionID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar configuraciones de peso")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pesos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetPresentacionPeso obtiene una configuración de peso específica
func (h *PresentacionPesoHandler) GetPresentacionPeso(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	peso, err := h.Service.GetPresentacionPeso(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Configuración de peso no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      peso,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreatePresentacionPeso crea una nueva configuración de peso
func (h *PresentacionPesoHandler) CreatePresentacionPeso(c *fiber.Ctx) error {
	var peso PresentacionPeso
	if err := c.BodyParser(&peso); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&peso); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreatePresentacionPeso(&peso); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la configuración de peso")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      peso,
		"message":   "Configuración de peso creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdatePresentacionPeso actualiza una configuración de peso existente
func (h *PresentacionPesoHandler) UpdatePresentacionPeso(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var peso PresentacionPeso
	if err := c.BodyParser(&peso); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&peso); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	updated, err := h.Service.UpdatePresentacionPeso(id, &peso)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la configuración de peso")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Configuración de peso actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeletePresentacionPeso elimina una configuración de peso
func (h *PresentacionPesoHandler) DeletePresentacionPeso(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	if err := h.Service.DeletePresentacionPeso(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la configuración de peso")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Configuración de peso eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedPresentacionesPeso inicializa múltiples configuraciones de peso (para administradores)
func (h *PresentacionPesoHandler) SeedPresentacionesPeso(c *fiber.Ctx) error {
	var pesos []PresentacionPeso
	if err := c.BodyParser(&pesos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedPresentacionesPeso(pesos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las configuraciones de peso")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Configuraciones de peso inicializadas exitosamente",
		"count":     len(pesos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchPresentacionPeso actualiza parcialmente una configuración de peso
func (h *PresentacionPesoHandler) PatchPresentacionPeso(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	peso, err := h.Service.PatchPresentacionPeso(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la configuración de peso")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      peso,
		"message":   "Configuración de peso actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesPresentacionPeso registra las rutas para la API de PresentacionPeso
func RegisterRoutesPresentacionPeso(app *fiber.App) {
	repo := NewPresentacionPesoRepository()
	service := NewPresentacionPesoService(repo)
	handler := NewPresentacionPesoHandler(service)

	api := app.Group("/api/v1/presentaciones-peso")

	api.Get("/", handler.GetAllPresentacionesPeso)
	api.Get("/:id", handler.GetPresentacionPeso)
	api.Post("/", handler.CreatePresentacionPeso)
	api.Put("/:id", handler.UpdatePresentacionPeso)
	api.Patch("/:id", handler.PatchPresentacionPeso)
	api.Delete("/:id", handler.DeletePresentacionPeso)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedPresentacionesPeso)
}

func init() {
	registry.RegisterModule(RegisterRoutesPresentacionPeso)
}
