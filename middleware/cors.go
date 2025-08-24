package middleware

import (
	"net/http"
)

// EnableCORS middleware để xử lý Cross-Origin Resource Sharing
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cho phép requests từ localhost:3000 (React app)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		// Cho phép các HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Cho phép các headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Cho phép credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Xử lý preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Tiếp tục xử lý request
		next.ServeHTTP(w, r)
	})
}
