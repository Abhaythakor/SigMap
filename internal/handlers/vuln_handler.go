package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Abhaythakor/SigMap/internal/vulnintel"
)

type VulnHandler struct {
	VulnService *vulnintel.Service
}

func NewVulnHandler(svc *vulnintel.Service) *VulnHandler {
	return &VulnHandler{VulnService: svc}
}

func (h *VulnHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	technology := chi.URLParam(r, "technology")
	log.Printf("DEBUG: VulnHandler.GetProfile called for technology: '%s'", technology)
	
	if technology == "" {
		http.Error(w, "Technology parameter required", http.StatusBadRequest)
		return
	}

	profile, err := h.VulnService.GetTechVulnProfile(r.Context(), technology)
	if err != nil {
		log.Printf("DEBUG: VulnService.GetTechVulnProfile error: %v", err)
		http.Error(w, "Failed to fetch vulnerability profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
