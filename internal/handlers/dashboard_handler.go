package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type DashboardHandler struct {
	Repo     *repositories.DashboardRepository
	template *template.Template
}

func NewDashboardHandler(repo *repositories.DashboardRepository) *DashboardHandler {
	h := &DashboardHandler{Repo: repo}
	h.parseTemplates()
	return h
}

func (h *DashboardHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "dashboard.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing dashboard templates: %v", err)
	}
	h.template = tmpl
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.GetStats(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch dashboard stats", http.StatusInternalServerError)
		return
	}

	trends, _ := h.Repo.GetTrendData(r.Context())
	dist, _ := h.Repo.GetDistributionData(r.Context())

	data := struct {
		CurrentPage  string
		Stats        repositories.DashboardStats
		Trends       []repositories.TrendPoint
		Distribution []repositories.DistributionPoint
	}{
		CurrentPage:  "dashboard",
		Stats:        stats,
		Trends:       trends,
		Distribution: dist,
	}

	if err := h.template.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering dashboard: %v", err)
	}
}
