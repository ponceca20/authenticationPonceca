package utils

import "github.com/gofiber/fiber/v2"

// SuccessResponse defines the structure for a successful API response.
type SuccessResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse defines the structure for an error API response.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SendSuccess sends a standardized success response.
func SendSuccess(c *fiber.Ctx, statusCode int, data interface{}, message ...string) error {
	res := SuccessResponse{
		Status: "success",
		Data:   data,
	}
	if len(message) > 0 {
		res.Message = message[0]
	}
	return c.Status(statusCode).JSON(res)
}

// SendError sends a standardized error response.
func SendError(c *fiber.Ctx, statusCode int, message string, err ...error) error {
	res := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	if len(err) > 0 && err[0] != nil {
		res.Error = err[0].Error()
	}
	return c.Status(statusCode).JSON(res)
}
