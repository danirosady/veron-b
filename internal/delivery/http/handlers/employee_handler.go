package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/dto/request"
	dtoResponse "github.com/tms/tyre/internal/dto/response"
	"github.com/tms/tyre/internal/usecase"
)

type EmployeeHandler struct {
	employeeUseCase *usecase.EmployeeUseCase
}

func NewEmployeeHandler(employeeUseCase *usecase.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{employeeUseCase: employeeUseCase}
}

// List returns a paginated list of employees
// GET /api/v1/employees?company_id=X&role=X&status=X
func (h *EmployeeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	role := c.Query("role")
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	employees, total, err := h.employeeUseCase.List(c.Request.Context(), page, perPage, uint(companyID), role, status)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data employee")
		return
	}

	response.SuccessWithPagination(c, "Success", dtoResponse.ToEmployeeResponses(employees), response.NewPagination(page, perPage, total))
}

// GetByID returns a single employee by ID
// GET /api/v1/employees/:id
func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	employee, err := h.employeeUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			response.NotFound(c, "Employee tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data employee")
		return
	}

	response.Success(c, http.StatusOK, "Success", dtoResponse.ToEmployeeResponse(employee))
}

// Create creates a new employee
// POST /api/v1/employees
func (h *EmployeeHandler) Create(c *gin.Context) {
	var req request.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	employee, err := h.employeeUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.BadRequest(c, "Company tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrEmployeeNameRequired):
			response.BadRequest(c, "Nama employee wajib diisi", nil)
		case errors.Is(err, usecase.ErrEmployeeNameTooLong):
			response.BadRequest(c, "Nama employee maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrEmployeeIDRequired):
			response.BadRequest(c, "Employee ID wajib diisi", nil)
		case errors.Is(err, usecase.ErrEmployeeIDTooLong):
			response.BadRequest(c, "Employee ID maksimal 50 karakter", nil)
		case errors.Is(err, usecase.ErrEmployeeIDExists):
			response.Error(c, http.StatusConflict, "Employee ID sudah digunakan untuk company dan role ini", nil)
		case errors.Is(err, usecase.ErrEmployeeRoleRequired):
			response.BadRequest(c, "Role wajib diisi", nil)
		case errors.Is(err, usecase.ErrEmployeeRoleInvalid):
			response.BadRequest(c, "Role harus salah satu dari: driver, supervisor, tyreman", nil)
		default:
			response.InternalError(c, "Gagal membuat employee")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Employee berhasil dibuat", dtoResponse.ToEmployeeResponse(employee))
}

// Update updates an existing employee
// PUT /api/v1/employees/:id
func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	employee, err := h.employeeUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmployeeNotFound):
			response.NotFound(c, "Employee tidak ditemukan")
		case errors.Is(err, usecase.ErrEmployeeNameRequired):
			response.BadRequest(c, "Nama employee wajib diisi", nil)
		case errors.Is(err, usecase.ErrEmployeeNameTooLong):
			response.BadRequest(c, "Nama employee maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrEmployeeIDRequired):
			response.BadRequest(c, "Employee ID wajib diisi", nil)
		case errors.Is(err, usecase.ErrEmployeeIDTooLong):
			response.BadRequest(c, "Employee ID maksimal 50 karakter", nil)
		case errors.Is(err, usecase.ErrEmployeeIDExists):
			response.Error(c, http.StatusConflict, "Employee ID sudah digunakan untuk company dan role ini", nil)
		default:
			response.InternalError(c, "Gagal mengupdate employee")
		}
		return
	}

	response.Success(c, http.StatusOK, "Employee berhasil diupdate", dtoResponse.ToEmployeeResponse(employee))
}

// Delete removes an employee
// DELETE /api/v1/employees/:id
func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.employeeUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			response.NotFound(c, "Employee tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal menghapus employee")
		return
	}

	response.Success(c, http.StatusOK, "Employee berhasil dihapus", nil)
}

// RegisterRoutes registers employee management routes on the given router group
func (h *EmployeeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	employees := rg.Group("/employees")
	{
		employees.GET("", h.List)
		employees.GET("/:id", h.GetByID)
		employees.POST("", h.Create)
		employees.PUT("/:id", h.Update)
		employees.DELETE("/:id", h.Delete)
	}
}
