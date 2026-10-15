package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"url-shortener/internal/handler/create"
	"url-shortener/internal/handler/get"
	"url-shortener/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repo := repository.New(ctx)
	defer repo.Close(ctx)

	rtr := chi.NewRouter()

	rtr.Use(middleware.RequestID)
	rtr.Use(middleware.Logger)

	rtr.Post("/", create.New(ctx, repo))
	rtr.Get("/{id}", get.New(ctx, repo))

	http.ListenAndServe(":8080", rtr)
}
