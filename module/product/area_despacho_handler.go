package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// AreaDespachoHandler maneja las solicitudes HTTP relacionadas con AreaDespacho
type AreaDespachoHandler struct {
	Service   *AreaDespachoService
	validator *validator.Validate
}

// NewAreaDespachoHandler crea una nueva instancia del handler
func NewAreaDespachoHandler(s *AreaDespachoService) *AreaDespachoHandler {
	return &AreaDespachoHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllAreasDespacho obtiene todas las áreas de despacho para la empresa del usuario autenticado
func (h *AreaDespachoHandler) GetAllAreasDespacho(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	areasDespacho, err := h.Service.GetAllAreasDespacho(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar áreas de despacho")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      areasDespacho,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetAreaDespacho obtiene un área de despacho específica
func (h *AreaDespachoHandler) GetAreaDespacho(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	areaDespacho, err := h.Service.GetAreaDespacho(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Área de despacho no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      areaDespacho,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateAreaDespacho crea una nueva área de despacho
func (h *AreaDespachoHandler) CreateAreaDespacho(c *fiber.Ctx) error {
	var areaDespacho AreaDespacho
	if err := c.BodyParser(&areaDespacho); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&areaDespacho); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	areaDespacho.EmpresaID = empresaID

	if err := h.Service.CreateAreaDespacho(&areaDespacho); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el área de despacho")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      areaDespacho,
		"message":   "Área de despacho creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateAreaDespacho actualiza un área de despacho existente
func (h *AreaDespachoHandler) UpdateAreaDespacho(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var areaDespacho AreaDespacho
	if err := c.BodyParser(&areaDespacho); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&areaDespacho); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	areaDespacho.EmpresaID = empresaID

	updated, err := h.Service.UpdateAreaDespacho(id, empresaID, &areaDespacho)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el área de despacho")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Área de despacho actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteAreaDespacho elimina un área de despacho
func (h *AreaDespachoHandler) DeleteAreaDespacho(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteAreaDespacho(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el área de despacho")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Área de despacho eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedAreasDespacho inicializa múltiples áreas de despacho (para administradores)
func (h *AreaDespachoHandler) SeedAreasDespacho(c *fiber.Ctx) error {
	var areasDespacho []AreaDespacho
	if err := c.BodyParser(&areasDespacho); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedAreasDespacho(areasDespacho); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las áreas de despacho")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Áreas de despacho inicializadas exitosamente",
		"count":     len(areasDespacho),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchAreaDespacho actualiza parcialmente un área de despacho
func (h *AreaDespachoHandler) PatchAreaDespacho(c *fiber.Ctx) error {
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
	fields["empresa_id"] = empresaID

	areaDespacho, err := h.Service.PatchAreaDespacho(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el área de despacho")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      areaDespacho,
		"message":   "Área de despacho actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesAreaDespacho registra las rutas para la API de AreaDespacho
func RegisterRoutesAreaDespacho(app *fiber.App) {
	repo := NewAreaDespachoRepository()
	service := NewAreaDespachoService(repo)
	handler := NewAreaDespachoHandler(service)

	api := app.Group("/api/v1/areas-despacho")

	api.Get("/", handler.GetAllAreasDespacho)
	api.Get("/:id", handler.GetAreaDespacho)
	api.Post("/", handler.CreateAreaDespacho)
	api.Put("/:id", handler.UpdateAreaDespacho)
	api.Patch("/:id", handler.PatchAreaDespacho)
	api.Delete("/:id", handler.DeleteAreaDespacho)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedAreasDespacho)
}
func init() {
	registry.RegisterModule(RegisterRoutesAreaDespacho)
}
