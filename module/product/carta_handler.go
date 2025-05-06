package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CartaHandler maneja las solicitudes HTTP relacionadas con Carta
type CartaHandler struct {
	Service   *CartaService
	validator *validator.Validate
}

// NewCartaHandler crea una nueva instancia del handler
func NewCartaHandler(s *CartaService) *CartaHandler {
	return &CartaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllCartas obtiene todas las cartas para la empresa del usuario autenticado
func (h *CartaHandler) GetAllCartas(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	cartas, err := h.Service.GetAllCartas(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar cartas")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      cartas,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetCarta obtiene una carta específica
func (h *CartaHandler) GetCarta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	carta, err := h.Service.GetCarta(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Carta no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      carta,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateCarta crea una nueva carta
func (h *CartaHandler) CreateCarta(c *fiber.Ctx) error {
	var carta Carta
	if err := c.BodyParser(&carta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&carta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	carta.EmpresaID = empresaID

	if err := h.Service.CreateCarta(&carta); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la carta")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      carta,
		"message":   "Carta creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateCarta actualiza una carta existente
func (h *CartaHandler) UpdateCarta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var carta Carta
	if err := c.BodyParser(&carta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&carta); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	carta.EmpresaID = empresaID

	updated, err := h.Service.UpdateCarta(id, empresaID, &carta)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la carta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Carta actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteCarta elimina una carta
func (h *CartaHandler) DeleteCarta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteCarta(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la carta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Carta eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedCartas inicializa múltiples cartas (para administradores)
func (h *CartaHandler) SeedCartas(c *fiber.Ctx) error {
	var cartas []Carta
	if err := c.BodyParser(&cartas); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedCartas(cartas); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las cartas")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Cartas inicializadas exitosamente",
		"count":     len(cartas),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchCarta actualiza parcialmente una carta
func (h *CartaHandler) PatchCarta(c *fiber.Ctx) error {
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

	// Proteger el campo empresaID contra modificaciones
	delete(fields, "empresa_id")

	carta, err := h.Service.PatchCarta(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la carta")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      carta,
		"message":   "Carta actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesCarta registra las rutas para la API de Carta
func RegisterRoutesCarta(app *fiber.App) {
	repo := NewCartaRepository()
	service := NewCartaService(repo)
	handler := NewCartaHandler(service)

	api := app.Group("/api/v1/cartas")

	api.Get("/", handler.GetAllCartas)
	api.Get("/:id", handler.GetCarta)
	api.Post("/", handler.CreateCarta)
	api.Put("/:id", handler.UpdateCarta)
	api.Patch("/:id", handler.PatchCarta)
	api.Delete("/:id", handler.DeleteCarta)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedCartas)
}

func init() {
	registry.RegisterModule(RegisterRoutesCarta)
}
