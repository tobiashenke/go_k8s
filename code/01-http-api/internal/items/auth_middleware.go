package items

import (
	"net/http"
	"strings"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
