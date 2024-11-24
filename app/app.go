package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	R          *chi.Mux
	Config     AppConfig
	Slog       *slog.Logger
	HttpLogger *httplog.Logger
}

func NewApp() *App {

	// Configuration
	var appConfig AppConfig
	cleanenv.ReadEnv(&appConfig)
	logger := httplog.NewLogger("httplog", httplog.Options{
		JSON:             false,
		LogLevel:         slog.LevelInfo,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/healthz",
		},
		QuietDownPeriod: 600 * time.Second,
	})

	r := chi.NewRouter()
	server := &App{
		R:          r,
		Config:     appConfig,
		HttpLogger: logger,
	}

	server.R.Use(middleware.NoCache)
	server.R.Use(httplog.RequestLogger(logger))

	return server
}

func (app *App) Run() {
	addr := fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port)
	server := &http.Server{Addr: addr, Handler: app.R}

	slog.Info("Started server.", "addr", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed starting server", "err", err)
		}
	}()

	// Capturing signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for SIGINT (kill -2)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed shutdown server", "err", err)
	}
	slog.Info("Server exited")

}

func HealthCheck(r *chi.Mux) {
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, http.StatusText(http.StatusOK))
	})
}
