package main

import (
	"net/http"
)

// cors wraps an HTTP handler with CORS (Cross-Origin Resource Sharing) middleware.
// It sets appropriate CORS headers for cross-origin requests and handles preflight OPTIONS requests.
// The middleware allows credentials and sets a max age of 3600 seconds for preflight caching.
func cors(h http.Handler) http.Handler {
	allowed := map[string]bool{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowed[origin] {
			// Always vary by Origin so caches don't mix responses
			w.Header().Add("Vary", "Origin")

			// Exact echo of the allowed origin (not "*") so credentials can work
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Methods you actually support
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			// Echo requested headers, or provide a safe default set
			reqHdrs := r.Header.Get("Access-Control-Request-Headers")
			if reqHdrs != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHdrs)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With")
			}

			// Optional: cache the preflight for an hour
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight early
		if r.Method == http.MethodOptions {
			// If you want to restrict paths, you can check r.URL.Path here
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
