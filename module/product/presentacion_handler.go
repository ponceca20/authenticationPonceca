package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PresentacionHandler maneja las solicitudes HTTP relacionadas con Presentaciones
type PresentacionHandler struct {
	Service   *PresentacionService
	validator *validator.Validate
}

// NewPresentacionHandler crea una nueva instancia del handler
func NewPresentacionHandler(s *PresentacionService) *PresentacionHandler {
	return &PresentacionHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllPresentaciones obtiene todas las presentaciones de un producto
func (h *PresentacionHandler) GetAllPresentaciones(c *fiber.Ctx) error {
	// Obtener ID del producto desde los parámetros
	productoIDStr := c.Params("producto_id")
	productoID, err := strconv.ParseUint(productoIDStr, 10, 64)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	// Obtener las presentaciones
	presentaciones, err := h.Service.GetAllPresentaciones(productoID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar presentaciones")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      presentaciones,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetAllPresentacionesByEmpresa obtiene todas las presentaciones filtradas por empresa_id
func (h *PresentacionHandler) GetAllPresentacionesByEmpresa(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	presentaciones, err := h.Service.GetAllPresentacionesByEmpresa(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar presentaciones")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      presentaciones,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetPresentacion obtiene una presentación específica
func (h *PresentacionHandler) GetPresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	presentacion, err := h.Service.GetPresentacion(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Presentación no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      presentacion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreatePresentacion crea una nueva presentación
func (h *PresentacionHandler) CreatePresentacion(c *fiber.Ctx) error {
	var presentacion Presentacion
	if err := c.BodyParser(&presentacion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&presentacion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreatePresentacion(&presentacion); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la presentación")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      presentacion,
		"message":   "Presentación creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdatePresentacion actualiza una presentación existente
func (h *PresentacionHandler) UpdatePresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var presentacion Presentacion
	if err := c.BodyParser(&presentacion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&presentacion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	updated, err := h.Service.UpdatePresentacion(id, &presentacion)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la presentación")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Presentación actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeletePresentacion elimina una presentación
func (h *PresentacionHandler) DeletePresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	if err := h.Service.DeletePresentacion(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la presentación")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Presentación eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedPresentaciones inicializa múltiples presentaciones (para administradores)
func (h *PresentacionHandler) SeedPresentaciones(c *fiber.Ctx) error {
	var presentaciones []Presentacion
	if err := c.BodyParser(&presentaciones); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedPresentaciones(presentaciones); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las presentaciones")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Presentaciones inicializadas exitosamente",
		"count":     len(presentaciones),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchPresentacion actualiza parcialmente una presentación
func (h *PresentacionHandler) PatchPresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	presentacion, err := h.Service.PatchPresentacion(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la presentación")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      presentacion,
		"message":   "Presentación actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesPresentaciones registra las rutas para la API de Presentaciones
func RegisterRoutesPresentaciones(app *fiber.App) {
	repo := NewPresentacionRepository()
	service := NewPresentacionService(repo)
	handler := NewPresentacionHandler(service)

	// Ruta base para presentaciones
	api := app.Group("/api/v1")

	// Rutas para gestionar presentaciones por producto
	api.Get("/productos/:producto_id/presentaciones", handler.GetAllPresentaciones)

	// Ruta para obtener todas las presentaciones por empresa
	api.Get("/presentaciones", handler.GetAllPresentacionesByEmpresa)

	// Rutas para gestionar presentaciones directamente
	presentacionesAPI := api.Group("/presentaciones")
	presentacionesAPI.Get("/:id", handler.GetPresentacion)
	presentacionesAPI.Post("/", handler.CreatePresentacion)
	presentacionesAPI.Put("/:id", handler.UpdatePresentacion)
	presentacionesAPI.Patch("/:id", handler.PatchPresentacion)
	presentacionesAPI.Delete("/:id", handler.DeletePresentacion)

	// Ruta de administrador para inicializar datos
	presentacionesAPI.Post("/seed", handler.SeedPresentaciones)
}

func init() {
	registry.RegisterModule(RegisterRoutesPresentaciones)
}
