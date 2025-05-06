package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductoGrupoConteoHandler maneja las solicitudes HTTP para ProductoGrupoConteo.
type ProductoGrupoConteoHandler struct {
	Service   *ProductoGrupoConteoService
	validator *validator.Validate
}

func NewProductoGrupoConteoHandler(s *ProductoGrupoConteoService) *ProductoGrupoConteoHandler {
	return &ProductoGrupoConteoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ProductoGrupoConteoHandler) GetAllProductoGrupoConteo(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	items, err := h.Service.GetAllProductoGrupoConteo(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error retrieving items")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      items,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoGrupoConteoHandler) GetProductoGrupoConteo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	item, err := h.Service.GetProductoGrupoConteo(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Item not found")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoGrupoConteoHandler) CreateProductoGrupoConteo(c *fiber.Ctx) error {
	var pg ProductoGrupoConteo
	if err := c.BodyParser(&pg); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	// Validar datos.
	if err := h.validator.Struct(&pg); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	// Se asume que la relación con Empresa se gestiona a través del producto.
	if err := h.Service.CreateProductoGrupoConteo(&pg); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not create item")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      pg,
		"message":   "Created successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoGrupoConteoHandler) UpdateProductoGrupoConteo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var pg ProductoGrupoConteo
	if err := c.BodyParser(&pg); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	if err := h.validator.Struct(&pg); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	updated, err := h.Service.UpdateProductoGrupoConteo(id, empresaID, &pg)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not update item")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Updated successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoGrupoConteoHandler) DeleteProductoGrupoConteo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	if err := h.Service.DeleteProductoGrupoConteo(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not delete item")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Deleted successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoGrupoConteoHandler) PatchProductoGrupoConteo(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	// Proteger posibles modificaciones al campo empresa_id.
	fields["empresa_id"] = empresaID
	item, err := h.Service.PatchProductoGrupoConteo(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not patch item")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"message":   "Patched successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesProductoGrupoConteo registra las rutas para la API.
func RegisterRoutesProductoGrupoConteo(app *fiber.App) {
	repo := NewProductoGrupoConteoRepository()
	service := NewProductoGrupoConteoService(repo)
	handler := NewProductoGrupoConteoHandler(service)

	api := app.Group("/api/v1/producto-grupo-conteo")
	api.Get("/", handler.GetAllProductoGrupoConteo)
	api.Get("/:id", handler.GetProductoGrupoConteo)
	api.Post("/", handler.CreateProductoGrupoConteo)
	api.Put("/:id", handler.UpdateProductoGrupoConteo)
	api.Patch("/:id", handler.PatchProductoGrupoConteo)
	api.Delete("/:id", handler.DeleteProductoGrupoConteo)

	registry.RegisterModule(RegisterRoutesProductoGrupoConteo)
}
