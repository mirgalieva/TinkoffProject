package httpgin

import (
	"context"
	"homework9/middleware"
	"net/http"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

type Server struct {
	svr *http.Server
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api := router.Group("/api/v1")
	api.Use(middleware.Logger)
	api.Use(middleware.Recover)
	AppRouter(api, a)
	return Server{&http.Server{Addr: port, Handler: router}}
}

func (s *Server) Handler() http.Handler {
	return s.svr.Handler
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.svr.Shutdown(ctx)
}

func (s *Server) ListenAndServe() error {
	return s.svr.ListenAndServe()
}
