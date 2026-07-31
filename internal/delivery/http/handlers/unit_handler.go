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

type UnitHandler struct {
	unitUseCase *usecase.UnitUseCase
}

func NewUnitHandler(unitUseCase *usecase.UnitUseCase) *UnitHandler {
	return &UnitHandler{unitUseCase: unitUseCase}
}

// List returns a paginated list of units
// GET /api/v1/units?company_id=X&project_id=Y
func (h *UnitHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)
	projectID, _ := strconv.ParseUint(c.DefaultQuery("project_id", "0"), 10, 32)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	units, total, err := h.unitUseCase.List(c.Request.Context(), page, perPage, uint(companyID), uint(projectID), status)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data unit")
		return
	}

	response.SuccessWithPagination(c, "Success", units, response.NewPagination(page, perPage, total))
}

// GetByID returns a single unit by ID
// GET /api/v1/units/:id
func (h *UnitHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	unit, err := h.unitUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrUnitNotFound) {
			response.NotFound(c, "Unit tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data unit")
		return
	}

	response.Success(c, http.StatusOK, "Success", unit)
}

// Create creates a new unit
// POST /api/v1/units
func (h *UnitHandler) Create(c *gin.Context) {
	var req request.CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	unit, err := h.unitUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.BadRequest(c, "Company tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrProjectNotFound):
			response.BadRequest(c, "Project tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrUnitIDRequired),
			errors.Is(err, usecase.ErrUnitModelRequired),
			errors.Is(err, usecase.ErrUnitTyreSizeRequired),
			errors.Is(err, usecase.ErrUnitTypeRequired):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, usecase.ErrUnitIDTooLong):
			response.BadRequest(c, "Unit ID maksimal 50 karakter", nil)
		case errors.Is(err, usecase.ErrUnitInvalidPosition):
			response.BadRequest(c, "Max position harus antara 1 dan 20", nil)
		case errors.Is(err, usecase.ErrUnitIDExists):
			response.Error(c, http.StatusConflict, "Unit ID sudah digunakan", nil)
		default:
			response.InternalError(c, "Gagal membuat unit")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Unit berhasil dibuat", unit)
}

// Update updates an existing unit
// PUT /api/v1/units/:id
func (h *UnitHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	unit, err := h.unitUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUnitNotFound):
			response.NotFound(c, "Unit tidak ditemukan")
		case errors.Is(err, usecase.ErrProjectNotFound):
			response.BadRequest(c, "Project tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrUnitModelRequired),
			errors.Is(err, usecase.ErrUnitTyreSizeRequired),
			errors.Is(err, usecase.ErrUnitTypeRequired):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, usecase.ErrUnitInvalidPosition):
			response.BadRequest(c, "Max position harus antara 1 dan 20", nil)
		default:
			response.InternalError(c, "Gagal mengupdate unit")
		}
		return
	}

	response.Success(c, http.StatusOK, "Unit berhasil diupdate", unit)
}

// Delete removes a unit
// DELETE /api/v1/units/:id
func (h *UnitHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.unitUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUnitNotFound):
			response.NotFound(c, "Unit tidak ditemukan")
		case errors.Is(err, usecase.ErrUnitHasMountedTyres):
			response.Error(c, http.StatusConflict, "Unit masih memiliki tyre yang terpasang", nil)
		case errors.Is(err, usecase.ErrUnitHasReplacements):
			response.Error(c, http.StatusConflict, "Unit memiliki riwayat replacement", nil)
		default:
			response.InternalError(c, "Gagal menghapus unit")
		}
		return
	}

	response.Success(c, http.StatusOK, "Unit berhasil dihapus", nil)
}

// UpdateHM updates the current HM of a unit
// PUT /api/v1/units/:id/hm
func (h *UnitHandler) UpdateHM(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateUnitHMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	unit, err := h.unitUseCase.UpdateHM(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUnitNotFound):
			response.NotFound(c, "Unit tidak ditemukan")
		case errors.Is(err, usecase.ErrUnitHMDecrease):
			response.BadRequest(c, "Current HM tidak boleh lebih kecil dari nilai saat ini", nil)
		default:
			response.InternalError(c, "Gagal mengupdate HM")
		}
		return
	}

	response.Success(c, http.StatusOK, "HM berhasil diupdate", unit)
}

// GetTyres returns tyres mounted on a unit
// GET /api/v1/units/:id/tyres
func (h *UnitHandler) GetTyres(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	tyres, err := h.unitUseCase.GetTyres(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrUnitNotFound) {
			response.NotFound(c, "Unit tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data tyre")
		return
	}

	response.Success(c, http.StatusOK, "Success", tyres)
}

// RegisterRoutes registers unit management routes on the given router group
func (h *UnitHandler) RegisterRoutes(rg *gin.RouterGroup) {
	units := rg.Group("/units")
	{
		units.GET("", h.List)
		units.GET("/:id", h.GetByID)
		units.POST("", h.Create)
		units.PUT("/:id", h.Update)
		units.DELETE("/:id", h.Delete)
		units.PUT("/:id/hm", h.UpdateHM)
		units.GET("/:id/tyres", h.GetTyres)
	}
}