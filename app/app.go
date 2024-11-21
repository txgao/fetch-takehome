package app

import (
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikejav/gosts"
)

type App struct {
	R      *chi.Mux
	Config AppConfig
	Slog   *slog.Logger
}

type Option func(*App)

func NewApp(opts ...Option) *App {
	server := &App{}
	for _, opt := range opts {
		opt(server)
	}
	if server.R == nil {
		server.R = chi.NewRouter()
	}

	server.R.Use(middleware.RequestID)
	server.R.Use(middleware.RealIP)
	server.R.Use(middleware.Recoverer)
	server.R.Use(middleware.NoCache)

	// config for hsts middleware
	hstsConf := &gosts.Info{
		MaxAge:               60 * 60 * 24,
		Expires:              time.Now().Add(24 * time.Hour),
		IncludeSubDomains:    true,
		SendPreloadDirective: false,
	}
	// middleware
	gosts.Configure(hstsConf)
	server.R.Use(gosts.Header)

	return server
}
