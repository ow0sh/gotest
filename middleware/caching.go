package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
)

func Caching(c *cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			base := chi.URLParam(r, "base")
			quote := chi.URLParam(r, "quote")

			if value, ok := c.Get(fmt.Sprintf(base + "-" + quote)); ok {
				if byteSlice, ok := value.([]byte); ok {
					fmt.Println("taken from cache")
					w.Header().Set("Content-Type", "application/json")
					w.Write(byteSlice)
					return
				}
			}
			fmt.Println("middlware")
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
