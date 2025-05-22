package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "go-todo/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/ping", s.PingHandler)
	mux.HandleFunc("/health", s.healthHandler)

	// Only serve Swagger UI if not in production
	if os.Getenv("APP_ENV") != "production" {
		mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	}

	// Todo routes
	mux.HandleFunc("/todos", s.getTodosHandler)
	mux.HandleFunc("/todo/", s.getTodoHandler)
	mux.HandleFunc("/todo/create", s.createTodoHandler)
	mux.HandleFunc("/todo/update/", s.updateTodoHandler)
	mux.HandleFunc("/todo/delete/", s.deleteTodoHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// @Summary: PingHandler godoc
// @Description: Ping the server
// @Tags: ping
// @Produce: json
// @Success: 200 {object} map[string]string
// @Router: /ping [get]
func (s *Server) PingHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "pong"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// @Summary: HealthHandler godoc
// @Description: Check the health of the server
// @Tags: health
// @Produce: json
// @Success: 200 {object} map[string]string
// @Router: /health [get]
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
