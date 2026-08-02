package http

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/tms/tyre/configs"
	"github.com/tms/tyre/internal/delivery/http/handlers"
	"github.com/tms/tyre/internal/delivery/http/middleware"
	"github.com/tms/tyre/internal/infrastructure/jwt"
	"github.com/tms/tyre/internal/infrastructure/repository"
	"github.com/tms/tyre/internal/usecase"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *configs.Config) {
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	repos := repository.NewRepositories(db)

	// Use cases
	authUseCase := usecase.NewAuthUseCase(repos.User, jwtService)
	companyUseCase := usecase.NewCompanyUseCase(repos.Company, repos.Project, repos.Unit, repos.Driver)
	projectUseCase := usecase.NewProjectUseCase(repos.Project, repos.Company, repos.Unit)
	unitUseCase := usecase.NewUnitUseCase(repos.Unit, repos.Project, repos.Company, repos.Tyre, repos.Master)
	driverUseCase := usecase.NewDriverUseCase(repos.Driver, repos.Company)
	tyreUseCase := usecase.NewTyreUseCase(repos.Tyre, repos.Company, repos.Master)
	masterUseCase := usecase.NewMasterUseCase(repos.Master)
	replacementUseCase := usecase.NewReplacementUseCase(repos.Replacement, repos.Tyre, repos.Unit, repos.Driver)
	reportUseCase := usecase.NewReportUseCase(db, repos.Tyre, repos.Replacement, repos.Unit, repos.Company)
	dashboardUseCase := usecase.NewDashboardUseCase(db)

	// Handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(repos.User)
	companyHandler := handlers.NewCompanyHandler(companyUseCase)
	projectHandler := handlers.NewProjectHandler(projectUseCase)
	unitHandler := handlers.NewUnitHandler(unitUseCase)
	driverHandler := handlers.NewDriverHandler(driverUseCase)
	tyreHandler := handlers.NewTyreHandler(tyreUseCase)
	masterHandler := handlers.NewMasterHandler(masterUseCase)
	replacementHandler := handlers.NewReplacementHandler(replacementUseCase)
	reportHandler := handlers.NewReportHandler(reportUseCase)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUseCase)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "Tyre Management System"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(jwtService))
		{
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/password", authHandler.ChangePassword)

			userHandler.RegisterRoutes(protected)
			companyHandler.RegisterRoutes(protected)
			projectHandler.RegisterRoutes(protected)
			unitHandler.RegisterRoutes(protected)
			driverHandler.RegisterRoutes(protected)
			tyreHandler.RegisterRoutes(protected)
			masterHandler.RegisterRoutes(protected)
			replacementHandler.RegisterRoutes(protected)
			reportHandler.RegisterRoutes(protected)
			dashboardHandler.RegisterRoutes(protected)
		}
	}
}