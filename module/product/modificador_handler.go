package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ModificadorHandler maneja las solicitudes HTTP para Modificador
type ModificadorHandler struct {
	Service   *ModificadorService
	validator *validator.Validate
}

func NewModificadorHandler(s *ModificadorService) *ModificadorHandler {
	return &ModificadorHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ModificadorHandler) GetAllModificadores(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	mods, err := h.Service.GetAllModificadores(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar modificadores")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mods,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) GetModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	mod, err := h.Service.GetModificador(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Modificador no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mod,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) CreateModificador(c *fiber.Ctx) error {
	var mod Modificador
	if err := c.BodyParser(&mod); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&mod); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateModificador(&mod); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el modificador")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      mod,
		"message":   "Modificador creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) UpdateModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var mod Modificador
	if err := c.BodyParser(&mod); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&mod); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	updated, err := h.Service.UpdateModificador(id, empresaID, &mod)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Modificador actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) DeleteModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteModificador(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Modificador eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) SeedModificadores(c *fiber.Ctx) error {
	var mods []Modificador
	if err := c.BodyParser(&mods); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedModificadores(mods); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los modificadores")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Modificadores inicializados exitosamente",
		"count":     len(mods),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModificadorHandler) PatchModificador(c *fiber.Ctx) error {
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
	mod, err := h.Service.PatchModificador(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el modificador")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      mod,
		"message":   "Modificador actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func RegisterRoutesModificador(app *fiber.App) {
	repo := NewModificadorRepository()
	service := NewModificadorService(repo)
	handler := NewModificadorHandler(service)
	api := app.Group("/api/v1/modificador")
	api.Get("/", handler.GetAllModificadores)
	api.Get("/:id", handler.GetModificador)
	api.Post("/", handler.CreateModificador)
	api.Put("/:id", handler.UpdateModificador)
	api.Patch("/:id", handler.PatchModificador)
	api.Delete("/:id", handler.DeleteModificador)
	api.Post("/seed", handler.SeedModificadores)
}

func init() {
	registry.RegisterModule(RegisterRoutesModificador)
}
