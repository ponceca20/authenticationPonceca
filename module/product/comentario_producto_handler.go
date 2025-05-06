package product

import (
	"strconv"
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ComentarioProductoHandler maneja las solicitudes HTTP relacionadas con ComentarioProducto
type ComentarioProductoHandler struct {
	Service   *ComentarioProductoService
	validator *validator.Validate
}

// NewComentarioProductoHandler crea una nueva instancia del handler
func NewComentarioProductoHandler(s *ComentarioProductoService) *ComentarioProductoHandler {
	return &ComentarioProductoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllComentariosProducto obtiene todos los comentarios para un producto específico o todos si no se especifica producto
func (h *ComentarioProductoHandler) GetAllComentariosProducto(c *fiber.Ctx) error {
	// Verificar si se proporciona un ID de producto como parámetro de consulta
	productoIDStr := c.Query("producto_id")
	var productoID *uint64

	if productoIDStr != "" {
		id, err := strconv.ParseUint(productoIDStr, 10, 64)
		if err != nil {
			return RespondWithError(c, fiber.StatusBadRequest, err, "ID de producto inválido")
		}
		productoID = &id
	}

	comentarios, err := h.Service.GetAllComentariosProducto(productoID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar comentarios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      comentarios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetComentariosGenericos obtiene todos los comentarios genéricos
func (h *ComentarioProductoHandler) GetComentariosGenericos(c *fiber.Ctx) error {
	comentarios, err := h.Service.GetComentariosGenericos()
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar comentarios genéricos")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      comentarios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetComentarioProducto obtiene un comentario específico
func (h *ComentarioProductoHandler) GetComentarioProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	comentario, err := h.Service.GetComentarioProducto(id)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Comentario no encontrado")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      comentario,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateComentarioProducto crea un nuevo comentario
func (h *ComentarioProductoHandler) CreateComentarioProducto(c *fiber.Ctx) error {
	var comentario ComentarioProducto
	if err := c.BodyParser(&comentario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&comentario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	if err := h.Service.CreateComentarioProducto(&comentario); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el comentario")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      comentario,
		"message":   "Comentario creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateComentarioProducto actualiza un comentario existente
func (h *ComentarioProductoHandler) UpdateComentarioProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var comentario ComentarioProducto
	if err := c.BodyParser(&comentario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&comentario); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	updated, err := h.Service.UpdateComentarioProducto(id, &comentario)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el comentario")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Comentario actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteComentarioProducto elimina un comentario
func (h *ComentarioProductoHandler) DeleteComentarioProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	if err := h.Service.DeleteComentarioProducto(id); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el comentario")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Comentario eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedComentariosProducto inicializa múltiples comentarios (para administradores)
func (h *ComentarioProductoHandler) SeedComentariosProducto(c *fiber.Ctx) error {
	var comentarios []ComentarioProducto
	if err := c.BodyParser(&comentarios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedComentariosProducto(comentarios); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los comentarios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Comentarios inicializados exitosamente",
		"count":     len(comentarios),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchComentarioProducto actualiza parcialmente un comentario
func (h *ComentarioProductoHandler) PatchComentarioProducto(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var fields map[string]interface{}
	if err := c.BodyParser(&fields); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	comentario, err := h.Service.PatchComentarioProducto(id, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el comentario")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      comentario,
		"message":   "Comentario actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesComentarioProducto registra las rutas para la API de ComentarioProducto
func RegisterRoutesComentarioProducto(app *fiber.App) {
	repo := NewComentarioProductoRepository()
	service := NewComentarioProductoService(repo)
	handler := NewComentarioProductoHandler(service)

	api := app.Group("/api/v1/comentarios")

	api.Get("/", handler.GetAllComentariosProducto)
	api.Get("/genericos", handler.GetComentariosGenericos)
	api.Get("/:id", handler.GetComentarioProducto)
	api.Post("/", handler.CreateComentarioProducto)
	api.Put("/:id", handler.UpdateComentarioProducto)
	api.Patch("/:id", handler.PatchComentarioProducto)
	api.Delete("/:id", handler.DeleteComentarioProducto)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedComentariosProducto)
}

func init() {
	registry.RegisterModule(RegisterRoutesComentarioProducto)
}
