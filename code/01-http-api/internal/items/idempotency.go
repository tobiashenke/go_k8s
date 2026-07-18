package items

import (
	"encoding/json"
	"net/http"
	"time"
)

type cachedResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func IdempotencyMiddleware(c *ItemCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Response format of the cache for GET and SAVE
			var data cachedResponse
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			response, err := c.GetResponse(r.Context(), key)
			// if key present return early
			if err == nil {
				json.Unmarshal(response, &data)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(data.Status)
				w.Write([]byte(data.Body))
				return
			}
			// if key not present, store key in cache
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			next.ServeHTTP(rw, r)
			data = cachedResponse{
				Status: rw.statusCode,
				Body:   rw.body.String(),
			}
			body, _ := json.Marshal(data)
			c.SetResponse(r.Context(), key, body, 24*time.Hour)
		})
	}
}
