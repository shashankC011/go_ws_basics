package main

// import "net/http"

// func enableCORS(w http.ResponseWriter) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all (use specific domain in prod)
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// }

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		enableCORS(w)

// 		// Handle preflight request
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
