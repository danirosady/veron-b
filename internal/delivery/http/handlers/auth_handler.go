package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/middleware"
	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/dto/request"
	"github.com/tms/tyre/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// Login handles user authentication
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authUseCase.Login(&req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) || errors.Is(err, usecase.ErrUserInactive) {
			response.Error(c, http.StatusUnauthorized, "Email atau password salah", nil)
			return
		}
		if errors.Is(err, usecase.ErrCompanyNotFound) {
			response.Error(c, http.StatusUnauthorized, "Data company tidak ditemukan", nil)
			return
		}
		response.InternalError(c, "Gagal login")
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", result)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authUseCase.RefreshToken(&req)
	if err != nil {
		if errors.Is(err, usecase.ErrTokenExpired) {
			response.Error(c, http.StatusUnauthorized, "Token sudah expired", nil)
			return
		}
		response.Error(c, http.StatusUnauthorized, "Token tidak valid", nil)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed", result)
}

// GetProfile returns the authenticated user's profile
// GET /api/v1/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authUseCase.GetProfile(userID.(uint))
	if err != nil || user == nil {
		response.NotFound(c, "User tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "Success", user)
}

// ChangePassword allows a user to change their own password
// PUT /api/v1/auth/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserID)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.authUseCase.ChangePassword(userID.(uint), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			response.Error(c, http.StatusBadRequest, "Password lama salah", nil)
			return
		}
		response.InternalError(c, "Gagal mengubah password")
		return
	}

	response.Success(c, http.StatusOK, "Password berhasil diubah", nil)
}

// RegisterPublicRoutes registers auth routes on the given router group
func (h *AuthHandler) RegisterPublicRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.GET("/profile", h.GetProfile)
		auth.PUT("/password", h.ChangePassword)
	}
}
