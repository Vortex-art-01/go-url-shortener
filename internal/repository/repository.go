package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func (r *Repository) CreateShortURL(ctx context.Context, url, shortURL string) error {
	var exists bool
	err := r.conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM url WHERE full_url = $1)", url).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("url already exists")
	}

	_, err = r.conn.Exec(ctx, "INSERT INTO url (short_url, full_url) VALUES ($1, $2)", shortURL, url)
	if err != nil {
		return fmt.Errorf("insert url: %w", err)
	}

	return nil
}
func (r Repository) GetURL(ctx context.Context, shortURL string) (string, error) {
	row := r.conn.QueryRow(ctx, "SELECT full_url FROM url WHERE short_url = $1", shortURL)
	var fullURL string
	err := row.Scan(&fullURL)
	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func (r *Repository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func New(ctx context.Context) *Repository {
	conn, err := pgx.Connect(ctx, "postgres://postgres:admin@localhost:5432/url_shortener?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}

	return &Repository{conn: conn}
}
