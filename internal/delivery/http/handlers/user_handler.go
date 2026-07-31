package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/middleware"
	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
	"github.com/tms/tyre/pkg/utils"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// List returns a paginated list of users
// GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	users, total, err := h.userRepo.List(page, perPage)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data user")
		return
	}

	response.SuccessWithPagination(c, "Success", users, response.NewPagination(page, perPage, total))
}

// GetByID returns a single user by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil || user == nil {
		response.NotFound(c, "User tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Success", user)
}

// Create creates a new user
// POST /api/v1/users
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Check email uniqueness
	existing, _ := h.userRepo.GetByEmail(req.Email)
	if existing != nil {
		response.Error(c, http.StatusConflict, "Email sudah digunakan", nil)
		return
	}

	// Validate role + company_id
	if req.Role == string(entity.RoleAdminCompany) && (req.CompanyID == nil || *req.CompanyID == 0) {
		response.Error(c, http.StatusBadRequest, "Company wajib dipilih untuk Admin Company", nil)
		return
	}
	if req.Role == string(entity.RoleSuperadmin) && req.CompanyID != nil && *req.CompanyID != 0 {
		response.Error(c, http.StatusBadRequest, "Superadmin tidak boleh terkait dengan company", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response.InternalError(c, "Gagal membuat user")
		return
	}

	user := &entity.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       req.Role,
		CompanyID:  req.CompanyID,
		Status:     "active",
	}

	if err := h.userRepo.Create(user); err != nil {
		response.InternalError(c, "Gagal membuat user")
		return
	}

	response.Success(c, http.StatusCreated, "User berhasil dibuat", user)
}

// Update updates an existing user
// PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil || user == nil {
		response.NotFound(c, "User tidak ditemukan")
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// Check email uniqueness (excluding current user)
	existing, _ := h.userRepo.GetByEmail(req.Email)
	if existing != nil && existing.ID != uint(id) {
		response.Error(c, http.StatusConflict, "Email sudah digunakan", nil)
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role
	user.CompanyID = req.CompanyID
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "Gagal mengupdate user")
		return
	}

	response.Success(c, http.StatusOK, "User berhasil diupdate", user)
}

// Delete removes a user
// DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	user, err := h.userRepo.GetByID(uint(id))
	if err != nil || user == nil {
		response.NotFound(c, "User tidak ditemukan")
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		response.InternalError(c, "Gagal menghapus user")
		return
	}

	response.Success(c, http.StatusOK, "User berhasil dihapus", nil)
}

// RegisterRoutes registers user management routes on the given router group
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.Use(middleware.RequireRoles(string(entity.RoleSuperadmin)))
	{
		users.GET("", h.List)
		users.GET("/:id", h.GetByID)
		users.POST("", h.Create)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
