package users

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// RolHandler maneja las solicitudes HTTP relacionadas con Rol.
type RolHandler struct {
	Service   *RolService
	validator *validator.Validate
}

// NewRolHandler crea una nueva instancia del handler.
func NewRolHandler(s *RolService) *RolHandler {
	return &RolHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *RolHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.Service.GetAllRoles()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar roles")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      roles,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) GetRol(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	rol, err := h.Service.GetRol(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Rol no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      rol,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) CreateRol(c *fiber.Ctx) error {
	var rol Rol
	if err := c.BodyParser(&rol); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&rol); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateRol(&rol); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el rol")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      rol,
		"message":   "Rol creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) UpdateRol(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var rol Rol
	if err := c.BodyParser(&rol); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&rol); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdateRol(id, &rol)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el rol")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Rol actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) DeleteRol(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeleteRol(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el rol")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Rol eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) SeedRoles(c *fiber.Ctx) error {
	var roles []Rol
	if err := c.BodyParser(&roles); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.Service.SeedRoles(roles); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los roles")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Roles inicializados exitosamente",
		"count":     len(roles),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *RolHandler) PatchRol(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	rol, err := h.Service.PatchRol(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el rol")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      rol,
		"message":   "Rol actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesRol registra las rutas del API para Rol.
func RegisterRoutesRol(app *fiber.App) {
	repo := NewRolRepository()
	service := NewRolService(repo)
	h := NewRolHandler(service)

	api := app.Group("/api/v1/roles")
	api.Get("/", h.GetAllRoles)
	api.Post("/seed", h.SeedRoles)
	api.Get("/:id", h.GetRol)
	api.Post("/", h.CreateRol)
	api.Put("/:id", h.UpdateRol)
	api.Patch("/:id", h.PatchRol)
	api.Delete("/:id", h.DeleteRol)
}
