package main

import (
	"context"
	"fetch-takehome/app"
	"fmt"
	"net/url"
	"os"

	receiptH "fetch-takehome/api/receipt"
	receiptSvc "fetch-takehome/pkg/receipt"

	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Port     uint16 `env:"PORT" env-default:"5433"`
	Host     string `env:"HOST" env-default:"localhost"`
	Database string `env:"DATABASE" env-default:"postgres"`
	User     string `env:"USER" env-default:"content"`
	Password string `env:"PASSWORD" env-default:"pwd"`
}

func (c DBConfig) toDatabaseUrl() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Database,
	}
	return u.String()
}

func main() {

	var config DBConfig
	cleanenv.ReadEnv(&config)

	server := app.NewApp()
	dbconn, err := pgxpool.New(context.Background(), config.toDatabaseUrl())
	if err != nil {
		os.Exit(-1)
	}
	receiptSvc := receiptSvc.NewService(dbconn)

	receiptHandler := receiptH.Handle{
		ReceiptService: receiptSvc,
	}

	app.HealthCheck(server.R)
	server.R.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Mount("/receipts", receiptH.Handler(receiptHandler))
		})
	})

	server.Run()
}
