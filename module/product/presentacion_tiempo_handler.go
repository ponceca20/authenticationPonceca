package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PresentacionTiempoHandler maneja las solicitudes HTTP relacionadas con PresentacionTiempo
type PresentacionTiempoHandler struct {
	Service   *PresentacionTiempoService
	validator *validator.Validate
}

// NewPresentacionTiempoHandler crea una nueva instancia del handler
func NewPresentacionTiempoHandler(s *PresentacionTiempoService) *PresentacionTiempoHandler {
	return &PresentacionTiempoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllPresentacionesTiempo obtiene todas las configuraciones de tiempo para una presentación
func (h *PresentacionTiempoHandler) GetAllPresentacionesTiempo(c *fiber.Ctx) error {
	presentacionIDStr := c.Query("presentacion_id")
	if presentacionIDStr == "" {
		return RespondWithError(c, fiber.StatusBadRequest, nil, "Se requiere el ID de la presentación")
	}

	presentacionID, err := strconv.ParseUint(presentacionIDStr, 10, 64)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de presentación inválido")
	}

	tiempos, err := h.Service.GetAllPresentacionesTiempo(presentacionID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar configuraciones de tiempo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tiempos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetPresentacionTiempo obtiene una configuración de tiempo específica
func (h *PresentacionTiempoHandler) GetPresentacionTiempo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	tiempo, err := h.Service.GetPresentacionTiempo(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Configuración de tiempo no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tiempo,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreatePresentacionTiempo crea una nueva configuración de tiempo
func (h *PresentacionTiempoHandler) CreatePresentacionTiempo(c *fiber.Ctx) error {
	var tiempo PresentacionTiempo
	if err := c.BodyParser(&tiempo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&tiempo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreatePresentacionTiempo(&tiempo); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la configuración de tiempo")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      tiempo,
		"message":   "Configuración de tiempo creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdatePresentacionTiempo actualiza una configuración de tiempo existente
func (h *PresentacionTiempoHandler) UpdatePresentacionTiempo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var tiempo PresentacionTiempo
	if err := c.BodyParser(&tiempo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&tiempo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	updated, err := h.Service.UpdatePresentacionTiempo(id, &tiempo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la configuración de tiempo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Configuración de tiempo actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeletePresentacionTiempo elimina una configuración de tiempo
func (h *PresentacionTiempoHandler) DeletePresentacionTiempo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	if err := h.Service.DeletePresentacionTiempo(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la configuración de tiempo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Configuración de tiempo eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedPresentacionesTiempo inicializa múltiples configuraciones de tiempo (para administradores)
func (h *PresentacionTiempoHandler) SeedPresentacionesTiempo(c *fiber.Ctx) error {
	var tiempos []PresentacionTiempo
	if err := c.BodyParser(&tiempos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedPresentacionesTiempo(tiempos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las configuraciones de tiempo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Configuraciones de tiempo inicializadas exitosamente",
		"count":     len(tiempos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchPresentacionTiempo actualiza parcialmente una configuración de tiempo
func (h *PresentacionTiempoHandler) PatchPresentacionTiempo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	tiempo, err := h.Service.PatchPresentacionTiempo(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la configuración de tiempo")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tiempo,
		"message":   "Configuración de tiempo actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesPresentacionTiempo registra las rutas para la API de PresentacionTiempo
func RegisterRoutesPresentacionTiempo(app *fiber.App) {
	repo := NewPresentacionTiempoRepository()
	service := NewPresentacionTiempoService(repo)
	handler := NewPresentacionTiempoHandler(service)

	api := app.Group("/api/v1/presentaciones-tiempo")

	api.Get("/", handler.GetAllPresentacionesTiempo)
	api.Get("/:id", handler.GetPresentacionTiempo)
	api.Post("/", handler.CreatePresentacionTiempo)
	api.Put("/:id", handler.UpdatePresentacionTiempo)
	api.Patch("/:id", handler.PatchPresentacionTiempo)
	api.Delete("/:id", handler.DeletePresentacionTiempo)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedPresentacionesTiempo)
}

func init() {
	registry.RegisterModule(RegisterRoutesPresentacionTiempo)
}
