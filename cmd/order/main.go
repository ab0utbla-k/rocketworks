package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/ab0utbla-k/rocketworks/pkg/httputil"
	"github.com/ab0utbla-k/rocketworks/pkg/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	// Can be safely removed if I use content type in WriteJSON. Think about it
	r.Use(render.SetContentType(render.ContentTypeJSON))

	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"environment": "dev",
			"version":     "0.1.0",
		},
	}

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		err := httputil.WriteJSON(w, http.StatusOK, data, nil)
		if err != nil {
			httputil.LogError(logger, r, err)
		}
	})

	srv := server.New(logger)
	addr := net.JoinHostPort("localhost", "8080")
	err := srv.Serve(addr, r)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
