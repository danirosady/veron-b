package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	if ve, ok := err.(interface{ GetErrors() []interface{} }); ok {
		for _, e := range ve.GetErrors() {
			if fe, ok := e.(interface{ GetField() string; GetMessage() string }); ok {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   fe.GetField(),
					Message: fe.GetMessage(),
				})
			}
		}
	}
	Error(c, http.StatusUnprocessableEntity, "Validation failed", fieldErrors)
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
