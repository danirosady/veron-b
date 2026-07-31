package http

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/tms/tyre/configs"
	"github.com/tms/tyre/internal/delivery/http/middleware"
)

type Server struct {
	engine *gin.Engine
	cfg    *configs.Config
	db     *gorm.DB
}

func NewServer(cfg *configs.Config, db *gorm.DB) *Server {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Server{
		engine: engine,
		cfg:    cfg,
		db:     db,
	}
}

func (s *Server) Setup() *gin.Engine {
	s.engine.Use(middleware.CORS(s.cfg.App.Origins))
	s.engine.Use(middleware.RequestID())
	s.engine.Use(middleware.Logger())
	s.engine.Use(middleware.Recovery())

	RegisterRoutes(s.engine, s.db, s.cfg)

	return s.engine
}
