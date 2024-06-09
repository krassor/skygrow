package middleware

import (
	"net/http"
)

func Authorization(accessToken string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
