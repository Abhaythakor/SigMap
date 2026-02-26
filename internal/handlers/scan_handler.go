package handlers

import (
	"log"
	"net/http"

	"github.com/Abhaythakor/SigMap/internal/repositories"
	"github.com/Abhaythakor/SigMap/internal/services"
)

type ScanHandler struct {
	DomainRepo *repositories.DomainRepository
	IngestSvc  *services.IngestionService
}

func NewScanHandler(repo *repositories.DomainRepository, ingestSvc *services.IngestionService) *ScanHandler {
	return &ScanHandler{DomainRepo: repo, IngestSvc: ingestSvc}
}

func (h *ScanHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	domainName := r.FormValue("domain")
	if domainName == "" {
		http.Error(w, "Domain required", http.StatusBadRequest)
		return
	}

	log.Printf("Triggering scan for domain: %s", domainName)
	
	ctx := r.Context()
	domainID, err := h.DomainRepo.EnsureDomain(ctx, domainName)
	if err != nil {
		http.Error(w, "Failed to ensure domain", http.StatusInternalServerError)
		return
	}

	// Simulate detection results (integration point for real scanner)
	techs := []string{"React", "Nginx", "Cloudflare", "Google Analytics"}
	for _, t := range techs {
		err = h.DomainRepo.AddDetection(ctx, domainID, t, "https://"+domainName, "v1.0.0", 95, "Live Scanner")
		if err != nil {
			log.Printf("Error adding detection for %s: %v", t, err)
		}
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
