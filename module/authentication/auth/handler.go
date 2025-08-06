package auth

import (
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles the HTTP requests for authentication.
type AuthHandler struct {
	service AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register is the handler for the user registration endpoint.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var dto RegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errors := utils.ValidateStruct(&dto); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	identity, err := h.service.Register(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusConflict, err.Error())
	}

	// Return a simplified, public version of the identity
	return utils.SendSuccess(c, fiber.StatusCreated, ToIdentityProfileDTO(identity), "User registered successfully")
}

// Login is the handler for the user login endpoint.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var dto LoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errors := utils.ValidateStruct(&dto); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	tokenResponse, err := h.service.Login(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, tokenResponse, "Login successful")
}

// RefreshToken is the handler for refreshing an access token.
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var dto RefreshTokenDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errors := utils.ValidateStruct(&dto); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	tokenResponse, err := h.service.RefreshToken(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Failed to refresh token", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, tokenResponse, "Token refreshed successfully")
}

// Logout is the handler for logging out a user.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var dto RefreshTokenDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := h.service.Logout(&dto); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to logout", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Logout successful")
}

// ForgotPassword is the handler for initiating a password reset.
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var dto ForgotPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	err := h.service.ForgotPassword(&dto)
	if err != nil {
		// Do not reveal if the email was found or not for security reasons.
		// Log the internal error.
	}

	// Always return a success message to prevent email enumeration attacks.
	return utils.SendSuccess(c, fiber.StatusOK, nil, "If an account with that email exists, a password reset link has been sent.")
}

// ResetPassword is the handler for resetting a password with a token.
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var dto ResetPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.ResetPassword(&dto); err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid token or failed to reset password", err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Password has been reset successfully.")
}

// ChangePassword is the handler for changing password (authenticated user).
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var dto ChangePasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	// Get identity ID from JWT token (set by middleware)
	identityID := c.Locals("identity_id")
	if identityID == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Authentication required")
	}

	if err := h.service.ChangePassword(identityID.(string), &dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Password changed successfully")
}

// VerifyEmail is the handler for verifying email address with token.
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var dto VerifyEmailDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.VerifyEmail(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Email verified successfully")
}

// ResendVerification is the handler for resending email verification.
func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	var dto ResendVerificationDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if errs := utils.ValidateStruct(&dto); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	if err := h.service.ResendVerification(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "If an unverified account with that email exists, a verification link has been sent")
}
