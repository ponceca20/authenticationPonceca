package users

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ModuloHandler maneja las solicitudes HTTP relacionadas con módulos
type ModuloHandler struct {
	Service   *ModuloService
	validator *validator.Validate
}

// NewModuloHandler crea una nueva instancia del handler
func NewModuloHandler(s *ModuloService) *ModuloHandler {
	return &ModuloHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ModuloHandler) GetAllModulos(c *fiber.Ctx) error {
	modulos, err := h.Service.GetAllModulos()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar módulos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      modulos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) GetModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	modulo, err := h.Service.GetModulo(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Módulo no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      modulo,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) CreateModulo(c *fiber.Ctx) error {
	var modulo Modulo
	if err := c.BodyParser(&modulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&modulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	if err := h.Service.CreateModulo(&modulo); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el módulo")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      modulo,
		"message":   "Módulo creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) UpdateModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var modulo Modulo
	if err := c.BodyParser(&modulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&modulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	updatedModulo, err := h.Service.UpdateModulo(id, &modulo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el módulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updatedModulo,
		"message":   "Módulo actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) DeleteModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	if err := h.Service.DeleteModulo(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el módulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Módulo eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) SeedModulos(c *fiber.Ctx) error {
	var modulos []Modulo
	if err := c.BodyParser(&modulos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedModulos(modulos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los módulos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Módulos inicializados exitosamente",
		"count":     len(modulos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ModuloHandler) PatchModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	modulo, err := h.Service.PatchModulo(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el módulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      modulo,
		"message":   "Módulo actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesModulo registra las rutas HTTP para módulos
func RegisterRoutesModulo(app *fiber.App) {
	repo := NewModuloRepository()
	service := NewModuloService(repo)
	h := NewModuloHandler(service)

	api := app.Group("/api/v1/modulos")
	api.Get("/", h.GetAllModulos)
	api.Post("/seed", h.SeedModulos)
	api.Get("/:id", h.GetModulo)
	api.Post("/", h.CreateModulo)
	api.Put("/:id", h.UpdateModulo)
	api.Patch("/:id", h.PatchModulo)
	api.Delete("/:id", h.DeleteModulo)
}
