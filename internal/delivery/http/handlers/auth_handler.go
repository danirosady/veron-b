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
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "Email atau kata sandi salah", nil)
			return
		}
		if errors.Is(err, usecase.ErrUserInactive) {
			response.Error(c, http.StatusUnauthorized, "Akun Anda tidak aktif. Hubungi administrator", nil)
			return
		}
		if errors.Is(err, usecase.ErrCompanyNotFound) {
			response.Error(c, http.StatusUnauthorized, "Data perusahaan tidak ditemukan. Hubungi administrator", nil)
			return
		}
		response.InternalError(c, "Login gagal. Silakan coba lagi nanti")
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
			response.Error(c, http.StatusUnauthorized, "Session expired. Please login again", nil)
			return
		}
		response.Error(c, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	response.Success(c, http.StatusOK, "Session refreshed", result)
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
		response.NotFound(c, "User not found")
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
			response.Error(c, http.StatusBadRequest, "Current password is incorrect", nil)
			return
		}
		response.InternalError(c, "Failed to change password. Please try again")
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

// RegisterPublicRoutes registers auth routes on the given router group
func (h *AuthHandler) RegisterPublicRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.GET("/profile", h.GetProfile)
		auth.PUT("/password", h.ChangePassword)
	}
}
