package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ListaPreciosHandler maneja las solicitudes HTTP relacionadas con ListaPrecios
type ListaPreciosHandler struct {
	Service   *ListaPreciosService
	validator *validator.Validate
}

// NewListaPreciosHandler crea una nueva instancia del handler
func NewListaPreciosHandler(s *ListaPreciosService) *ListaPreciosHandler {
	return &ListaPreciosHandler{
		Service:   s,
		validator: validator.New(),
	}
}

// GetAllListaPrecios obtiene todas las listas de precios para la empresa del usuario autenticado
func (h *ListaPreciosHandler) GetAllListaPrecios(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	listaPrecios, err := h.Service.GetAllListaPrecios(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar listas de precios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      listaPrecios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetListaPreciosWithPresentaciones obtiene una lista de precios con sus presentaciones asociadas
func (h *ListaPreciosHandler) GetListaPreciosWithPresentaciones(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	listaPrecios, err := h.Service.GetListaPreciosWithPresentaciones(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Lista de precios no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      listaPrecios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// GetListaPrecios obtiene una lista de precios específica
func (h *ListaPreciosHandler) GetListaPrecios(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	listaPrecios, err := h.Service.GetListaPrecios(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Lista de precios no encontrada")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      listaPrecios,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// CreateListaPrecios crea una nueva lista de precios
func (h *ListaPreciosHandler) CreateListaPrecios(c *fiber.Ctx) error {
	var listaPrecios ListaPrecios
	if err := c.BodyParser(&listaPrecios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&listaPrecios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	// Establecer el EmpresaID desde el contexto
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	listaPrecios.EmpresaID = empresaID

	if err := h.Service.CreateListaPrecios(&listaPrecios); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear la lista de precios")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      listaPrecios,
		"message":   "Lista de precios creada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// UpdateListaPrecios actualiza una lista de precios existente
func (h *ListaPreciosHandler) UpdateListaPrecios(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	var listaPrecios ListaPrecios
	if err := c.BodyParser(&listaPrecios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	// Validar datos
	if err := h.validator.Struct(&listaPrecios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	// Asegurar que no se manipula el EmpresaID
	listaPrecios.EmpresaID = empresaID

	updated, err := h.Service.UpdateListaPrecios(id, empresaID, &listaPrecios)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar la lista de precios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Lista de precios actualizada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DeleteListaPrecios elimina una lista de precios
func (h *ListaPreciosHandler) DeleteListaPrecios(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}

	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}

	if err := h.Service.DeleteListaPrecios(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar la lista de precios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Lista de precios eliminada exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SeedListaPrecios inicializa múltiples listas de precios (para administradores)
func (h *ListaPreciosHandler) SeedListaPrecios(c *fiber.Ctx) error {
	var listaPrecios []ListaPrecios
	if err := c.BodyParser(&listaPrecios); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}

	if err := h.Service.SeedListaPrecios(listaPrecios); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar las listas de precios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Listas de precios inicializadas exitosamente",
		"count":     len(listaPrecios),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// PatchListaPrecios actualiza parcialmente una lista de precios
func (h *ListaPreciosHandler) PatchListaPrecios(c *fiber.Ctx) error {
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

	listaPrecios, err := h.Service.PatchListaPrecios(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente la lista de precios")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      listaPrecios,
		"message":   "Lista de precios actualizada parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesListaPrecios registra las rutas para la API de ListaPrecios
func RegisterRoutesListaPrecios(app *fiber.App) {
	repo := NewListaPreciosRepository()
	service := NewListaPreciosService(repo)
	handler := NewListaPreciosHandler(service)

	api := app.Group("/api/v1/listas-precios")

	api.Get("/", handler.GetAllListaPrecios)
	api.Get("/:id", handler.GetListaPrecios)
	api.Get("/:id/presentaciones", handler.GetListaPreciosWithPresentaciones)
	api.Post("/", handler.CreateListaPrecios)
	api.Put("/:id", handler.UpdateListaPrecios)
	api.Patch("/:id", handler.PatchListaPrecios)
	api.Delete("/:id", handler.DeleteListaPrecios)

	// Ruta de administrador para inicializar datos
	api.Post("/seed", handler.SeedListaPrecios)
}

func init() {
	registry.RegisterModule(RegisterRoutesListaPrecios)
}
