package product

import (
	"time"
	"context"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// SubCategoriaHandler gestiona las solicitudes de SubCategoria
type SubCategoriaHandler struct {
	Service   *SubCategoriaService
	validator *validator.Validate
}

// NewSubCategoriaHandler crea una nueva instancia de SubCategoriaHandler
func NewSubCategoriaHandler(s *SubCategoriaService) *SubCategoriaHandler {
	return &SubCategoriaHandler{Service: s, validator: validator.New()}
}

// GetAllSubCategorias lista todas las subcategorias
func (h *SubCategoriaHandler) GetAllSubCategorias(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	subs, err := h.Service.GetAllSubCategorias(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar subcategorias")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": subs, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// GetSubCategoria obtiene una SubCategoria por ID
func (h *SubCategoriaHandler) GetSubCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	sub, err := h.Service.GetSubCategoria(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "SubCategoria no encontrada")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": sub, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// CreateSubCategoria crea una nueva SubCategoria
func (h *SubCategoriaHandler) CreateSubCategoria(c *fiber.Ctx) error {
	var sub SubCategoria
	if err := c.BodyParser(&sub); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&sub); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validación fallida")
	}
	// verify authentication but do not assign EmpresaID (no such field on SubCategoria)
	if _, err := GetCompanyID(c); err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.CreateSubCategoria(&sub); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la subcategoria")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": sub, "message": "SubCategoria creada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// UpdateSubCategoria actualiza una SubCategoria existente
func (h *SubCategoriaHandler) UpdateSubCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var sub SubCategoria
	if err := c.BodyParser(&sub); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	if err := h.validator.Struct(&sub); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Validación fallida")
	}
	// verify authentication (context) but no EmpresaID field
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	updated, err := h.Service.UpdateSubCategoria(id, empresaID, &sub)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la subcategoria")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": updated, "message": "SubCategoria actualizada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// PatchSubCategoria actualiza parcialmente una SubCategoria
func (h *SubCategoriaHandler) PatchSubCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo inválido")
	}
	// En los handlers de subcategoria, ignora empresa_id en los PATCH y no lo pases al servicio ni al repo.
	// Si hay lógica que obtiene empresa_id del contexto, solo úsala para categoria, no subcategoria.
	sub, err := h.Service.PatchSubCategoria(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la subcategoria")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": sub, "message": "SubCategoria actualizada parcialmente", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// DeleteSubCategoria elimina una SubCategoria
func (h *SubCategoriaHandler) DeleteSubCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteSubCategoria(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la subcategoria")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "SubCategoria eliminada", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// RegisterRoutesSubCategoria registra las rutas para la API de SubCategoria
func RegisterRoutesSubCategoria(app *fiber.App) {
	repo := NewSubCategoriaRepository()
	ctx := context.Background()
	productoSearch, _ := NewProductoSearchService(ctx)
	service := NewSubCategoriaService(repo, productoSearch)
	handler := NewSubCategoriaHandler(service)
	api := app.Group("/api/v1/subcategorias")
	api.Get("/", handler.GetAllSubCategorias)
	api.Get("/:id", handler.GetSubCategoria)
	api.Post("/", handler.CreateSubCategoria)
	api.Put("/:id", handler.UpdateSubCategoria)
	api.Patch("/:id", handler.PatchSubCategoria)
	api.Delete("/:id", handler.DeleteSubCategoria)
}

func init() {
	registry.RegisterModule(RegisterRoutesSubCategoria)
}
