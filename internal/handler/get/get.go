package get

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type shortURLGetter interface {
	GetUrl(ctx context.Context, shortURL string) (string, error)
}

func New(ctx context.Context, getter shortURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if len(id) == 0 {
			fmt.Println("id is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url, err := getter.GetUrl(ctx, id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
