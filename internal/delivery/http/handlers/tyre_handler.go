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

type TyreHandler struct {
	tyreUseCase *usecase.TyreUseCase
}

func NewTyreHandler(tyreUseCase *usecase.TyreUseCase) *TyreHandler {
	return &TyreHandler{tyreUseCase: tyreUseCase}
}

// List returns a paginated list of tyres
// GET /api/v1/tyres?company_id=X
func (h *TyreHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")
	brandID := c.Query("brand_id")
	sizeID := c.Query("size_id")
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	tyres, total, err := h.tyreUseCase.List(c.Request.Context(), page, perPage, uint(companyID), status, brandID, sizeID)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data tyre")
		return
	}

	response.SuccessWithPagination(c, "Success", tyres, response.NewPagination(page, perPage, total))
}

// GetByID returns a single tyre by ID
// GET /api/v1/tyres/:id
func (h *TyreHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	tyre, err := h.tyreUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrTyreNotFound) {
			response.NotFound(c, "Tyre tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data tyre")
		return
	}

	response.Success(c, http.StatusOK, "Success", tyre)
}

// GetByBarcode returns a tyre by barcode
// GET /api/v1/tyres/barcode/:barcode
func (h *TyreHandler) GetByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")
	if barcode == "" {
		response.BadRequest(c, "Barcode wajib diisi", nil)
		return
	}

	tyre, err := h.tyreUseCase.GetByBarcode(c.Request.Context(), barcode)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrTyreBarcodeRequired):
			response.BadRequest(c, "Barcode wajib diisi", nil)
		case errors.Is(err, usecase.ErrTyreNotFound):
			response.NotFound(c, "Tyre tidak ditemukan")
		default:
			response.InternalError(c, "Gagal mengambil data tyre")
		}
		return
	}

	response.Success(c, http.StatusOK, "Success", tyre)
}

// GetSpareTyres returns all spare (and dismounted) tyres for a company
// GET /api/v1/tyres/spare?company_id=X
func (h *TyreHandler) GetSpareTyres(c *gin.Context) {
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)
	if companyID == 0 {
		response.BadRequest(c, "Company ID wajib diisi", nil)
		return
	}

	tyres, err := h.tyreUseCase.GetSpareTyres(c.Request.Context(), uint(companyID))
	if err != nil {
		response.InternalError(c, "Gagal mengambil data spare tyre")
		return
	}

	response.Success(c, http.StatusOK, "Success", tyres)
}

// Create creates a new tyre master record
// POST /api/v1/tyres?company_id=X  (or company_id in body)
func (h *TyreHandler) Create(c *gin.Context) {
	// Support company_id from query param or body
	companyIDParam := c.DefaultQuery("company_id", "")
	var companyID uint64
	if companyIDParam != "" {
		companyID, _ = strconv.ParseUint(companyIDParam, 10, 32)
	}
	if companyID == 0 {
		// Try reading from body
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err == nil {
			if cid, ok := body["company_id"].(float64); ok {
				companyID = uint64(cid)
			}
		}
	}
	if companyID == 0 {
		response.BadRequest(c, "Company ID wajib diisi", nil)
		return
	}

	var req request.CreateTyreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	tyre, err := h.tyreUseCase.Create(c.Request.Context(), uint(companyID), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.BadRequest(c, "Company tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrTyreSizeNotFound):
			response.BadRequest(c, "Tyre size tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrTyreBrandNotFound):
			response.BadRequest(c, "Tyre brand tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrTyrePatternNotFound):
			response.BadRequest(c, "Tyre pattern tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrTyreBarcodeRequired),
			errors.Is(err, usecase.ErrTyreSerialRequired):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, usecase.ErrTyreBarcodeTooLong),
			errors.Is(err, usecase.ErrTyreSerialTooLong):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, usecase.ErrTyreBarcodeExists),
			errors.Is(err, usecase.ErrTyreSerialExists):
			response.Error(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, usecase.ErrTyreRTDExceedsOTD),
			errors.Is(err, usecase.ErrTyreRTDInvalid):
			response.BadRequest(c, err.Error(), nil)
		case errors.Is(err, usecase.ErrTyrePSIOutOfRange):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalError(c, "Gagal membuat tyre")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Tyre berhasil dibuat", tyre)
}

// Update updates an existing tyre master record
// PUT /api/v1/tyres/:id
func (h *TyreHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateTyreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	tyre, err := h.tyreUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrTyreNotFound):
			response.NotFound(c, "Tyre tidak ditemukan")
		case errors.Is(err, usecase.ErrTyreRTDExceedsOTD),
			errors.Is(err, usecase.ErrTyreRTDInvalid),
			errors.Is(err, usecase.ErrTyrePSIOutOfRange):
			response.BadRequest(c, err.Error(), nil)
		default:
			response.InternalError(c, "Gagal mengupdate tyre")
		}
		return
	}

	response.Success(c, http.StatusOK, "Tyre berhasil diupdate", tyre)
}

// Delete removes a tyre
// DELETE /api/v1/tyres/:id
func (h *TyreHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.tyreUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrTyreNotFound) {
			response.NotFound(c, "Tyre tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal menghapus tyre")
		return
	}

	response.Success(c, http.StatusOK, "Tyre berhasil dihapus", nil)
}

// RegisterRoutes registers tyre management routes on the given router group
func (h *TyreHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tyres := rg.Group("/tyres")
	{
		tyres.GET("", h.List)
		tyres.POST("", h.Create)
		tyres.GET("/spare", h.GetSpareTyres)
		tyres.GET("/barcode/:barcode", h.GetByBarcode)
		tyres.GET("/:id", h.GetByID)
		tyres.PUT("/:id", h.Update)
		tyres.DELETE("/:id", h.Delete)
	}
}
