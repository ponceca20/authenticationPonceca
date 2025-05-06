package product

import (
	"practicev2/registry"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ImpuestoHandler maneja las solicitudes HTTP relacionadas con Impuesto
type ImpuestoHandler struct {
	Service   *ImpuestoService
	validator *validator.Validate
}

func NewImpuestoHandler(s *ImpuestoService) *ImpuestoHandler {
	return &ImpuestoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *ImpuestoHandler) GetAllImpuestos(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	impuestos, err := h.Service.GetAllImpuestos(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error retrieving impuestos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      impuestos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) GetImpuesto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	impuesto, err := h.Service.GetImpuesto(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Impuesto not found")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      impuesto,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) CreateImpuesto(c *fiber.Ctx) error {
	var impuesto Impuesto
	if err := c.BodyParser(&impuesto); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}
	if err := h.validator.Struct(&impuesto); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation failed")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	impuesto.EmpresaID = empresaID
	if err := h.Service.CreateImpuesto(&impuesto); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error creating impuesto")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      impuesto,
		"message":   "Impuesto created successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) UpdateImpuesto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var impuesto Impuesto
	if err := c.BodyParser(&impuesto); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}
	if err := h.validator.Struct(&impuesto); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validation failed")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	impuesto.EmpresaID = empresaID
	updated, err := h.Service.UpdateImpuesto(id, empresaID, &impuesto)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error updating impuesto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Impuesto updated successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) DeleteImpuesto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	if err := h.Service.DeleteImpuesto(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error deleting impuesto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Impuesto deleted successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) SeedImpuestos(c *fiber.Ctx) error {
	var impuestos []Impuesto
	if err := c.BodyParser(&impuestos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}
	if err := h.Service.SeedImpuestos(impuestos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error seeding impuestos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Impuestos seeded successfully",
		"count":     len(impuestos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ImpuestoHandler) PatchImpuesto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid ID")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Invalid request body")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Authentication error")
	}
	fields["empresa_id"] = empresaID
	impuesto, err := h.Service.PatchImpuesto(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error patching impuesto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      impuesto,
		"message":   "Impuesto patched successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func RegisterRoutesImpuesto(app *fiber.App) {
	repo := NewImpuestoRepository()
	service := NewImpuestoService(repo)
	handler := NewImpuestoHandler(service)

	api := app.Group("/api/v1/impuestos")
	api.Get("/", handler.GetAllImpuestos)
	api.Get("/:id", handler.GetImpuesto)
	api.Post("/", handler.CreateImpuesto)
	api.Put("/:id", handler.UpdateImpuesto)
	api.Patch("/:id", handler.PatchImpuesto)
	api.Delete("/:id", handler.DeleteImpuesto)
	api.Post("/seed", handler.SeedImpuestos)
}

func init() {
	registry.RegisterModule(RegisterRoutesImpuesto)
}
