package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductoModificadorHandler struct {
	Service   *ProductoModificadorService
	validator *validator.Validate
}

func NewProductoModificadorHandler(s *ProductoModificadorService) *ProductoModificadorHandler {
	return &ProductoModificadorHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// getProductoIDFromQuery extrae y valida el producto_id de la query.
// Retorna (productoID, true) si existe y es válido, o (0, false) si no está presente.
func getProductoIDFromQuery(c *fiber.Ctx) (uint64, bool, error) {
	productoIDStr := c.Query("producto_id")
	if productoIDStr == "" {
		return 0, false, nil // No está presente, no es error
	}
	productoID, err := strconv.ParseUint(productoIDStr, 10, 64)
	if err != nil {
		return 0, false, err
	}
	return productoID, true, nil
}

func (h *ProductoModificadorHandler) GetAllProductoModificadores(c *fiber.Ctx) error {
	productoID, hasProductoID, err := getProductoIDFromQuery(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	var pms []ProductoModificador
	if hasProductoID {
		pms, err = h.Service.GetAllProductoModificadores(productoID)
		if err != nil {
			return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar modificadores de producto")
		}
	} else {
		// Si no se pasa producto_id, retorna todos los modificadores (opcional, comenta si no lo deseas)
		pms, err = h.Service.GetAllProductoModificadores(0)
		if err != nil {
			return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar todos los modificadores de producto")
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pms,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) GetProductoModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	pm, err := h.Service.GetProductoModificador(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Modificador de producto no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pm,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) CreateProductoModificador(c *fiber.Ctx) error {
	var pm ProductoModificador
	if err := c.BodyParser(&pm); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&pm); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateProductoModificador(&pm); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el modificador de producto")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      pm,
		"message":   "Modificador de producto creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) UpdateProductoModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var pm ProductoModificador
	if err := c.BodyParser(&pm); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&pm); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	updated, err := h.Service.UpdateProductoModificador(id, &pm)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el modificador de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Modificador de producto actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) DeleteProductoModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	if err := h.Service.DeleteProductoModificador(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el modificador de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Modificador de producto eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) SeedProductoModificadores(c *fiber.Ctx) error {
	var pms []ProductoModificador
	if err := c.BodyParser(&pms); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedProductoModificadores(pms); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los modificadores de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Modificadores de producto inicializados exitosamente",
		"count":     len(pms),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *ProductoModificadorHandler) PatchProductoModificador(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	pm, err := h.Service.PatchProductoModificador(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el modificador de producto")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      pm,
		"message":   "Modificador de producto actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func RegisterRoutesProductoModificador(app *fiber.App) {
	repo := NewProductoModificadorRepository()
	service := NewProductoModificadorService(repo)
	handler := NewProductoModificadorHandler(service)

	api := app.Group("/api/v1/producto-modificadores")

	api.Get("/", handler.GetAllProductoModificadores)
	api.Get("/:id", handler.GetProductoModificador)
	api.Post("/", handler.CreateProductoModificador)
	api.Put("/:id", handler.UpdateProductoModificador)
	api.Patch("/:id", handler.PatchProductoModificador)
	api.Delete("/:id", handler.DeleteProductoModificador)
	api.Post("/seed", handler.SeedProductoModificadores)
}

func init() {
	registry.RegisterModule(RegisterRoutesProductoModificador)
}
