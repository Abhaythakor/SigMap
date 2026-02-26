package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type DeltaHandler struct {
	Repo     *repositories.DomainRepository
	template *template.Template
}

func NewDeltaHandler(repo *repositories.DomainRepository) *DeltaHandler {
	h := &DeltaHandler{Repo: repo}
	h.parseTemplates()
	return h
}

func (h *DeltaHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "delta.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing delta templates: %v", err)
	}
	h.template = tmpl
}

func (h *DeltaHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.ListDelta(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch delta data", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage string
		Deltas      []repositories.DeltaListItem
	}{
		CurrentPage: "delta",
		Deltas:      items,
	}

	if err := h.template.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering delta: %v", err)
	}
}
