package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Authorization(accessToken string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("%s %s %s Body: %s\n", r.RemoteAddr, r.Method, r.URL, r.Body)
		next.ServeHTTP(w, r)
	})
}
