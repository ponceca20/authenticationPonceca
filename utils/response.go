package utils

import "github.com/gofiber/fiber/v2"

// SuccessResponse defines the structure for a successful API response.
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse defines the structure for an error API response.
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// SendSuccessResponse sends a standardized success response.
func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse sends a standardized error response.
func SendErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := ErrorResponse{
		Status:  "error",
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}
