package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PrecioDescuentoHandler maneja las solicitudes HTTP para PrecioDescuentoPorCantidad.
type PrecioDescuentoHandler struct {
	Service   *PrecioDescuentoService
	validator *validator.Validate
}

func NewPrecioDescuentoHandler(s *PrecioDescuentoService) *PrecioDescuentoHandler {
	return &PrecioDescuentoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *PrecioDescuentoHandler) GetAllPrecioDescuento(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	items, err := h.Service.GetAllPrecioDescuento(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error retrieving items")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      items,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioDescuentoHandler) GetPrecioDescuento(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	item, err := h.Service.GetPrecioDescuento(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Item not found")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioDescuentoHandler) CreatePrecioDescuento(c *fiber.Ctx) error {
	var pd PrecioDescuentoPorCantidad
	if err := c.BodyParser(&pd); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	if err := h.validator.Struct(&pd); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	if err := h.Service.CreatePrecioDescuento(&pd); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not create item")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      pd,
		"message":   "Created successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioDescuentoHandler) UpdatePrecioDescuento(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var pd PrecioDescuentoPorCantidad
	if err := c.BodyParser(&pd); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	if err := h.validator.Struct(&pd); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	updated, err := h.Service.UpdatePrecioDescuento(id, empresaID, &pd)
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

func (h *PrecioDescuentoHandler) DeletePrecioDescuento(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	if err := h.Service.DeletePrecioDescuento(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not delete item")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Deleted successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioDescuentoHandler) PatchPrecioDescuento(c *fiber.Ctx) error {
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
	fields["empresa_id"] = empresaID
	item, err := h.Service.PatchPrecioDescuento(id, empresaID, fields)
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

// RegisterRoutesPrecioDescuento registra las rutas para la API.
func RegisterRoutesPrecioDescuento(app *fiber.App) {
	repo := NewPrecioDescuentoRepository()
	service := NewPrecioDescuentoService(repo)
	handler := NewPrecioDescuentoHandler(service)

	api := app.Group("/api/v1/precio-descuento")
	api.Get("/", handler.GetAllPrecioDescuento)
	api.Get("/:id", handler.GetPrecioDescuento)
	api.Post("/", handler.CreatePrecioDescuento)
	api.Put("/:id", handler.UpdatePrecioDescuento)
	api.Patch("/:id", handler.PatchPrecioDescuento)
	api.Delete("/:id", handler.DeletePrecioDescuento)

	registry.RegisterModule(RegisterRoutesPrecioDescuento)
}
