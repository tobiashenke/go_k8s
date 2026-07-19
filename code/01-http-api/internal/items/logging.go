package items

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrId := fmt.Sprintf("%d", time.Now().UnixNano())
		slog.Info("request", "correlation_id", corrId, "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
