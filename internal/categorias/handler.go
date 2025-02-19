package categorias

import (
	"errors"
	"practicev2/registry"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	Service *CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{Service: NewCategoryService()}
}
func respondWithError(c *fiber.Ctx, status int, _ error, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.Service.GetAllCategories()
	if err != nil {
		return respondWithError(c, fiber.StatusInternalServerError, err, "Error")
	}
	return c.JSON(categories)
}
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return respondWithError(c, fiber.StatusBadRequest, errors.New("invalid id"), "ID inválido")
	}
	category, err := h.Service.GetCategory(id)
	if err != nil {
		return respondWithError(c, fiber.StatusNotFound, err, "No encontrado")
	}
	return c.JSON(category)
}
func (h *CategoryHandler) NewCategory(c *fiber.Ctx) error {
	var category Categoria
	if err := c.BodyParser(&category); err != nil {
		return respondWithError(c, fiber.StatusBadRequest, err, "Datos inválidos")
	}
	if err := h.Service.NewCategory(&category); err != nil {
		return respondWithError(c, fiber.StatusInternalServerError, err, "Error")
	}
	return c.Status(fiber.StatusCreated).JSON(category)
}
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var data Categoria
	if err := c.BodyParser(&data); err != nil {
		return respondWithError(c, fiber.StatusBadRequest, err, "Datos inválidos")
	}
	category, err := h.Service.UpdateCategory(id, &data)
	if err != nil {
		return respondWithError(c, fiber.StatusInternalServerError, err, "Error")
	}
	return c.JSON(category)
}
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.Service.DeleteCategory(id); err != nil {
		return respondWithError(c, fiber.StatusInternalServerError, err, "Error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *CategoryHandler) SeedCategories(c *fiber.Ctx) error {
	categories := []Categoria{
		{Name: "Ficción"},
		{Name: "No ficción"},
		{Name: "Ciencia ficción"},
		{Name: "Fantasía"},
		{Name: "Terror"},
		{Name: "Romance"},
	}
	if err := h.Service.SeedCategories(categories); err != nil {
		return respondWithError(c, fiber.StatusInternalServerError, err, "Error")
	}
	//return c.SendStatus(fiber.StatusNoContent)
	return c.JSON(categories)
}
func RegisterRoutes2(app *fiber.App) {
	h := NewCategoryHandler()
	api := app.Group("/api/categorias")
	api.Get("/", h.GetAllCategories)
	// Register static routes before dynamic ones
	api.Get("/seed", h.SeedCategories)
	api.Get("/:id", h.GetCategory)
	api.Post("/", h.NewCategory)
	api.Put("/:id", h.UpdateCategory)
	api.Delete("/:id", h.DeleteCategory)
}

// Added init function to register this module automatically.
func init() {
	registry.RegisterModule(RegisterRoutes2)
}
