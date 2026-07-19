package items

import (
	"net/http"
	"os"
	"strings"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// disable authentication for local development
		if os.Getenv("AUTH_DISABLED") == "true" {
			next.ServeHTTP(w, r)
			return
		}

		var prefix string = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, prefix) {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, prefix)
		_, err := ValidateToken(tokenString)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
