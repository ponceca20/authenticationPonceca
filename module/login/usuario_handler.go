package users

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UsuarioHandler maneja las solicitudes HTTP relacionadas con usuarios
type UsuarioHandler struct {
	Service   *UsuarioService
	validator *validator.Validate
}

// NewUsuarioHandler crea una nueva instancia del handler con validación
func NewUsuarioHandler(s *UsuarioService) *UsuarioHandler {
	v := validator.New()
	// Registro de validaciones personalizadas si es necesario (ejemplo: email único)
	return &UsuarioHandler{
		Service:   s,
		validator: v,
	}
}

// GetAllUsuarios obtiene todos los usuarios
func (h *UsuarioHandler) GetAllUsuarios(c *fiber.Ctx) error {
	usuarios, err := h.Service.GetAllUsuarios()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar usuarios")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      usuarios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetUsuario obtiene un usuario por ID
func (h *UsuarioHandler) GetUsuario(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	usuario, err := h.Service.GetUsuario(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Usuario no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      usuario,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateUsuario crea un nuevo usuario con validación
func (h *UsuarioHandler) CreateUsuario(c *fiber.Ctx) error {
	var usuario Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&usuario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	if err := h.Service.CreateUsuario(&usuario); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el usuario")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      usuario,
		"message":   "Usuario creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateUsuario actualiza un usuario existente
func (h *UsuarioHandler) UpdateUsuario(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var usuario Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&usuario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación de datos")
	}
	updatedUsuario, err := h.Service.UpdateUsuario(id, &usuario)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el usuario")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updatedUsuario,
		"message":   "Usuario actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteUsuario elimina un usuario
func (h *UsuarioHandler) DeleteUsuario(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	if err := h.Service.DeleteUsuario(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el usuario")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Usuario eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedUsuarios inicializa múltiples usuarios
func (h *UsuarioHandler) SeedUsuarios(c *fiber.Ctx) error {
	var usuarios []Usuario
	if err := c.BodyParser(&usuarios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedUsuarios(usuarios); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los usuarios")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Usuarios inicializados exitosamente",
		"count":     len(usuarios),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchUsuario actualiza parcialmente un usuario
func (h *UsuarioHandler) PatchUsuario(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido proporcionado")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	usuario, err := h.Service.PatchUsuario(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el usuario")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      usuario,
		"message":   "Usuario actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutes2 registra las rutas del API
func RegisterRoutes2(app *fiber.App) {
	repo := NewUsuarioRepository()
	service := NewUsuarioService(repo)
	h := NewUsuarioHandler(service)

	api := app.Group("/api/v1/users")
	api.Get("/", h.GetAllUsuarios)
	api.Post("/seed", h.SeedUsuarios)
	api.Get("/:id", h.GetUsuario)
	api.Post("/", h.CreateUsuario)
	api.Put("/:id", h.UpdateUsuario)
	api.Patch("/:id", h.PatchUsuario)
	api.Delete("/:id", h.DeleteUsuario)
}
