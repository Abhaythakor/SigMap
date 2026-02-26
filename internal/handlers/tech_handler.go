package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type TechHandler struct {
	Repo     *repositories.TechRepository
	template *template.Template
}

func NewTechHandler(repo *repositories.TechRepository) *TechHandler {
	h := &TechHandler{Repo: repo}
	h.parseTemplates()
	return h
}

func (h *TechHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "technologies.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing tech templates: %v", err)
	}
	h.template = tmpl
}

func (h *TechHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context(), 50, 0)
	if err != nil {
		http.Error(w, "Failed to fetch technologies", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage  string
		Technologies []repositories.TechListItem
	}{
		CurrentPage:  "technologies",
		Technologies: items,
	}

	if err := h.template.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering technologies: %v", err)
	}
}
