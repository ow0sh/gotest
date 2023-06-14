package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info(fmt.Sprintf("Method: [%v]; Address: [%v], URL: [%v] in %v", r.Method, r.RemoteAddr, r.URL, time.Since(start)))
		}

		return http.HandlerFunc(fn)
	}
}
