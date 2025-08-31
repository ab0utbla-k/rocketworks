package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/ab0utbla-k/rocketworks/pkg/utils"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"environment": "dev",
			"version":     "0.1.0",
		},
	}

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_ = utils.WriteJSON(w, http.StatusOK, data, nil)
	})

	srv := &http.Server{
		Addr:         net.JoinHostPort("localhost", "8080"),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// change
		log.Println("caught signal", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// change
		log.Println("completing background tasks", "addr", srv.Addr)
		shutdownError <- nil
	}()

	log.Println("starting server", "addr", srv.Addr, "env", "dev")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Boom")
		os.Exit(1)
	}

	err = <-shutdownError
	if err != nil {
		log.Fatalln("Boom2")
		os.Exit(1)
	}
}
