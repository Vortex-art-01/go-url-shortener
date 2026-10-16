package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
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

	fmt.Println("server listening on :8080")
	http.ListenAndServe(":8080", rtr)
}
