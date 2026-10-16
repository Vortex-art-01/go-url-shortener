package memstorage

import (
	"errors"
	"fmt"
)

type MemRepository struct {
	storage map[string]string
}

func (r *MemRepository) GetShortURL(url string) (string, error) {
	for s, u := range r.storage {
		if u == url {
			return s, nil
		}
	}

	return "", errors.New("shortUrl not found")
}

func (r *MemRepository) CreateShortURL(url, shortURL string) error {
	r.storage[shortURL] = url
	return nil
}

func (r MemRepository) GetURL(shortURL string) (string, error) {
	fmt.Println("shortURL", shortURL)
	fmt.Println(r.storage)
	if url, ok := r.storage[shortURL]; ok {
		fmt.Println("GetUrl 1")

		return url, nil
	}

	fmt.Println("GetUrl 2")
	return "", errors.New("url not found")
}

func New() *MemRepository {
	return &MemRepository{storage: make(map[string]string)}
}
