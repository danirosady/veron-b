package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type APIResponse struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	Data      interface{}    `json:"data,omitempty"`
	Errors    []FieldError   `json:"errors,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, errors []FieldError) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func Paginated(c *gin.Context, data interface{}, meta *PaginationMeta) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   "OK",
		Data:      data,
		Pagination: meta,
	})
}

func BadRequest(c *gin.Context, message string, errors []FieldError) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, nil)
}

// ValidationError returns a 422 response with field-level errors from binding validation
func ValidationError(c *gin.Context, err error) {
	var fieldErrors []FieldError

	// Map of validation tags to user-friendly Indonesian messages
	tagMessages := map[string]string{
		"required":          "Field ini wajib diisi",
		"email":            "Format email tidak valid",
		"min":              "Nilai terlalu pendek",
		"max":              "Nilai terlalu panjang",
		"eqfield":          "Nilai tidak cocok",
		"password_strong":   "Kata sandi harus minimal 8 karakter dengan huruf besar, huruf kecil, angka, dan karakter khusus",
	}

	// Try direct type assertion first (validator.ValidationErrors is a []FieldError)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			// fe.Namespace() returns e.g. "CreateReplacementRequest.driver_id"
			// Extract just the JSON field name (last segment after the dot)
			fieldName := fe.Field()
			parts := strings.Split(fe.Namespace(), ".")
			if len(parts) >= 2 {
				fieldName = parts[len(parts)-1]
			}

			// Get user-friendly message based on tag
			msg, ok := tagMessages[fe.Tag()]
			if !ok {
				msg = fe.Tag()
			}

			// For min/max, add the actual limit
			if fe.Tag() == "min" || fe.Tag() == "max" {
				msg = fmt.Sprintf("%s (minimum %s)", msg, fe.Param())
			}

			fieldErrors = append(fieldErrors, FieldError{
				Field:   fieldName,
				Message: msg,
			})
		}
	}

	// If still empty, log error type for debugging
	if len(fieldErrors) == 0 {
		fmt.Printf("[DEBUG] ValidationError type: %T, msg: %s\n", err, err.Error())
	}

	Error(c, http.StatusUnprocessableEntity, "Mohon perbaiki kesalahan pada form", fieldErrors)
}

// SuccessWithPagination returns a success response with pagination metadata
func SuccessWithPagination(c *gin.Context, message string, data interface{}, pagination *PaginationMeta) {
	c.JSON(http.StatusOK, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// NewPagination builds a PaginationMeta from page, perPage, and total.
func NewPagination(page, perPage int, total int64) *PaginationMeta {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	totalPages := int(total / int64(perPage))
	if total%int64(perPage) > 0 {
		totalPages++
	}
	return &PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
