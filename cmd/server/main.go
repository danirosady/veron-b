package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tms/tyre/configs"
	httpdelivery "github.com/tms/tyre/internal/delivery/http"
	"github.com/tms/tyre/internal/infrastructure/database"
	"github.com/tms/tyre/internal/infrastructure/logger"
)

// runDownMigrationOnly applies the last down migration and exits
func runDownMigrationOnly() {
	cfg := configs.LoadFromEnv()
	db, err := database.NewPostgres(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	if err := database.RunDownMigrations(db); err != nil {
		fmt.Fprintf(os.Stderr, "Migration down failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Migration down completed.")
}

// runUpMigrationOnly runs pending up migrations and exits
func runUpMigrationOnly() {
	cfg := configs.LoadFromEnv()
	db, err := database.NewPostgres(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	if err := database.RunMigrations(db); err != nil {
		fmt.Fprintf(os.Stderr, "Migration up failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Migration up completed.")
}

func main() {
	// Check for CLI-only flags before starting the server
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--migrate-down":
			runDownMigrationOnly()
			return
		case "--migrate-up":
			runUpMigrationOnly()
			return
		}
	}

	cfg := configs.LoadFromEnv()

	logger.InitLogger(cfg.App.Env)
	logger.Info("Starting Tyre Management System server...")

	db, err := database.NewPostgres(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("Database connection established")

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := database.RunMigrations(db); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Database migrations applied")

	ginServer := httpdelivery.NewServer(cfg, db)
	router := ginServer.Setup()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server listening on %s:%s", cfg.App.Host, cfg.App.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error("Error closing database connection", "error", err)
	}

	logger.Info("Server exited gracefully")
	log.Println("Server stopped")
}
