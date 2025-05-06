package product

import (
	"time"
	"context"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TipoProductoHandler maneja las solicitudes HTTP relacionadas con TipoProducto
type TipoProductoHandler struct {
	Service   *TipoProductoService
	validator *validator.Validate
}

// NewTipoProductoHandler crea una nueva instancia del handler
func NewTipoProductoHandler(s *TipoProductoService) *TipoProductoHandler {
	return &TipoProductoHandler{Service: s, validator: validator.New()}
}

func (h *TipoProductoHandler) GetAllTipoProductos(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	items, err := h.Service.GetAllTipoProductos(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar tipos de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": items, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *TipoProductoHandler) GetTipoProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	item, err := h.Service.GetTipoProducto(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Tipo de producto no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": item, "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *TipoProductoHandler) CreateTipoProducto(c *fiber.Ctx) error {
	var tp TipoProducto
	if err := c.BodyParser(&tp); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&tp); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	tp.EmpresaID = empresaID
	if err := h.Service.CreateTipoProducto(&tp); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el tipo de producto")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": tp, "message": "Tipo de producto creado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *TipoProductoHandler) UpdateTipoProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var tp TipoProducto
	if err := c.BodyParser(&tp); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&tp); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	tp.EmpresaID = empresaID
	updated, err := h.Service.UpdateTipoProducto(id, empresaID, &tp)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el tipo de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": updated, "message": "Tipo de producto actualizado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *TipoProductoHandler) PatchTipoProducto(c *fiber.Ctx) error {
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
	item, err := h.Service.PatchTipoProducto(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el tipo de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": item, "message": "Tipo de producto actualizado parcialmente", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *TipoProductoHandler) DeleteTipoProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteTipoProducto(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el tipo de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Tipo de producto eliminado", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

// RegisterRoutesTipoProducto registra las rutas para la API de TipoProducto
func RegisterRoutesTipoProducto(app *fiber.App) {
	repo := NewTipoProductoRepository()
	ctx := context.Background()
	productoSearch, _ := NewProductoSearchService(ctx)
	service := NewTipoProductoService(repo, productoSearch)
	handler := NewTipoProductoHandler(service)
	api := app.Group("/api/v1/tipos-productos")
	api.Get("/", handler.GetAllTipoProductos)
	api.Get("/:id", handler.GetTipoProducto)
	api.Post("/", handler.CreateTipoProducto)
	api.Put("/:id", handler.UpdateTipoProducto)
	api.Patch("/:id", handler.PatchTipoProducto)
	api.Delete("/:id", handler.DeleteTipoProducto)
}

func init() {
	registry.RegisterModule(RegisterRoutesTipoProducto)
}
