package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/dto/request"
	"github.com/tms/tyre/internal/usecase"
)

type DriverHandler struct {
	driverUseCase *usecase.DriverUseCase
}

func NewDriverHandler(driverUseCase *usecase.DriverUseCase) *DriverHandler {
	return &DriverHandler{driverUseCase: driverUseCase}
}

// List returns a paginated list of drivers
// GET /api/v1/drivers?company_id=X
func (h *DriverHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	drivers, total, err := h.driverUseCase.List(c.Request.Context(), page, perPage, uint(companyID), status)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data driver")
		return
	}

	response.SuccessWithPagination(c, "Success", drivers, response.NewPagination(page, perPage, total))
}

// GetByID returns a single driver by ID
// GET /api/v1/drivers/:id
func (h *DriverHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	driver, err := h.driverUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) {
			response.NotFound(c, "Driver tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data driver")
		return
	}

	response.Success(c, http.StatusOK, "Success", driver)
}

// Create creates a new driver
// POST /api/v1/drivers
func (h *DriverHandler) Create(c *gin.Context) {
	var req request.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	driver, err := h.driverUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.BadRequest(c, "Company tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrDriverNameRequired):
			response.BadRequest(c, "Nama driver wajib diisi", nil)
		case errors.Is(err, usecase.ErrDriverNameTooLong):
			response.BadRequest(c, "Nama driver maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeReq):
			response.BadRequest(c, "Employee ID wajib diisi", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeLong):
			response.BadRequest(c, "Employee ID maksimal 50 karakter", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeExist):
			response.Error(c, http.StatusConflict, "Employee ID sudah digunakan untuk company ini", nil)
		default:
			response.InternalError(c, "Gagal membuat driver")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Driver berhasil dibuat", driver)
}

// Update updates an existing driver
// PUT /api/v1/drivers/:id
func (h *DriverHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	driver, err := h.driverUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrDriverNotFound):
			response.NotFound(c, "Driver tidak ditemukan")
		case errors.Is(err, usecase.ErrDriverNameRequired):
			response.BadRequest(c, "Nama driver wajib diisi", nil)
		case errors.Is(err, usecase.ErrDriverNameTooLong):
			response.BadRequest(c, "Nama driver maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeReq):
			response.BadRequest(c, "Employee ID wajib diisi", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeLong):
			response.BadRequest(c, "Employee ID maksimal 50 karakter", nil)
		case errors.Is(err, usecase.ErrDriverEmployeeExist):
			response.Error(c, http.StatusConflict, "Employee ID sudah digunakan untuk company ini", nil)
		default:
			response.InternalError(c, "Gagal mengupdate driver")
		}
		return
	}

	response.Success(c, http.StatusOK, "Driver berhasil diupdate", driver)
}

// Delete removes a driver
// DELETE /api/v1/drivers/:id
func (h *DriverHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.driverUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) {
			response.NotFound(c, "Driver tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal menghapus driver")
		return
	}

	response.Success(c, http.StatusOK, "Driver berhasil dihapus", nil)
}

// RegisterRoutes registers driver management routes on the given router group
func (h *DriverHandler) RegisterRoutes(rg *gin.RouterGroup) {
	drivers := rg.Group("/drivers")
	{
		drivers.GET("", h.List)
		drivers.GET("/:id", h.GetByID)
		drivers.POST("", h.Create)
		drivers.PUT("/:id", h.Update)
		drivers.DELETE("/:id", h.Delete)
	}
}