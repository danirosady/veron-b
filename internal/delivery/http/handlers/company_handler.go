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

type CompanyHandler struct {
	companyUseCase *usecase.CompanyUseCase
}

func NewCompanyHandler(companyUseCase *usecase.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{companyUseCase: companyUseCase}
}

// List returns a paginated list of companies
// GET /api/v1/companies
func (h *CompanyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	companies, total, err := h.companyUseCase.List(c.Request.Context(), page, perPage, status)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data company")
		return
	}

	response.SuccessWithPagination(c, "Success", companies, response.NewPagination(page, perPage, total))
}

// GetByID returns a single company by ID
// GET /api/v1/companies/:id
func (h *CompanyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	company, err := h.companyUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrCompanyNotFound) {
			response.NotFound(c, "Company tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data company")
		return
	}

	response.Success(c, http.StatusOK, "Success", company)
}

// Create creates a new company
// POST /api/v1/companies
func (h *CompanyHandler) Create(c *gin.Context) {
	var req request.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	company, err := h.companyUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNameRequired):
			response.BadRequest(c, "Nama company wajib diisi", nil)
		case errors.Is(err, usecase.ErrCompanyNameTooLong):
			response.BadRequest(c, "Nama company maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrCompanyNameExists):
			response.Error(c, http.StatusConflict, "Nama company sudah digunakan", nil)
		case errors.Is(err, usecase.ErrCompanyEmailInvalid):
			response.BadRequest(c, "Format email tidak valid", nil)
		default:
			response.InternalError(c, "Gagal membuat company")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Company berhasil dibuat", company)
}

// Update updates an existing company
// PUT /api/v1/companies/:id
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	company, err := h.companyUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.NotFound(c, "Company tidak ditemukan")
		case errors.Is(err, usecase.ErrCompanyNameRequired):
			response.BadRequest(c, "Nama company wajib diisi", nil)
		case errors.Is(err, usecase.ErrCompanyNameTooLong):
			response.BadRequest(c, "Nama company maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrCompanyNameExists):
			response.Error(c, http.StatusConflict, "Nama company sudah digunakan", nil)
		case errors.Is(err, usecase.ErrCompanyEmailInvalid):
			response.BadRequest(c, "Format email tidak valid", nil)
		default:
			response.InternalError(c, "Gagal mengupdate company")
		}
		return
	}

	response.Success(c, http.StatusOK, "Company berhasil diupdate", company)
}

// Delete removes a company
// DELETE /api/v1/companies/:id
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.companyUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.NotFound(c, "Company tidak ditemukan")
		case errors.Is(err, usecase.ErrCompanyHasActive):
			response.Error(c, http.StatusConflict, "Company masih memiliki data aktif (project, unit, atau driver)", nil)
		default:
			response.InternalError(c, "Gagal menghapus company")
		}
		return
	}

	response.Success(c, http.StatusOK, "Company berhasil dihapus", nil)
}

// RegisterRoutes registers company management routes on the given router group
func (h *CompanyHandler) RegisterRoutes(rg *gin.RouterGroup) {
	companies := rg.Group("/companies")
	{
		companies.GET("", h.List)
		companies.GET("/:id", h.GetByID)
		companies.POST("", h.Create)
		companies.PUT("/:id", h.Update)
		companies.DELETE("/:id", h.Delete)
	}
}