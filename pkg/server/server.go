package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	logger *slog.Logger
	// wg     *sync.WaitGroup
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
		// wg:     &sync.WaitGroup{},
	}
}

func (s *Server) Serve(addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
		// SET BY CONFIG???
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.logger.Info("caught signal", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		s.logger.Info("completing background tasks", "addr", srv.Addr)

		// s.wg.Wait()
		shutdownError <- nil
	}()

	/// !!!!!! FIX ENV WHILE SETTTING CONFIGG
	s.logger.Info("starting server", "addr", srv.Addr, "env", "dev!!!")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
