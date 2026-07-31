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

type ProjectHandler struct {
	projectUseCase *usecase.ProjectUseCase
}

func NewProjectHandler(projectUseCase *usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{projectUseCase: projectUseCase}
}

// List returns a paginated list of projects
// GET /api/v1/projects?company_id=X
func (h *ProjectHandler) List(c *gin.Context) {
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

	projects, total, err := h.projectUseCase.List(c.Request.Context(), page, perPage, uint(companyID), status)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data project")
		return
	}

	response.SuccessWithPagination(c, "Success", projects, response.NewPagination(page, perPage, total))
}

// GetByID returns a single project by ID
// GET /api/v1/projects/:id
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	project, err := h.projectUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrProjectNotFound) {
			response.NotFound(c, "Project tidak ditemukan")
			return
		}
		response.InternalError(c, "Gagal mengambil data project")
		return
	}

	response.Success(c, http.StatusOK, "Success", project)
}

// Create creates a new project
// POST /api/v1/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	project, err := h.projectUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrCompanyNotFound):
			response.BadRequest(c, "Company tidak ditemukan", nil)
		case errors.Is(err, usecase.ErrProjectNameRequired):
			response.BadRequest(c, "Nama project wajib diisi", nil)
		case errors.Is(err, usecase.ErrProjectNameTooLong):
			response.BadRequest(c, "Nama project maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrProjectInvalidDates):
			response.BadRequest(c, "Start date harus lebih kecil atau sama dengan end date", nil)
		default:
			response.InternalError(c, "Gagal membuat project")
		}
		return
	}

	response.Success(c, http.StatusCreated, "Project berhasil dibuat", project)
}

// Update updates an existing project
// PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	project, err := h.projectUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrProjectNotFound):
			response.NotFound(c, "Project tidak ditemukan")
		case errors.Is(err, usecase.ErrProjectNameRequired):
			response.BadRequest(c, "Nama project wajib diisi", nil)
		case errors.Is(err, usecase.ErrProjectNameTooLong):
			response.BadRequest(c, "Nama project maksimal 255 karakter", nil)
		case errors.Is(err, usecase.ErrProjectInvalidDates):
			response.BadRequest(c, "Start date harus lebih kecil atau sama dengan end date", nil)
		default:
			response.InternalError(c, "Gagal mengupdate project")
		}
		return
	}

	response.Success(c, http.StatusOK, "Project berhasil diupdate", project)
}

// Delete removes a project
// DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	err = h.projectUseCase.Delete(c.Request.Context(), uint(id))
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrProjectNotFound):
			response.NotFound(c, "Project tidak ditemukan")
		case errors.Is(err, usecase.ErrProjectHasUnits):
			response.Error(c, http.StatusConflict, "Project masih memiliki unit aktif", nil)
		default:
			response.InternalError(c, "Gagal menghapus project")
		}
		return
	}

	response.Success(c, http.StatusOK, "Project berhasil dihapus", nil)
}

// RegisterRoutes registers project management routes on the given router group
func (h *ProjectHandler) RegisterRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.GET("", h.List)
		projects.GET("/:id", h.GetByID)
		projects.POST("", h.Create)
		projects.PUT("/:id", h.Update)
		projects.DELETE("/:id", h.Delete)
	}
}