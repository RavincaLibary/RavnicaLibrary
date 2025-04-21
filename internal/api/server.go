package api

import (
	"RavnicaLibrary/internal/cards"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	router  *mux.Router
	httpSrv *http.Server
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	r := mux.NewRouter()

	srv := &Server{
		router: r,
		httpSrv: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}

	// Register routes
	srv.registerRoutes()

	return srv
}

// registerRoutes sets up all the routes for the server
func (s *Server) registerRoutes() {
	// Card handlers
	cardHandler := cards.NewHandler()
	cardHandler.RegisterRoutes(s.router)

	// Add middleware
	s.router.Use(loggingMiddleware)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s", s.httpSrv.Addr)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Server is shutting down...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped")
	return nil
}

// loggingMiddleware logs information about each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
