package usecase

import (
	"errors"
	"time"

	"github.com/tms/tyre/internal/domain/entity"
	"github.com/tms/tyre/internal/domain/repository"
	"github.com/tms/tyre/internal/dto/request"
	"github.com/tms/tyre/internal/dto/response"
	"github.com/tms/tyre/internal/infrastructure/jwt"
	"github.com/tms/tyre/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
	ErrCompanyNotFound    = errors.New("company not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthUseCase struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtService *jwt.JWTService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *AuthUseCase) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if user.Role == string(entity.RoleAdminCompany) && user.CompanyID == nil {
		return nil, ErrCompanyNotFound
	}

	claims := jwt.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
	if user.CompanyID != nil {
		claims.TenantID = *user.CompanyID
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(claims)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         NewUserDTO(user),
	}, nil
}

func (uc *AuthUseCase) RefreshToken(req *request.RefreshTokenRequest) (*response.LoginResponse, error) {
	claims, err := uc.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrExpiredToken) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	user, err := uc.userRepo.GetByID(claims.UserID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	newClaims := jwt.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
	if user.CompanyID != nil {
		newClaims.TenantID = *user.CompanyID
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(newClaims)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         NewUserDTO(user),
	}, nil
}

func (uc *AuthUseCase) GetProfile(userID uint) (*entity.User, error) {
	return uc.userRepo.GetByID(userID)
}

func (uc *AuthUseCase) ChangePassword(userID uint, req *request.ChangePasswordRequest) error {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return uc.userRepo.Update(user)
}

// NewUserDTO creates a UserDTO from an entity.User
func NewUserDTO(user *entity.User) *response.UserDTO {
	return &response.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		Status:    user.Status,
	}
}

// Ensure time is used (Resolve linter warning for unused import)
var _ = time.Time{}
