package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductoMediaHandler maneja las solicitudes HTTP relacionadas con ProductoMedia
type ProductoMediaHandler struct {
	Service   *ProductoMediaService
	validator *validator.Validate
}

// NewProductoMediaHandler crea una nueva instancia del handler
func NewProductoMediaHandler(s *ProductoMediaService) *ProductoMediaHandler {
	return &ProductoMediaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllProductoMedias obtiene todos los medios para la empresa del usuario autenticado
func (h *ProductoMediaHandler) GetAllProductoMedias(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoMedias, err := h.Service.GetAllProductoMedias(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar medios de productos")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoMedias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetProductoMediasByProductoID obtiene todos los medios para un producto específico
func (h *ProductoMediaHandler) GetProductoMediasByProductoID(c *fiber.Ctx) error {
	productoID, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoMedias, err := h.Service.GetProductoMediasByProductoID(productoID, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar medios para el producto especificado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoMedias,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetProductoMedia obtiene un medio específico
func (h *ProductoMediaHandler) GetProductoMedia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	productoMedia, err := h.Service.GetProductoMedia(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Medio no encontrado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoMedia,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateProductoMedia crea un nuevo medio para un producto
func (h *ProductoMediaHandler) CreateProductoMedia(c *fiber.Ctx) error {
	var productoMedia ProductoMedia
	if err := c.BodyParser(&productoMedia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&productoMedia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreateProductoMedia(&productoMedia); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el medio para el producto")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      productoMedia,
		"message":   "Medio creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateProductoMedia actualiza un medio existente
func (h *ProductoMediaHandler) UpdateProductoMedia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var productoMedia ProductoMedia
	if err := c.BodyParser(&productoMedia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&productoMedia); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	updated, err := h.Service.UpdateProductoMedia(id, empresaID, &productoMedia)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el medio")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Medio actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteProductoMedia elimina un medio
func (h *ProductoMediaHandler) DeleteProductoMedia(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteProductoMedia(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el medio")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Medio eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedProductoMedias inicializa múltiples medios (para administradores)
func (h *ProductoMediaHandler) SeedProductoMedias(c *fiber.Ctx) error {
	var productoMedias []ProductoMedia
	if err := c.BodyParser(&productoMedias); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedProductoMedias(productoMedias); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los medios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Medios inicializados exitosamente",
		"count":     len(productoMedias),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchProductoMedia actualiza parcialmente un medio
func (h *ProductoMediaHandler) PatchProductoMedia(c *fiber.Ctx) error {
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

	productoMedia, err := h.Service.PatchProductoMedia(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el medio")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      productoMedia,
		"message":   "Medio actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SetMainImage establece una imagen como principal para un producto
func (h *ProductoMediaHandler) SetMainImage(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	productoID, err := c.ParamsInt("productoId")
	if err != nil || productoID <= 0 {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.SetMainImage(uint64(id), uint64(productoID), empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo establecer la imagen como principal")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Imagen principal establecida exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReorderMedia reordena los medios de un producto
func (h *ProductoMediaHandler) ReorderMedia(c *fiber.Ctx) error {
	productoID, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
	}

	var requestBody struct {
		MediaIDs []uint64 `json:"media_ids"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if len(requestBody.MediaIDs) == 0 {
		return RespondWithError(c, fiber.StatusBadRequest, nil, "La lista de IDs de medios está vacía")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.ReorderMedia(productoID, empresaID, requestBody.MediaIDs); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron reordenar los medios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Medios reordenados exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesProductoMedia registra las rutas para la API de ProductoMedia
func RegisterRoutesProductoMedia(app *fiber.App) {
	repo := NewProductoMediaRepository()
	service := NewProductoMediaService(repo)
	handler := NewProductoMediaHandler(service)

	api := app.Group("/api/v1/producto-medias")

	api.Get("/", handler.GetAllProductoMedias)
	api.Get("/:id", handler.GetProductoMedia)
	api.Get("/producto/:id", handler.GetProductoMediasByProductoID)
	api.Post("/", handler.CreateProductoMedia)
	api.Put("/:id", handler.UpdateProductoMedia)
	api.Patch("/:id", handler.PatchProductoMedia)
	api.Delete("/:id", handler.DeleteProductoMedia)
	api.Post("/principal/:id/producto/:productoId", handler.SetMainImage)
	api.Post("/reordenar/:id", handler.ReorderMedia)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedProductoMedias)
}

func init() {
	registry.RegisterModule(RegisterRoutesProductoMedia)
}
