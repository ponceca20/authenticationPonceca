package product

import (
	"errors"
	"fmt"
	"practicev2/module/middleware"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseID extrae y valida el parámetro ID desde la solicitud
func ParseID(c *fiber.Ctx) (uint64, error) {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
	}
	return uint64(id), nil
}

// RespondWithError genera una respuesta de error estandarizada
func RespondWithError(c *fiber.Ctx, status int, err error, message string) error {
	// Log para depuración
	fmt.Printf("RespondWithError called with: c=%v, status=%d, err=%v, message=%s\n", c, status, err, message)
	if c == nil {
		fmt.Println("fiber.Ctx is nil!")
		return errors.New("fiber.Ctx is nil")
	}
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	} else {
		errMsg = ""
	}
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"error":   errMsg,
	})
}

// MapErrorStatus mapea errores a códigos de estado HTTP apropiados
func MapErrorStatus(err error) int {
	if err == nil {
		return fiber.StatusOK
	}
	lower := strings.ToLower(err.Error())
	switch {
	case strings.Contains(lower, "no encontrado"):
		return fiber.StatusNotFound
	case strings.Contains(lower, "validación") || strings.Contains(lower, "inválido"):
		return fiber.StatusBadRequest
	case strings.Contains(lower, "duplicado"):
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}
}

func getAuthenticatedUser(c *fiber.Ctx) (uint64, error) {
	uid, ok := c.Locals("usuario_id").(uint64)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Usuario no autenticado")
	}
	return uid, nil
}

// GetCompanyID recupera el ID de la primera empresa del usuario desde el contexto
func GetCompanyID(c *fiber.Ctx) (uint64, error) {
	companies, ok := c.Locals("companies").([]middleware.MinimalCompany)
	if !ok || len(companies) == 0 {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Usuario sin empresa asignada")
	}
	return companies[0].EmpresaID, nil
}
