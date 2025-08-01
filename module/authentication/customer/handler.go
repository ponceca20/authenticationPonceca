package customer

import (
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// CustomerHandler handles the HTTP requests for e-commerce customers.
type CustomerHandler struct {
	service CustomerService
}

// NewCustomerHandler creates a new instance of CustomerHandler.
func NewCustomerHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// Register is the handler for creating a new customer account.
func (h *CustomerHandler) Register(c *fiber.Ctx) error {
	var dto CustomerRegistrationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// El servicio maneja validación más detallada.
	customer, err := h.service.RegisterCustomer(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusConflict, err.Error())
	}

	// Construir respuesta completa con información del cliente
	customerResponse := fiber.Map{
		"id":                customer.ID,
		"customer_number":   customer.CustomerNumber,
		"email":             dto.Email, // Usar el email del DTO ya que está validado
		"accepts_marketing": customer.AcceptsMarketing,
	}

	// Devolver información del cliente registrado
	return utils.SendSuccess(c, fiber.StatusCreated, fiber.Map{
		"customer": customerResponse,
	}, "Cliente registrado exitosamente")
}

// GetProfile is the handler for retrieving the authenticated customer's profile.
// This route must be protected by the SmartAuthMiddleware.
func (h *CustomerHandler) GetProfile(c *fiber.Ctx) error {
	// The SmartAuthMiddleware should have placed the AuthContext in locals.
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	profileDTO, err := h.service.GetCustomerProfile(authCtx.Identity.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Customer profile not found", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, profileDTO)
}

// AddAddress is the handler for adding a new shipping address.
func (h *CustomerHandler) AddAddress(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	var dto AddressDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	address, err := h.service.AddAddress(authCtx.Identity.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to add address", err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, address)
}

// ListAddresses is the handler for listing a customer's shipping addresses.
func (h *CustomerHandler) ListAddresses(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	addresses, err := h.service.ListAddresses(authCtx.Identity.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve addresses", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, addresses)
}

// UpdateAddress is the handler for updating a shipping address.
func (h *CustomerHandler) UpdateAddress(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	addressID := c.Params("id")
	if addressID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Address ID is required")
	}

	var dto AddressDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	address, err := h.service.UpdateAddress(authCtx.Identity.ID, addressID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update address", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, address)
}

// DeleteAddress is the handler for deleting a shipping address.
func (h *CustomerHandler) DeleteAddress(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	addressID := c.Params("id")
	if addressID == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Address ID is required")
	}

	if err := h.service.DeleteAddress(authCtx.Identity.ID, addressID); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete address", err)
	}

	return utils.SendSuccess(c, fiber.StatusNoContent, nil)
}

// GetPreferences is the handler for retrieving a customer's preferences.
func (h *CustomerHandler) GetPreferences(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	prefs, err := h.service.GetPreferences(authCtx.Identity.ID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve preferences", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, prefs)
}

// DeleteProfile is a placeholder for deleting a customer's profile.
func (h *CustomerHandler) DeleteProfile(c *fiber.Ctx) error {
	return utils.SendError(c, fiber.StatusNotImplemented, "Delete profile not implemented yet")
}

// SetDefaultAddress is a placeholder for setting a default shipping address.
func (h *CustomerHandler) SetDefaultAddress(c *fiber.Ctx) error {
	return utils.SendError(c, fiber.StatusNotImplemented, "Set default address not implemented yet")
}

// UpdatePreferences is the handler for updating a customer's preferences.
func (h *CustomerHandler) UpdatePreferences(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("authContext").(*middleware.AuthContext)
	if !ok || authCtx.IsGuest {
		return utils.SendError(c, fiber.StatusUnauthorized, "User not authenticated")
	}

	var dto PreferencesDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	prefs, err := h.service.UpdatePreferences(authCtx.Identity.ID, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to update preferences", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, prefs)
}
