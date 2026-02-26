package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type CategoryHandler struct {
	Repo     *repositories.CategoryRepository
	template *template.Template
}

func NewCategoryHandler(repo *repositories.CategoryRepository) *CategoryHandler {
	h := &CategoryHandler{Repo: repo}
	h.parseTemplates()
	return h
}

func (h *CategoryHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "categories.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing category templates: %v", err)
	}
	h.template = tmpl
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.List(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage string
		Categories  []repositories.CategoryListItem
	}{
		CurrentPage: "categories",
		Categories:  items,
	}

	if err := h.template.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering categories: %v", err)
	}
}
