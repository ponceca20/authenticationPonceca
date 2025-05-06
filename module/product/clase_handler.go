package product

import (
	"time"
	"context"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ClaseHandler maneja solicitudes HTTP para Clase
type ClaseHandler struct {
	Service   *ClaseService
	validator *validator.Validate
}

// NewClaseHandler crea handler de Clase
func NewClaseHandler(s *ClaseService) *ClaseHandler {
	return &ClaseHandler{Service: s, validator: validator.New()}
}

// GetAllClases lista todas las clases
func (h *ClaseHandler) GetAllClases(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	clases, err := h.Service.GetAllClases(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar clases")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": clases, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// GetClase obtiene una clase por ID
func (h *ClaseHandler) GetClase(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	clase, err := h.Service.GetClase(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Clase no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": clase, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// CreateClase crea una nueva clase
func (h *ClaseHandler) CreateClase(c *fiber.Ctx) error {
	var clase Clase
	if err := c.BodyParser(&clase); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&clase); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validación fallida")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	clase.EmpresaID = empresaID
	if err := h.Service.CreateClase(&clase); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la clase")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": clase, "message": "Clase creada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// UpdateClase actualiza una clase existente
func (h *ClaseHandler) UpdateClase(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var clase Clase
	if err := c.BodyParser(&clase); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&clase); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validación fallida")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	clase.EmpresaID = empresaID
	updated, err := h.Service.UpdateClase(id, empresaID, &clase)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la clase")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": updated, "message": "Clase actualizada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// DeleteClase elimina una clase
func (h *ClaseHandler) DeleteClase(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteClase(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la clase")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Clase eliminada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// PatchClase actualiza parcialmente una clase
func (h *ClaseHandler) PatchClase(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	fields["empresa_id"] = empresaID
	clase, err := h.Service.PatchClase(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la clase")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": clase, "message": "Clase actualizada parcialmente", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// RegisterRoutesClase registra rutas de Clase
func RegisterRoutesClase(app *fiber.App) {
	repo := NewClaseRepository()
	ctx := context.Background()
	productoSearch, _ := NewProductoSearchService(ctx)
	service := NewClaseService(repo, productoSearch)
	handler := NewClaseHandler(service)
	api := app.Group("/api/v1/clases")
	api.Get("/", handler.GetAllClases)
	api.Get("/:id", handler.GetClase)
	api.Post("/", handler.CreateClase)
	api.Put("/:id", handler.UpdateClase)
	api.Patch("/:id", handler.PatchClase)
	api.Delete("/:id", handler.DeleteClase)
}

func init() {
	registry.RegisterModule(RegisterRoutesClase)
}
