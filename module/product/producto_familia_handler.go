package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductoFamiliaHandler maneja las solicitudes HTTP relacionadas con ProductoFamilia
type ProductoFamiliaHandler struct {
	Service   *ProductoFamiliaService
	validator *validator.Validate
}

// NewProductoFamiliaHandler crea una nueva instancia del handler
func NewProductoFamiliaHandler(s *ProductoFamiliaService) *ProductoFamiliaHandler {
	return &ProductoFamiliaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllProductoFamilias obtiene todas las relaciones ProductoFamilia para la empresa del usuario autenticado
func (h *ProductoFamiliaHandler) GetAllProductoFamilias(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoFamilias, err := h.Service.GetAllProductoFamilias(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar relaciones producto-familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetProductoFamiliasByFamiliaID obtiene todas las relaciones ProductoFamilia para una familia específica
func (h *ProductoFamiliaHandler) GetProductoFamiliasByFamiliaID(c *fiber.Ctx) error {
	familiaID, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de familia inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoFamilias, err := h.Service.GetProductoFamiliasByFamiliaID(familiaID, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar productos para la familia especificada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetProductoFamiliasByProductoID obtiene todas las relaciones ProductoFamilia para un producto específico
func (h *ProductoFamiliaHandler) GetProductoFamiliasByProductoID(c *fiber.Ctx) error {
	productoID, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoFamilias, err := h.Service.GetProductoFamiliasByProductoID(productoID, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar familias para el producto especificado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetProductoFamilia obtiene una relación ProductoFamilia específica
func (h *ProductoFamiliaHandler) GetProductoFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoFamilia, err := h.Service.GetProductoFamilia(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Relación producto-familia no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilia,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateProductoFamilia crea una nueva relación ProductoFamilia
func (h *ProductoFamiliaHandler) CreateProductoFamilia(c *fiber.Ctx) error {
	var productoFamilia ProductoFamilia
	if err := c.BodyParser(&productoFamilia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&productoFamilia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreateProductoFamilia(&productoFamilia); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la relación producto-familia")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilia,
		"message":   "Relación producto-familia creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateProductoFamilia actualiza una relación ProductoFamilia existente
func (h *ProductoFamiliaHandler) UpdateProductoFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var productoFamilia ProductoFamilia
	if err := c.BodyParser(&productoFamilia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&productoFamilia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	updated, err := h.Service.UpdateProductoFamilia(id, empresaID, &productoFamilia)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la relación producto-familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Relación producto-familia actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteProductoFamilia elimina una relación ProductoFamilia
func (h *ProductoFamiliaHandler) DeleteProductoFamilia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteProductoFamilia(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la relación producto-familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Relación producto-familia eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedProductoFamilias inicializa múltiples relaciones ProductoFamilia (para administradores)
func (h *ProductoFamiliaHandler) SeedProductoFamilias(c *fiber.Ctx) error {
	var productoFamilias []ProductoFamilia
	if err := c.BodyParser(&productoFamilias); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedProductoFamilias(productoFamilias); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las relaciones producto-familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Relaciones producto-familia inicializadas exitosamente",
		"count":     len(productoFamilias),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchProductoFamilia actualiza parcialmente una relación ProductoFamilia
func (h *ProductoFamiliaHandler) PatchProductoFamilia(c *fiber.Ctx) error {
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

	productoFamilia, err := h.Service.PatchProductoFamilia(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la relación producto-familia")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoFamilia,
		"message":   "Relación producto-familia actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesProductoFamilia registra las rutas para la API de ProductoFamilia
func RegisterRoutesProductoFamilia(app *fiber.App) {
	repo := NewProductoFamiliaRepository()
	service := NewProductoFamiliaService(repo)
	handler := NewProductoFamiliaHandler(service)

	api := app.Group("/api/v1/productos-familias")

	api.Get("/", handler.GetAllProductoFamilias)
	api.Get("/:id", handler.GetProductoFamilia)
	api.Get("/familia/:id", handler.GetProductoFamiliasByFamiliaID)
	api.Get("/producto/:id", handler.GetProductoFamiliasByProductoID)
	api.Post("/", handler.CreateProductoFamilia)
	api.Put("/:id", handler.UpdateProductoFamilia)
	api.Patch("/:id", handler.PatchProductoFamilia)
	api.Delete("/:id", handler.DeleteProductoFamilia)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedProductoFamilias)
}

func init() {
	registry.RegisterModule(RegisterRoutesProductoFamilia)
}
