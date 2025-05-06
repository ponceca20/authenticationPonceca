package product

import (
	"time"
	"context"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductoHandler maneja las solicitudes HTTP relacionadas con Producto
type ProductoHandler struct {
	Service   *ProductoService
	validator *validator.Validate
}

// NewProductoHandler crea una nueva instancia del handler
func NewProductoHandler(s *ProductoService) *ProductoHandler {
	return &ProductoHandler{Service: s, validator: validator.New()}
}

// GetAllProductos obtiene todos los productos para la empresa autenticada
func (h *ProductoHandler) GetAllProductos(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	items, err := h.Service.GetAllProductos(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar productos")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": items, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// GetProducto obtiene un producto específico por ID
func (h *ProductoHandler) GetProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	item, err := h.Service.GetProducto(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Producto no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": item, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// CreateProducto crea un nuevo producto
func (h *ProductoHandler) CreateProducto(c *fiber.Ctx) error {
	var p Producto
	if err := c.BodyParser(&p); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&p); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	p.EmpresaID = empresaID
	if err := h.Service.CreateProducto(&p); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el producto")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": p, "message": "Producto creado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// UpdateProducto actualiza un producto existente
func (h *ProductoHandler) UpdateProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var p Producto
	if err := c.BodyParser(&p); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&p); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	p.EmpresaID = empresaID
	updated, err := h.Service.UpdateProducto(id, empresaID, &p)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": updated, "message": "Producto actualizado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// DeleteProducto elimina un producto
func (h *ProductoHandler) DeleteProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteProducto(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Producto eliminado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// PatchProducto actualiza parcialmente un producto
func (h *ProductoHandler) PatchProducto(c *fiber.Ctx) error {
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
	item, err := h.Service.PatchProducto(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": item, "message": "Producto actualizado parcialmente", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// RegisterRoutesProducto registra las rutas para la API de Producto
func RegisterRoutesProducto(app *fiber.App) {
	repo := NewProductoRepository()
	ctx := context.Background()
	svcSearch, _ := NewProductoSearchService(ctx) // Ignora error para fallback limpio
	service := NewProductoService(repo, svcSearch)
	handler := NewProductoHandler(service)
	api := app.Group("/api/v1/productos")
	api.Get("/", handler.GetAllProductos)
	api.Get("/:id", handler.GetProducto)
	api.Post("/", handler.CreateProducto)
	api.Put("/:id", handler.UpdateProducto)
	api.Patch("/:id", handler.PatchProducto)
	api.Delete("/:id", handler.DeleteProducto)
}

func init() {
	registry.RegisterModule(RegisterRoutesProducto)
}
