package auth

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Nuevo: instancia global del validador para evitar recrearlo en cada handler.
var validate = validator.New()

// Nuevo: función helper para obtener el timestamp formateado.
func getTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// RolModuloHandler maneja las solicitudes HTTP relacionadas con RolModulo
type RolModuloHandler struct {
	Service   *RolModuloService
	validator *validator.Validate
}

// Se modifica NewRolModuloHandler para usar el validador global.
func NewRolModuloHandler(s *RolModuloService) *RolModuloHandler {
	return &RolModuloHandler{
		Service:   s,
		validator: validate,
	}
}

// GetAllRolModulos obtiene todos los rolmodulos
func (h *RolModuloHandler) GetAllRolModulos(c *fiber.Ctx) error {
	rolmodulos, err := h.Service.GetAllRolModulos()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar rolmodulos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      rolmodulos,
		"timestamp": getTimestamp(),
	})
}

// GetRolModulo obtiene un rolmodulo por ID
func (h *RolModuloHandler) GetRolModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	rolmodulo, err := h.Service.GetRolModulo(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "RolModulo no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      rolmodulo,
		"timestamp": getTimestamp(),
	})
}

// CreateRolModulo crea un nuevo rolmodulo con validación
func (h *RolModuloHandler) CreateRolModulo(c *fiber.Ctx) error {
	var rolmodulo RolModulo
	if err := c.BodyParser(&rolmodulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&rolmodulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	if err := h.Service.CreateRolModulo(&rolmodulo); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el rolmodulo")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      rolmodulo,
		"message":   "RolModulo creado exitosamente",
		"timestamp": getTimestamp(),
	})
}

// UpdateRolModulo actualiza un rolmodulo existente
func (h *RolModuloHandler) UpdateRolModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var rolmodulo RolModulo
	if err := c.BodyParser(&rolmodulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&rolmodulo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	updatedRolModulo, err := h.Service.UpdateRolModulo(id, &rolmodulo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el rolmodulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updatedRolModulo,
		"message":   "RolModulo actualizado exitosamente",
		"timestamp": getTimestamp(),
	})
}

// DeleteRolModulo elimina un rolmodulo
func (h *RolModuloHandler) DeleteRolModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	if err := h.Service.DeleteRolModulo(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el rolmodulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "RolModulo eliminado exitosamente",
		"timestamp": getTimestamp(),
	})
}

// SeedRolModulos inicializa múltiples rolmodulos
func (h *RolModuloHandler) SeedRolModulos(c *fiber.Ctx) error {
	var rolmodulos []RolModulo
	if err := c.BodyParser(&rolmodulos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedRolModulos(rolmodulos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los rolmodulos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "RolModulos inicializados exitosamente",
		"count":     len(rolmodulos),
		"timestamp": getTimestamp(),
	})
}

// PatchRolModulo actualiza parcialmente un rolmodulo
func (h *RolModuloHandler) PatchRolModulo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	rolmodulo, err := h.Service.PatchRolModulo(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el rolmodulo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      rolmodulo,
		"message":   "RolModulo actualizado parcialmente con éxito",
		"timestamp": getTimestamp(),
	})
}

// RegisterRoutesRolModulo registra las rutas del API para rolmodulos
func RegisterRoutesRolModulo(app *fiber.App) {
	repo := NewRolModuloRepository()
	service := NewRolModuloService(repo)
	h := NewRolModuloHandler(service)

	api := app.Group("/api/v1/rolmodulos")
	api.Get("/", h.GetAllRolModulos)
	api.Post("/seed", h.SeedRolModulos)
	api.Get("/:id", h.GetRolModulo)
	api.Post("/", h.CreateRolModulo)
	api.Put("/:id", h.UpdateRolModulo)
	api.Patch("/:id", h.PatchRolModulo)
	api.Delete("/:id", h.DeleteRolModulo)
}
