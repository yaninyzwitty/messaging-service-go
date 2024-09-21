package middleware

import "net/http"

// var allowedOrigins = map[string]bool{
// 	"https://example.com":  true,
// 	"https://example1.com": true,
// }

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// origin := r.Header.Get("Origin")

		// if allowedOrigins[origin] {
		// Set CORS headers for allowed origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// }
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		//  pass the request to the next handler

		next.ServeHTTP(w, r)
	})
}
