package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// OfertaPresentacionHandler maneja las solicitudes HTTP relacionadas con OfertaPresentacion
type OfertaPresentacionHandler struct {
	Service   *OfertaPresentacionService
	validator *validator.Validate
}

func NewOfertaPresentacionHandler(s *OfertaPresentacionService) *OfertaPresentacionHandler {
	return &OfertaPresentacionHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *OfertaPresentacionHandler) GetAllOfertaPresentaciones(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	ops, err := h.Service.GetAllOfertaPresentaciones(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar ofertas presentaciones")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      ops,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) GetOfertaPresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	op, err := h.Service.GetOfertaPresentacion(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Oferta presentación no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      op,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) CreateOfertaPresentacion(c *fiber.Ctx) error {
	var op OfertaPresentacion
	if err := c.BodyParser(&op); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&op); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateOfertaPresentacion(&op); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la oferta presentación")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      op,
		"message":   "Oferta presentación creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) UpdateOfertaPresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var op OfertaPresentacion
	if err := c.BodyParser(&op); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&op); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	updated, err := h.Service.UpdateOfertaPresentacion(id, empresaID, &op)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la oferta presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Oferta presentación actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) DeleteOfertaPresentacion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteOfertaPresentacion(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la oferta presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Oferta presentación eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) SeedOfertaPresentaciones(c *fiber.Ctx) error {
	var ops []OfertaPresentacion
	if err := c.BodyParser(&ops); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedOfertaPresentaciones(ops); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las ofertas presentaciones")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Ofertas presentaciones inicializadas exitosamente",
		"count":     len(ops),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *OfertaPresentacionHandler) PatchOfertaPresentacion(c *fiber.Ctx) error {
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
	op, err := h.Service.PatchOfertaPresentacion(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la oferta presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      op,
		"message":   "Oferta presentación actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesOfertaPresentacion registra las rutas para OfertaPresentacion
func RegisterRoutesOfertaPresentacion(app *fiber.App) {
	repo := NewOfertaPresentacionRepository()
	service := NewOfertaPresentacionService(repo)
	handler := NewOfertaPresentacionHandler(service)
	api := app.Group("/api/v1/ofertas-presentacion")
	api.Get("/", handler.GetAllOfertaPresentaciones)
	api.Get("/:id", handler.GetOfertaPresentacion)
	api.Post("/", handler.CreateOfertaPresentacion)
	api.Put("/:id", handler.UpdateOfertaPresentacion)
	api.Patch("/:id", handler.PatchOfertaPresentacion)
	api.Delete("/:id", handler.DeleteOfertaPresentacion)
	api.Post("/seed", handler.SeedOfertaPresentaciones)
}

func init() {
	registry.RegisterModule(RegisterRoutesOfertaPresentacion)
}
