package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tms/tyre/internal/delivery/http/middleware"
	"github.com/tms/tyre/internal/delivery/http/response"
	"github.com/tms/tyre/internal/usecase"
)

type DashboardHandler struct {
	dashboardUseCase *usecase.DashboardUseCase
}

func NewDashboardHandler(dashboardUseCase *usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{dashboardUseCase: dashboardUseCase}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	var companyID *uint
	if claims.Role == "admin_company" {
		companyID = &claims.TenantID
	}

	stats, err := h.dashboardUseCase.GetStats(companyID)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data dashboard")
		return
	}

	response.Success(c, http.StatusOK, "Success", stats)
}

func (h *DashboardHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/dashboard/stats", h.GetStats)
}
