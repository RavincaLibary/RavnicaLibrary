package cards

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// Handler handles HTTP requests related to cards
type Handler struct {
	service *CardService
}

// NewHandler creates a new card handler
func NewHandler() *Handler {
	return &Handler{
		service: NewCardService(nil),
	}
}

// RegisterRoutes registers the card routes with the provided router
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/cards/search", h.SearchCards).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/cards/{id}", h.GetCardByID).Methods(http.MethodGet, http.MethodOptions)
}

// SearchCards handles the card search endpoint
func (h *Handler) SearchCards(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		h.setCORSHeaders(w)
		return
	}

	// Get the query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing required 'q' parameter", http.StatusBadRequest)
		return
	}

	// Search for cards
	result, err := h.service.SearchCards(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	h.setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	// Return the search result as JSON
	json.NewEncoder(w).Encode(result)
}

// GetCardByID handles fetching a specific card by ID
func (h *Handler) GetCardByID(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		h.setCORSHeaders(w)
		return
	}

	// Get card ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Card ID is required", http.StatusBadRequest)
		return
	}

	// Get card details
	card, err := h.service.GetCardByID(id)
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Card not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	h.setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	// Return the card as JSON
	json.NewEncoder(w).Encode(card)
}

// setCORSHeaders sets CORS headers on the response
func (h *Handler) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
