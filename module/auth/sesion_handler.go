package auth

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// SesionHandler maneja las solicitudes HTTP relacionadas con sesiones
type SesionHandler struct {
	Service   *SesionService
	validator *validator.Validate
}

// NewSesionHandler crea una nueva instancia del handler
func NewSesionHandler(s *SesionService) *SesionHandler {
	return &SesionHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *SesionHandler) GetAllSesions(c *fiber.Ctx) error {
	sesions, err := h.Service.GetAllSesions()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar sesiones")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      sesions,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) GetSesion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	sesion, err := h.Service.GetSesion(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Sesión no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      sesion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) CreateSesion(c *fiber.Ctx) error {
	var sesion Sesion
	if err := c.BodyParser(&sesion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&sesion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	if err := h.Service.CreateSesion(&sesion); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la sesión")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      sesion,
		"message":   "Sesión creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) UpdateSesion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var sesion Sesion
	if err := c.BodyParser(&sesion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&sesion); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	updatedSesion, err := h.Service.UpdateSesion(id, &sesion)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la sesión")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updatedSesion,
		"message":   "Sesión actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) DeleteSesion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	if err := h.Service.DeleteSesion(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la sesión")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Sesión eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) SeedSesions(c *fiber.Ctx) error {
	var sesions []Sesion
	if err := c.BodyParser(&sesions); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedSesions(sesions); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las sesiones")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Sesiones inicializadas exitosamente",
		"count":     len(sesions),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *SesionHandler) PatchSesion(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	sesion, err := h.Service.PatchSesion(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la sesión")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      sesion,
		"message":   "Sesión actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesSesion registra las rutas HTTP para sesiones
func RegisterRoutesSesion(app *fiber.App) {
	repo := NewSesionRepository()
	service := NewSesionService(repo)
	h := NewSesionHandler(service)

	api := app.Group("/api/v1/sesions")
	api.Get("/", h.GetAllSesions)
	api.Post("/seed", h.SeedSesions)
	api.Get("/:id", h.GetSesion)
	api.Post("/", h.CreateSesion)
	api.Put("/:id", h.UpdateSesion)
	api.Patch("/:id", h.PatchSesion)
	api.Delete("/:id", h.DeleteSesion)
}
