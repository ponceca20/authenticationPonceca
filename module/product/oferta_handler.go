package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// OfertaHandler maneja las solicitudes HTTP relacionadas con Oferta
type OfertaHandler struct {
	Service   *OfertaService
	validator *validator.Validate
}

// NewOfertaHandler crea una nueva instancia del handler de Oferta
func NewOfertaHandler(s *OfertaService) *OfertaHandler {
	return &OfertaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllOfertas obtiene todas las ofertas para la empresa del usuario autenticado
func (h *OfertaHandler) GetAllOfertas(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	ofertas, err := h.Service.GetAllOfertas(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar ofertas")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      ofertas,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetOferta obtiene una oferta específica
func (h *OfertaHandler) GetOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	oferta, err := h.Service.GetOferta(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Oferta no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      oferta,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateOferta crea una nueva oferta
func (h *OfertaHandler) CreateOferta(c *fiber.Ctx) error {
	var oferta Oferta
	if err := c.BodyParser(&oferta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.validator.Struct(&oferta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	oferta.EmpresaID = empresaID

	if err := h.Service.CreateOferta(&oferta); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la oferta")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      oferta,
		"message":   "Oferta creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateOferta actualiza una oferta existente
func (h *OfertaHandler) UpdateOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var oferta Oferta
	if err := c.BodyParser(&oferta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.validator.Struct(&oferta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	oferta.EmpresaID = empresaID

	updated, err := h.Service.UpdateOferta(id, empresaID, &oferta)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la oferta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Oferta actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteOferta elimina una oferta
func (h *OfertaHandler) DeleteOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteOferta(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la oferta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Oferta eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedOfertas inicializa múltiples ofertas (para administradores)
func (h *OfertaHandler) SeedOfertas(c *fiber.Ctx) error {
	var ofertas []Oferta
	if err := c.BodyParser(&ofertas); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedOfertas(ofertas); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las ofertas")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Ofertas inicializadas exitosamente",
		"count":     len(ofertas),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchOferta actualiza parcialmente una oferta
func (h *OfertaHandler) PatchOferta(c *fiber.Ctx) error {
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
	// Proteger el campo empresa_id
	fields["empresa_id"] = empresaID

	oferta, err := h.Service.PatchOferta(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la oferta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      oferta,
		"message":   "Oferta actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesOferta registra las rutas para la API de Oferta
func RegisterRoutesOferta(app *fiber.App) {
	repo := NewOfertaRepository()
	service := NewOfertaService(repo)
	handler := NewOfertaHandler(service)

	api := app.Group("/api/v1/ofertas")
	api.Get("/", handler.GetAllOfertas)
	api.Get("/:id", handler.GetOferta)
	api.Post("/", handler.CreateOferta)
	api.Put("/:id", handler.UpdateOferta)
	api.Patch("/:id", handler.PatchOferta)
	api.Delete("/:id", handler.DeleteOferta)
	api.Post("/seed", handler.SeedOfertas)
}

func init() {
	registry.RegisterModule(RegisterRoutesOferta)
}
