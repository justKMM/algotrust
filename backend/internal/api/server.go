package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"rationalgo/internal/config"
	"rationalgo/internal/repository"
)

// Server is the HTTP API for the audit dashboard.
type Server struct {
	cfg   config.Config
	store *repository.Store
	mux   *http.ServeMux
}

// NewServer creates an API server with seeded in-memory state.
func NewServer(cfg config.Config) *Server {
	s := &Server{
		cfg:   cfg,
		store: repository.NewStore(),
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /api/state", s.handleState)
	s.mux.HandleFunc("POST /api/state/reset", s.handleReset)
}

// ListenAndServe starts the HTTP server with CORS for local frontend dev.
func (s *Server) ListenAndServe() error {
	addr := s.cfg.HTTPAddr
	log.Printf("RationAlgo API listening on %s", addr)
	return http.ListenAndServe(addr, s.withCORS(s.mux))
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"phase":  "1",
	})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.State())
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.store.Reset()
	writeJSON(w, http.StatusOK, s.store.State())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Fprintf(w, `{"error":%q}`, strings.ReplaceAll(err.Error(), `"`, `\"`))
	}
}
