package auth

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UsuarioEmpresaHandler maneja las solicitudes HTTP para UsuarioEmpresa.
type UsuarioEmpresaHandler struct {
	Service   *UsuarioEmpresaService
	validator *validator.Validate
}

func NewUsuarioEmpresaHandler(s *UsuarioEmpresaService) *UsuarioEmpresaHandler {
	return &UsuarioEmpresaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *UsuarioEmpresaHandler) GetAllUsuarioEmpresas(c *fiber.Ctx) error {
	ues, err := h.Service.GetAllUsuarioEmpresas()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar usuario empresas")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      ues,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) GetUsuarioEmpresa(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	ue, err := h.Service.GetUsuarioEmpresa(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Usuario empresa no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      ue,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) CreateUsuarioEmpresa(c *fiber.Ctx) error {
	var ue UsuarioEmpresa
	if err := c.BodyParser(&ue); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&ue); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateUsuarioEmpresa(&ue); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el usuario empresa")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      ue,
		"message":   "Usuario empresa creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) UpdateUsuarioEmpresa(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var ue UsuarioEmpresa
	if err := c.BodyParser(&ue); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&ue); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdateUsuarioEmpresa(id, &ue)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el usuario empresa")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Usuario empresa actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) DeleteUsuarioEmpresa(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeleteUsuarioEmpresa(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la usuario empresa")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Usuario empresa eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) SeedUsuarioEmpresas(c *fiber.Ctx) error {
	var ues []UsuarioEmpresa
	if err := c.BodyParser(&ues); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.Service.SeedUsuarioEmpresas(ues); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar usuario empresas")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Usuario empresas inicializadas exitosamente",
		"count":     len(ues),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *UsuarioEmpresaHandler) PatchUsuarioEmpresa(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	ue, err := h.Service.PatchUsuarioEmpresa(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el usuario empresa")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      ue,
		"message":   "Usuario empresa actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesUsuarioEmpresa registra las rutas del API para UsuarioEmpresa.
func RegisterRoutesUsuarioEmpresa(app *fiber.App) {
	repo := NewUsuarioEmpresaRepository()
	service := NewUsuarioEmpresaService(repo)
	h := NewUsuarioEmpresaHandler(service)

	api := app.Group("/api/v1/usuarioempresas")
	api.Get("/", h.GetAllUsuarioEmpresas)
	api.Post("/seed", h.SeedUsuarioEmpresas)
	api.Get("/:id", h.GetUsuarioEmpresa)
	api.Post("/", h.CreateUsuarioEmpresa)
	api.Put("/:id", h.UpdateUsuarioEmpresa)
	api.Patch("/:id", h.PatchUsuarioEmpresa)
	api.Delete("/:id", h.DeleteUsuarioEmpresa)
}
