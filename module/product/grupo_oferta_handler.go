package product

import (
	"time"

	"practicev2/registry"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// GrupoOfertaHandler maneja las solicitudes HTTP relacionadas con GrupoOferta
type GrupoOfertaHandler struct {
	Service   *GrupoOfertaService
	validator *validator.Validate
}

func NewGrupoOfertaHandler(s *GrupoOfertaService) *GrupoOfertaHandler {
	return &GrupoOfertaHandler{
		Service:   s,
		validator: validator.New(),
	}
}

func (h *GrupoOfertaHandler) GetAllGrupoOfertas(c *fiber.Ctx) error {
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	grupos, err := h.Service.GetAllGrupoOfertas(empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Error al recuperar grupos de oferta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      grupos,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) GetGrupoOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	grupo, err := h.Service.GetGrupoOferta(id, empresaID)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "Grupo de oferta no encontrado")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      grupo,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) CreateGrupoOferta(c *fiber.Ctx) error {
	var grupo GrupoOferta
	if err := c.BodyParser(&grupo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&grupo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	if err := h.Service.CreateGrupoOferta(&grupo); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo crear el grupo de oferta")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":    "success",
		"data":      grupo,
		"message":   "Grupo de oferta creado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) UpdateGrupoOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	var grupo GrupoOferta
	if err := c.BodyParser(&grupo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.validator.Struct(&grupo); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Fallo en la validación")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	updated, err := h.Service.UpdateGrupoOferta(id, empresaID, &grupo)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar el grupo de oferta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      updated,
		"message":   "Grupo de oferta actualizado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) DeleteGrupoOferta(c *fiber.Ctx) error {
	id, err := ParseID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "ID inválido")
	}
	empresaID, err := GetCompanyID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, err, "Error de autenticación")
	}
	if err := h.Service.DeleteGrupoOferta(id, empresaID); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo eliminar el grupo de oferta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Grupo de oferta eliminado exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) SeedGrupoOfertas(c *fiber.Ctx) error {
	var grupos []GrupoOferta
	if err := c.BodyParser(&grupos); err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, err, "Cuerpo de la solicitud inválido")
	}
	if err := h.Service.SeedGrupoOfertas(grupos); err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudieron inicializar los grupos de oferta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Grupos de oferta inicializados exitosamente",
		"count":     len(grupos),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *GrupoOfertaHandler) PatchGrupoOferta(c *fiber.Ctx) error {
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
	grupo, err := h.Service.PatchGrupoOferta(id, empresaID, fields)
	if err != nil {
		return RespondWithError(c, MapErrorStatus(err), err, "No se pudo actualizar parcialmente el grupo de oferta")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"data":      grupo,
		"message":   "Grupo de oferta actualizado parcialmente con éxito",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterRoutesGrupoOferta registra las rutas para GrupoOferta
func RegisterRoutesGrupoOferta(app *fiber.App) {
	repo := NewGrupoOfertaRepository()
	service := NewGrupoOfertaService(repo)
	handler := NewGrupoOfertaHandler(service)
	api := app.Group("/api/v1/grupos-oferta")
	api.Get("/", handler.GetAllGrupoOfertas)
	api.Get("/:id", handler.GetGrupoOferta)
	api.Post("/", handler.CreateGrupoOferta)
	api.Put("/:id", handler.UpdateGrupoOferta)
	api.Patch("/:id", handler.PatchGrupoOferta)
	api.Delete("/:id", handler.DeleteGrupoOferta)
	api.Post("/seed", handler.SeedGrupoOfertas)
}

func init() {
	registry.RegisterModule(RegisterRoutesGrupoOferta)
}
