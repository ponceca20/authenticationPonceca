package auth

import (
	"time"

	"practicev2/registry" // Add this import

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TipoDocumentoSunatHandler maneja las solicitudes HTTP relacionadas con TipoDocumentoSunat.
type TipoDocumentoSunatHandler struct {
	Service   *TipoDocumentoSunatService
	validator *validator.Validate
}

func NewTipoDocumentoSunatHandler(s *TipoDocumentoSunatService) *TipoDocumentoSunatHandler {
	return &TipoDocumentoSunatHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *TipoDocumentoSunatHandler) GetAllTiposDocumento(c *fiber.Ctx) error {
	tiposDocumento, err := h.Service.GetAllTiposDocumento()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar tipos de documento")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tiposDocumento,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) GetTipoDocumento(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	tipoDocumento, err := h.Service.GetTipoDocumento(uint(id))
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Tipo de documento no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tipoDocumento,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) GetTipoDocumentoByCodigo(c *fiber.Ctx) error {
	codigo := c.Params("codigo")
	if codigo == "" {
		return RespondWithError(c, fiber.StatusBadRequest, nil, "Código inválido")
	}
	tipoDocumento, err := h.Service.GetTipoDocumentoByCodigo(codigo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Tipo de documento no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tipoDocumento,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) CreateTipoDocumento(c *fiber.Ctx) error {
	var tipoDocumento TipoDocumentoSunat
	if err := c.BodyParser(&tipoDocumento); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&tipoDocumento); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateTipoDocumento(&tipoDocumento); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el tipo de documento")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      tipoDocumento,
		"message":   "Tipo de documento creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) UpdateTipoDocumento(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var tipoDocumento TipoDocumentoSunat
	if err := c.BodyParser(&tipoDocumento); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&tipoDocumento); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdateTipoDocumento(uint(id), &tipoDocumento)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el tipo de documento")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Tipo de documento actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) DeleteTipoDocumento(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeleteTipoDocumento(uint(id)); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el tipo de documento")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Tipo de documento eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) SeedTiposDocumento(c *fiber.Ctx) error {
	var tiposDocumento []TipoDocumentoSunat
	if err := c.BodyParser(&tiposDocumento); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.Service.SeedTiposDocumento(tiposDocumento); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los tipos de documento")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Tipos de documento inicializados exitosamente",
		"count":     len(tiposDocumento),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *TipoDocumentoSunatHandler) PatchTipoDocumento(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	tipoDocumento, err := h.Service.PatchTipoDocumento(uint(id), fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el tipo de documento")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      tipoDocumento,
		"message":   "Tipo de documento actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesTipoDocumentoSunat registra las rutas del API para TipoDocumentoSunat.
func RegisterRoutesTipoDocumentoSunat(app *fiber.App) {
	repo := NewTipoDocumentoSunatRepository()
	service := NewTipoDocumentoSunatService(repo)
	h := NewTipoDocumentoSunatHandler(service)

	api := app.Group("/api/v1/tipos-documento")
	api.Get("/", h.GetAllTiposDocumento)
	api.Post("/seed", h.SeedTiposDocumento)
	api.Get("/:id", h.GetTipoDocumento)
	api.Get("/codigo/:codigo", h.GetTipoDocumentoByCodigo)
	api.Post("/", h.CreateTipoDocumento)
	api.Put("/:id", h.UpdateTipoDocumento)
	api.Patch("/:id", h.PatchTipoDocumento)
	api.Delete("/:id", h.DeleteTipoDocumento)
}

// Añadido init para registrar directamente este módulo en el registry.
func init() {
	registry.RegisterModule(RegisterRoutesTipoDocumentoSunat)
}
