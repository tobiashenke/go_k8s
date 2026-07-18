package items

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func RateLimitMiddleware(c *ItemCache, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* r.RemoteAddr will be the proxy's IP, if the requestors API is behind
			a proxy or load balancer (like in k8s with an Ingress)
			To get the real client IP then use r.Header.Get("X-Forwarded-For") instead */
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			exceeded, err := c.ExceededRateLimitWithLua(r.Context(), ip, limit, window)
			if !exceeded || err != nil {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
		})
	}
}
