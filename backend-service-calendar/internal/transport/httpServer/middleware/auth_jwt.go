package middleware

import (
	"context"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/jwt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func Authorization(secret string) func(next http.Handler) http.Handler {
	op := "middleware.Authorization()"
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			splitToken := strings.Split(authHeader, "Bearer ")
			if len(splitToken) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tokenString := splitToken[1]
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			user, err := jwt.ParseToken(tokenString, secret)
			if err != nil {
				log.Error().Msgf("%s:%w", op, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
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
