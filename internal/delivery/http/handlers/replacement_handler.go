package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/middleware"
	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/dto/request"
	respDTO "github.com/tms/tyre/internal/dto/response"
	"github.com/tms/tyre/internal/usecase"
)

type ReplacementHandler struct {
	replacementUseCase *usecase.ReplacementUseCase
}

func NewReplacementHandler(replacementUseCase *usecase.ReplacementUseCase) *ReplacementHandler {
	return &ReplacementHandler{replacementUseCase: replacementUseCase}
}

// Create creates a new replacement record
// POST /api/v1/replacements
func (h *ReplacementHandler) Create(c *gin.Context) {
	operatorID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "Operator not found")
		return
	}

	var req request.CreateReplacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	replacement, err := h.replacementUseCase.Create(c.Request.Context(), &req, operatorID.(uint))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrReplacementUnitNotFound):
			response.BadRequest(c, "Unit tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrReplacementDriverNotFound):
			response.BadRequest(c, "Driver tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrReplacementTyreNotFound):
			response.BadRequest(c, "Tyre tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrReplacementTyreNotSpare):
			response.BadRequest(c, "Tyre baru harus berstatus spare untuk mount", nil)
		case errors.Is(err, usecase.ErrReplacementTyreNotMounted):
			response.BadRequest(c, "Tyre lama harus berstatus mounted untuk dismount", nil)
		case errors.Is(err, usecase.ErrReplacementTyreRequired):
			response.BadRequest(c, "New tyre ID wajib diisi untuk mount action", nil)
		case errors.Is(err, usecase.ErrReplacementSwapBothRequired):
			response.BadRequest(c, "Swap action requires both old_tyre_id and new_tyre_id", nil)
		default:
			if strings.Contains(err.Error(), "invalid date format") {
				response.BadRequest(c, err.Error(), nil)
				return
			}
			response.InternalError(c, "Gagal membuat replacement")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Replacement berhasil dibuat", map[string]interface{}{
		"id": replacement.ID,
		"unit_id": replacement.UnitID,
		"date": replacement.Date,
	})
}

// List returns a paginated list of replacements
// GET /api/v1/replacements
func (h *ReplacementHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	companyID, _ := strconv.ParseUint(c.DefaultQuery("company_id", "0"), 10, 32)
	projectID, _ := strconv.ParseUint(c.DefaultQuery("project_id", "0"), 10, 32)
	unitID, _ := strconv.ParseUint(c.DefaultQuery("unit_id", "0"), 10, 32)
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	role := ""
	tenantID := uint(0)
	if claims, ok := middleware.GetClaims(c); ok {
		role = claims.Role
		tenantID = claims.TenantID
	}

	filters := map[string]interface{}{}
	if companyID > 0 {
		filters["company_id"] = uint(companyID)
	} else if role != "superadmin" && tenantID > 0 {
		filters["company_id"] = tenantID
	}
	// superadmin sees all — no company_id filter needed
	if projectID > 0 {
		filters["project_id"] = uint(projectID)
	}
	if unitID > 0 {
		filters["unit_id"] = uint(unitID)
	}
	if dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo != "" {
		filters["date_to"] = dateTo
	}

	replacements, total, err := h.replacementUseCase.List(c.Request.Context(), page, perPage, filters)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data replacement")
		return
	}

	response.SuccessWithPagination(c, "Success", respDTO.ToReplacementResponses(replacements), response.NewPagination(page, perPage, total))
}

// GetByID returns a single replacement by ID
// GET /api/v1/replacements/:id
func (h *ReplacementHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	replacement, err := h.replacementUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrReplacementNotFound) {
			response.NotFound(c, "Replacement tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data replacement")
		return
	}

	response.Success(c, http.StatusOK, "Success", respDTO.ToReplacementResponse(replacement))
}

// Update updates an existing replacement record
// PUT /api/v1/replacements/:id
func (h *ReplacementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	operatorID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "Operator not found")
		return
	}

	var req request.UpdateReplacementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	replacement, err := h.replacementUseCase.Update(c.Request.Context(), uint(id), &req, operatorID.(uint))
	if err != nil {
		if errors.Is(err, usecase.ErrReplacementNotFound) {
			response.NotFound(c, "Replacement tidak ditemukan")
			return
		}
		if strings.Contains(err.Error(), "invalid date format") {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalError(c, "Gagal mengupdate replacement")
		return
	}

	response.Success(c, http.StatusOK, "Replacement berhasil diupdate", replacement)
}

// Delete removes a replacement record
// DELETE /api/v1/replacements/:id
func (h *ReplacementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.replacementUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrReplacementNotFound) {
			response.NotFound(c, "Replacement tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal menghapus replacement")
		return
	}

	response.Success(c, http.StatusOK, "Replacement berhasil dihapus", nil)
}

// GetLastByUnit returns the most recent replacement for a unit
// GET /api/v1/replacements/unit/:id/last
func (h *ReplacementHandler) GetLastByUnit(c *gin.Context) {
	unitID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID unit tidak valid", nil)
		return
	}

	replacement, err := h.replacementUseCase.GetLastByUnitID(uint(unitID))
	if err != nil {
		response.InternalError(c, "Gagal mengambil data replacement")
		return
	}

	if replacement == nil {
		response.Success(c, http.StatusOK, "Success", nil)
		return
	}
	response.Success(c, http.StatusOK, "Success", respDTO.ToReplacementResponse(replacement))
}

// RegisterRoutes registers replacement management routes on the given router group
func (h *ReplacementHandler) RegisterRoutes(rg *gin.RouterGroup) {
	replacements := rg.Group("/replacements")
	{
		replacements.GET("", h.List)
		replacements.GET("/:id", h.GetByID)
		replacements.GET("/unit/:id/last", h.GetLastByUnit)
		replacements.POST("", h.Create)
		replacements.PUT("/:id", h.Update)
		replacements.DELETE("/:id", h.Delete)
	}
}
