package main

import (
	"bank_app/internal/app"
	"bank_app/internal/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			os.Exit(1)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	deps := app.InitDependencies(cfg)
	defer deps.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: deps.Router.GetEngine(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	serverError := make(chan error, 1)

	go func() {
		deps.Logger.Info("Server is starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			deps.Logger.Error("Failed to run server:", "error", err)
			serverError <- err
		}
	}()

	select {
	case <-quit:
		deps.Logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			deps.Logger.Error("Server forced to shutdown:", "error", err)
		}
	case err := <-serverError:
		deps.Logger.Error("Server error:", "error", err)
	}

	deps.Logger.Info("Server exited")
}
