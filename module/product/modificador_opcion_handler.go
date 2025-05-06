package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ModificadorOpcionHandler maneja las solicitudes HTTP para ModificadorOpcion
type ModificadorOpcionHandler struct {
	Service   *ModificadorOpcionService
	validator *validator.Validate
}

func NewModificadorOpcionHandler(s *ModificadorOpcionService) *ModificadorOpcionHandler {
	return &ModificadorOpcionHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ModificadorOpcionHandler) GetAllModificadorOpciones(c *fiber.Ctx) error {
	modificadorID := c.Query("modificador_id")
	if modificadorID == "" {
		return RespondWithError(c, fiber.StatusBadRequest, nil, "Se requiere el ID del modificador")
	}
	modID, err := strconv.ParseUint(modificadorID, 10, 64)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de modificador inválido")
	}
	mos, err := h.Service.GetAllModificadorOpciones(modID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar opciones de modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) GetModificadorOpcion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	mo, err := h.Service.GetModificadorOpcion(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Opción de modificador no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mo,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) CreateModificadorOpcion(c *fiber.Ctx) error {
	var mo ModificadorOpcion
	if err := c.BodyParser(&mo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&mo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateModificadorOpcion(&mo); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la opción de modificador")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      mo,
		"message":   "Opción de modificador creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) UpdateModificadorOpcion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var mo ModificadorOpcion
	if err := c.BodyParser(&mo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&mo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdateModificadorOpcion(id, &mo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la opción de modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Opción de modificador actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) DeleteModificadorOpcion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeleteModificadorOpcion(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la opción de modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Opción de modificador eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) PatchModificadorOpcion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	mo, err := h.Service.PatchModificadorOpcion(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la opción de modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mo,
		"message":   "Opción de modificador actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorOpcionHandler) SeedModificadorOpciones(c *fiber.Ctx) error {
	var mos []ModificadorOpcion
	if err := c.BodyParser(&mos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedModificadorOpciones(mos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las opciones de modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Opciones de modificador inicializadas exitosamente",
		"count":     len(mos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func RegisterRoutesModificadorOpcion(app *fiber.App) {
	repo := NewModificadorOpcionRepository()
	service := NewModificadorOpcionService(repo)
	handler := NewModificadorOpcionHandler(service)
	api := app.Group("/api/v1/modificador-opcion")
	api.Get("/", handler.GetAllModificadorOpciones)
	api.Get("/:id", handler.GetModificadorOpcion)
	api.Post("/", handler.CreateModificadorOpcion)
	api.Put("/:id", handler.UpdateModificadorOpcion)
	api.Patch("/:id", handler.PatchModificadorOpcion)
	api.Delete("/:id", handler.DeleteModificadorOpcion)
	api.Post("/seed", handler.SeedModificadorOpciones)
}

func init() {
	registry.RegisterModule(RegisterRoutesModificadorOpcion)
}
