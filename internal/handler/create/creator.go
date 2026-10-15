package create

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

type shortURLCreator interface {
	CreateShortUrl(ctx context.Context, url, shortURL string) error
}

func New(ctx context.Context, creator shortURLCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("err cannot read body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		if len(url) == 0 {
			fmt.Println("err empty body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCode, err := NewCode(6)
		if err != nil {
			fmt.Println("err NewCode", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = creator.CreateShortUrl(ctx, string(url), newCode)
		if err != nil {
			fmt.Println("err CreateShortUrl", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("r,uri", r.URL.RequestURI())

		w.Write([]byte("http://localhost:8080/" + newCode))
	}
}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewCode(n int) (string, error) {
	max := big.NewInt(int64(len(alphabet)))
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idx.Int64()]
	}
	return string(b), nil
}
