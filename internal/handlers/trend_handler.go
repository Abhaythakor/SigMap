package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type TrendHandler struct {
	Repo     *repositories.TrendRepo
	template *template.Template
}

func NewTrendHandler(repo *repositories.TrendRepo) *TrendHandler {
	h := &TrendHandler{Repo: repo}
	h.parseTemplates()
	return h
}

func (h *TrendHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "trends.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing trend templates: %v", err)
	}
	h.template = tmpl
}

func (h *TrendHandler) List(w http.ResponseWriter, r *http.Request) {
	trends, err := h.Repo.GetTrends(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch trends", http.StatusInternalServerError)
		return
	}

	velocity, _ := h.Repo.GetVelocityData(r.Context())

	data := struct {
		CurrentPage string
		Trends      []repositories.TrendStat
		Velocity    []repositories.TrendPoint
	}{
		CurrentPage: "trends",
		Trends:      trends,
		Velocity:    velocity,
	}

	if err := h.template.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering trends: %v", err)
	}
}
