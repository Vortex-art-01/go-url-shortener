package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/handler/create"
	"url-shortener/internal/handler/get"
	"url-shortener/internal/repository/memstorage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repo := memstorage.New()

	rtr := chi.NewRouter()

	rtr.Use(middleware.RequestID)
	rtr.Use(middleware.Logger)

	rtr.Post("/", create.New(ctx, repo))
	rtr.Get("/{id}", get.New(ctx, repo))

	server := http.Server{Addr: ":8080", Handler: rtr}

	go func() {
		log.Println("server listening on :8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
