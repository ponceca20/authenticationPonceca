package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// GrupoConteoHandler maneja las solicitudes HTTP para GrupoConteo.
type GrupoConteoHandler struct {
	Service   *GrupoConteoService
	validator *validator.Validate
}

// NewGrupoConteoHandler crea una nueva instancia del handler.
func NewGrupoConteoHandler(s *GrupoConteoService) *GrupoConteoHandler {
	return &GrupoConteoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *GrupoConteoHandler) GetAll(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	items, err := h.Service.GetAll(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar grupos de conteo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      items,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoConteoHandler) GetByID(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	item, err := h.Service.GetByID(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Grupo de conteo no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoConteoHandler) Create(c *fiber.Ctx) error {
	var item GrupoConteo
	if err := c.BodyParser(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	item.EmpresaID = empresaID
	if err := h.Service.Create(&item); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el grupo de conteo")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"message":   "Grupo de conteo creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoConteoHandler) Update(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var item GrupoConteo
	if err := c.BodyParser(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	item.EmpresaID = empresaID
	updated, err := h.Service.Update(id, empresaID, &item)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el grupo de conteo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Grupo de conteo actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoConteoHandler) Delete(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.Delete(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el grupo de conteo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Grupo de conteo eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoConteoHandler) Patch(c *fiber.Ctx) error {
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
	fields["empresa_id"] = empresaID
	item, err := h.Service.Patch(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el grupo de conteo")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"message":   "Grupo de conteo actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesGrupoConteo registra las rutas para la API.
func RegisterRoutesGrupoConteo(app *fiber.App) {
	repo := NewGrupoConteoRepository()
	service := NewGrupoConteoService(repo)
	handler := NewGrupoConteoHandler(service)

	api := app.Group("/api/v1/grupo-conteo")
	api.Get("/", handler.GetAll)
	api.Get("/:id", handler.GetByID)
	api.Post("/", handler.Create)
	api.Put("/:id", handler.Update)
	api.Patch("/:id", handler.Patch)
	api.Delete("/:id", handler.Delete)

	registry.RegisterModule(RegisterRoutesGrupoConteo)
}
