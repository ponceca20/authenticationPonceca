package product

import (
	"context"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CategoriaHandler maneja las solicitudes HTTP relacionadas con Categoria
type CategoriaHandler struct {
	Service   *CategoriaService
	validator *validator.Validate
}

// NewCategoriaHandler crea una nueva instancia del handler
func NewCategoriaHandler(s *CategoriaService) *CategoriaHandler {
	return &CategoriaHandler{Service: s, validator: validator.New()}
}

// GetAllCategorias obtiene todas las categorias para la empresa del usuario autenticado
func (h *CategoriaHandler) GetAllCategorias(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	cats, err := h.Service.GetAllCategorias(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar categorias")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      cats,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetCategoria obtiene una categoria específica
func (h *CategoriaHandler) GetCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	cat, err := h.Service.GetCategoria(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Categoria no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      cat,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateCategoria crea una nueva categoria
func (h *CategoriaHandler) CreateCategoria(c *fiber.Ctx) error {
	var cat Categoria
	if err := c.BodyParser(&cat); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.validator.Struct(&cat); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	cat.EmpresaID = empresaID

	if err := h.Service.CreateCategoria(&cat); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la categoria")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      cat,
		"message":   "Categoria creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateCategoria actualiza una categoria existente
func (h *CategoriaHandler) UpdateCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var cat Categoria
	if err := c.BodyParser(&cat); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.validator.Struct(&cat); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	cat.EmpresaID = empresaID

	updated, err := h.Service.UpdateCategoria(id, empresaID, &cat)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la categoria")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Categoria actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchCategoria actualiza parcialmente una categoria
func (h *CategoriaHandler) PatchCategoria(c *fiber.Ctx) error {
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
	// Asegurar companyID en campos si viene
	fields["empresa_id"] = empresaID

	cat, err := h.Service.PatchCategoria(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la categoria")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      cat,
		"message":   "Categoria actualizada parcialmente exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteCategoria elimina una categoria
func (h *CategoriaHandler) DeleteCategoria(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteCategoria(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la categoria")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Categoria eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesCategoria registra las rutas para la API de Categoria
func RegisterRoutesCategoria(app *fiber.App) {
	repo := NewCategoriaRepository()
	ctx := context.Background()
	productoSearch, _ := NewProductoSearchService(ctx)
	service := NewCategoriaService(repo, productoSearch)
	handler := NewCategoriaHandler(service)

	api := app.Group("/api/v1/categorias")

	api.Get("/", handler.GetAllCategorias)
	api.Get("/:id", handler.GetCategoria)
	api.Post("/", handler.CreateCategoria)
	api.Put("/:id", handler.UpdateCategoria)
	api.Patch("/:id", handler.PatchCategoria)
	api.Delete("/:id", handler.DeleteCategoria)
}

func init() {
	registry.RegisterModule(RegisterRoutesCategoria)
}
