package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type ExportHandler struct {
	Repo *repositories.DomainRepository
}

func NewExportHandler(repo *repositories.DomainRepository) *ExportHandler {
	return &ExportHandler{Repo: repo}
}

func (h *ExportHandler) Domains(w http.ResponseWriter, r *http.Request) {
	filters := repositories.DomainFilters{
		Search:       r.URL.Query().Get("search"),
		Category:     r.URL.Query().Get("category"),
		Confidence:   r.URL.Query().Get("confidence"),
		IsBookmarked: r.URL.Query().Get("bookmarked") == "true",
	}

	// Fetch all matching records (limit 10000 for export)
	items, err := h.Repo.List(r.Context(), 10000, 0, filters)
	if err != nil {
		log.Printf("Export error: %v", err)
		http.Error(w, "Failed to fetch data for export", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=domains_export.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header
	writer.Write([]string{"Domain", "Technologies", "Categories", "Confidence", "Last Seen"})

	for _, item := range items {
		techNames := make([]string, len(item.Technologies))
		for i, t := range item.Technologies {
			techNames[i] = t.Name
		}

		writer.Write([]string{
			item.Name,
			strings.Join(techNames, "; "),
			strings.Join(item.Categories, "; "),
			fmt.Sprintf("%d%%", item.Confidence),
			item.LastSeen,
		})
	}
}
