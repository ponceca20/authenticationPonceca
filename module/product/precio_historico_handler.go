package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// PrecioHistoricoHandler maneja las solicitudes HTTP para PrecioHistorico.
type PrecioHistoricoHandler struct {
	Service   *PrecioHistoricoService
	validator *validator.Validate
}

func NewPrecioHistoricoHandler(s *PrecioHistoricoService) *PrecioHistoricoHandler {
	return &PrecioHistoricoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *PrecioHistoricoHandler) GetAllPrecioHistorico(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	items, err := h.Service.GetAllPrecioHistorico(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error retrieving items")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      items,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioHistoricoHandler) GetPrecioHistorico(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	item, err := h.Service.GetPrecioHistorico(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Item not found")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      item,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioHistoricoHandler) CreatePrecioHistorico(c *fiber.Ctx) error {
	var ph PrecioHistorico
	if err := c.BodyParser(&ph); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	if err := h.validator.Struct(&ph); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	if err := h.Service.CreatePrecioHistorico(&ph); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not create item")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      ph,
		"message":   "Created successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioHistoricoHandler) UpdatePrecioHistorico(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var ph PrecioHistorico
	if err := c.BodyParser(&ph); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid body")
	}
	if err := h.validator.Struct(&ph); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation error")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	updated, err := h.Service.UpdatePrecioHistorico(id, empresaID, &ph)
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

func (h *PrecioHistoricoHandler) DeletePrecioHistorico(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	if err := h.Service.DeletePrecioHistorico(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Could not delete item")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Deleted successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *PrecioHistoricoHandler) PatchPrecioHistorico(c *fiber.Ctx) error {
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
	item, err := h.Service.PatchPrecioHistorico(id, empresaID, fields)
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

// RegisterRoutesPrecioHistorico registra las rutas para la API.
func RegisterRoutesPrecioHistorico(app *fiber.App) {
	repo := NewPrecioHistoricoRepository()
	service := NewPrecioHistoricoService(repo)
	handler := NewPrecioHistoricoHandler(service)

	api := app.Group("/api/v1/precio-historico")
	api.Get("/", handler.GetAllPrecioHistorico)
	api.Get("/:id", handler.GetPrecioHistorico)
	api.Post("/", handler.CreatePrecioHistorico)
	api.Put("/:id", handler.UpdatePrecioHistorico)
	api.Patch("/:id", handler.PatchPrecioHistorico)
	api.Delete("/:id", handler.DeletePrecioHistorico)

	registry.RegisterModule(RegisterRoutesPrecioHistorico)
}
