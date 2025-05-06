package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ListaPreciosPresentacionHandler maneja las solicitudes HTTP para ListaPreciosPresentacion.
type ListaPreciosPresentacionHandler struct {
	Service   *ListaPreciosPresentacionService
	validator *validator.Validate
}

// NewListaPreciosPresentacionHandler crea una nueva instancia del handler.
func NewListaPreciosPresentacionHandler(s *ListaPreciosPresentacionService) *ListaPreciosPresentacionHandler {
	return &ListaPreciosPresentacionHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ListaPreciosPresentacionHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.Service.GetAll()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar presentaciones")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      items,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ListaPreciosPresentacionHandler) GetByID(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	item, err := h.Service.GetByID(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Presentación no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ListaPreciosPresentacionHandler) Create(c *fiber.Ctx) error {
	var item ListaPreciosPresentacion
	if err := c.BodyParser(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	// Validar datos
	if err := h.validator.Struct(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.Create(&item); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la presentación")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"message":   "Presentación creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ListaPreciosPresentacionHandler) Update(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var item ListaPreciosPresentacion
	if err := c.BodyParser(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&item); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.Update(id, &item)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Presentación actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ListaPreciosPresentacionHandler) Delete(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.Delete(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Presentación eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ListaPreciosPresentacionHandler) Patch(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	item, err := h.Service.Patch(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la presentación")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"message":   "Presentación actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesListaPreciosPresentacion registra las rutas para la API.
func RegisterRoutesListaPreciosPresentacion(app *fiber.App) {
	repo := NewListaPreciosPresentacionRepository()
	service := NewListaPreciosPresentacionService(repo)
	handler := NewListaPreciosPresentacionHandler(service)

	api := app.Group("/api/v1/lista-precios-presentacion")
	api.Get("/", handler.GetAll)
	api.Get("/:id", handler.GetByID)
	api.Post("/", handler.Create)
	api.Put("/:id", handler.Update)
	api.Patch("/:id", handler.Patch)
	api.Delete("/:id", handler.Delete)

}
func init() {
	registry.RegisterModule(RegisterRoutesListaPreciosPresentacion)
}
