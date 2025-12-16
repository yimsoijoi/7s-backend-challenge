package http

import (
	"log"
	"net/http"
	"time"

	"github.com/yimsoijoi/7s-backend-challenge/internal/infrastructure"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func Auth(jwt infrastructure.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		userID, err := jwt.Validate(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Add("user-id", userID)
		next.ServeHTTP(w, r)
	})
}
