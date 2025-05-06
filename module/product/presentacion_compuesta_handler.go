package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PresentacionCompuestaHandler maneja las solicitudes HTTP relacionadas con PresentacionCompuesta
type PresentacionCompuestaHandler struct {
	Service   *PresentacionCompuestaService
	validator *validator.Validate
}

func NewPresentacionCompuestaHandler(s *PresentacionCompuestaService) *PresentacionCompuestaHandler {
	return &PresentacionCompuestaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *PresentacionCompuestaHandler) GetAllPresentacionCompuestas(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	pcs, err := h.Service.GetAllPresentacionCompuestas(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar presentacion compuesta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pcs,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) GetPresentacionCompuesta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	pc, err := h.Service.GetPresentacionCompuesta(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Presentacion compuesta no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pc,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) CreatePresentacionCompuesta(c *fiber.Ctx) error {
	var pc PresentacionCompuesta
	if err := c.BodyParser(&pc); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&pc); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreatePresentacionCompuesta(&pc); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la presentacion compuesta")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      pc,
		"message":   "Presentacion compuesta creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) UpdatePresentacionCompuesta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var pc PresentacionCompuesta
	if err := c.BodyParser(&pc); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&pc); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	updated, err := h.Service.UpdatePresentacionCompuesta(id, empresaID, &pc)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la presentacion compuesta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Presentacion compuesta actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) DeletePresentacionCompuesta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeletePresentacionCompuesta(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la presentacion compuesta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Presentacion compuesta eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) SeedPresentacionCompuestas(c *fiber.Ctx) error {
	var pcs []PresentacionCompuesta
	if err := c.BodyParser(&pcs); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedPresentacionCompuestas(pcs); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las presentacion compuesta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Presentacion compuesta inicializada exitosamente",
		"count":     len(pcs),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PresentacionCompuestaHandler) PatchPresentacionCompuesta(c *fiber.Ctx) error {
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
	pc, err := h.Service.PatchPresentacionCompuesta(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la presentacion compuesta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pc,
		"message":   "Presentacion compuesta actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesPresentacionCompuesta registra las rutas para PresentacionCompuesta
func RegisterRoutesPresentacionCompuesta(app *fiber.App) {
	repo := NewPresentacionCompuestaRepository()
	service := NewPresentacionCompuestaService(repo)
	handler := NewPresentacionCompuestaHandler(service)

	api := app.Group("/api/v1/presentacion-compuesta")

	api.Get("/", handler.GetAllPresentacionCompuestas)
	api.Get("/:id", handler.GetPresentacionCompuesta)
	api.Post("/", handler.CreatePresentacionCompuesta)
	api.Put("/:id", handler.UpdatePresentacionCompuesta)
	api.Patch("/:id", handler.PatchPresentacionCompuesta)
	api.Delete("/:id", handler.DeletePresentacionCompuesta)
	api.Post("/seed", handler.SeedPresentacionCompuestas)
}

func init() {
	registry.RegisterModule(RegisterRoutesPresentacionCompuesta)
}
