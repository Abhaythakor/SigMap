package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Abhaythakor/SigMap/internal/repositories"
	"github.com/Abhaythakor/SigMap/internal/services"
	"github.com/Abhaythakor/SigMap/internal/vulnintel"
)

type ScanHandler struct {
	DomainRepo *repositories.DomainRepository
	IngestSvc  *services.IngestionService
	VulnSvc    *vulnintel.Service
	ChaosSvc   *services.ChaosService
	HTTPXSvc   *services.HTTPXService
	NucleiSvc  *services.NucleiService
}

func NewScanHandler(repo *repositories.DomainRepository, ingestSvc *services.IngestionService, vulnSvc *vulnintel.Service, chaosSvc *services.ChaosService, httpxSvc *services.HTTPXService, nucleiSvc *services.NucleiService) *ScanHandler {
	return &ScanHandler{
		DomainRepo: repo,
		IngestSvc:  ingestSvc,
		VulnSvc:    vulnSvc,
		ChaosSvc:   chaosSvc,
		HTTPXSvc:   httpxSvc,
		NucleiSvc:  nucleiSvc,
	}
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

	log.Printf("Triggering full intelligence scan for: %s", domainName)
	
	ctx := r.Context()
	domainID, err := h.DomainRepo.EnsureDomain(ctx, domainName)
	if err != nil {
		http.Error(w, "Failed to ensure domain", http.StatusInternalServerError)
		return
	}

	// 1. Infrastructure Enrichment
	h.IngestSvc.LookupInfrastructure(ctx, domainID, domainName)

	// 2. Subdomain Discovery (Background)
	go h.ChaosSvc.DiscoverSubdomains(ctx, domainName)

	// 3. Live Tech Detection (via HTTPX)
	go func() {
		if err := h.HTTPXSvc.ScanDomain(context.Background(), domainName); err != nil {
			log.Printf("Scan error for %s: %v", domainName, err)
		}
		
		// 4. Active Vulnerability Scan (Nuclei) - Run after tech detection
		if err := h.NucleiSvc.ScanAndStore(context.Background(), domainID, domainName); err != nil {
			log.Printf("Nuclei error for %s: %v", domainName, err)
		}
	}()

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
