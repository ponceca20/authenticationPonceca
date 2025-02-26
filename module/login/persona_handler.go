package users

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PersonaHandler maneja las solicitudes HTTP relacionadas con Persona.
type PersonaHandler struct {
	Service   *PersonaService
	validator *validator.Validate
}

func NewPersonaHandler(s *PersonaService) *PersonaHandler {
	return &PersonaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *PersonaHandler) GetAllPersonas(c *fiber.Ctx) error {
	personas, err := h.Service.GetAllPersonas()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar personas")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      personas,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) GetPersona(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	persona, err := h.Service.GetPersona(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Persona no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      persona,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) CreatePersona(c *fiber.Ctx) error {
	var persona Persona
	if err := c.BodyParser(&persona); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&persona); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreatePersona(&persona); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la persona")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      persona,
		"message":   "Persona creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) UpdatePersona(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var persona Persona
	if err := c.BodyParser(&persona); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&persona); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdatePersona(id, &persona)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la persona")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Persona actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) DeletePersona(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeletePersona(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la persona")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Persona eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) SeedPersonas(c *fiber.Ctx) error {
	var personas []Persona
	if err := c.BodyParser(&personas); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.Service.SeedPersonas(personas); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las personas")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Personas inicializadas exitosamente",
		"count":     len(personas),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PersonaHandler) PatchPersona(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	persona, err := h.Service.PatchPersona(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la persona")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      persona,
		"message":   "Persona actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesPersona registra las rutas del API para Persona.
func RegisterRoutesPersona(app *fiber.App) {
	repo := NewPersonaRepository()
	service := NewPersonaService(repo)
	h := NewPersonaHandler(service)

	api := app.Group("/api/v1/personas")
	api.Get("/", h.GetAllPersonas)
	api.Post("/seed", h.SeedPersonas)
	api.Get("/:id", h.GetPersona)
	api.Post("/", h.CreatePersona)
	api.Put("/:id", h.UpdatePersona)
	api.Patch("/:id", h.PatchPersona)
	api.Delete("/:id", h.DeletePersona)
}
