package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type DomainHandler struct {
	Repo      *repositories.DomainRepository
	templates map[string]*template.Template
}

func NewDomainHandler(repo *repositories.DomainRepository) *DomainHandler {
	h := &DomainHandler{
		Repo:      repo,
		templates: make(map[string]*template.Template),
	}
	h.parseTemplates()
	return h
}

func (h *DomainHandler) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
	}
}

func (h *DomainHandler) parseTemplates() {
	funcMap := h.getFuncMap()

	// List View
	baseFiles := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "domains.html"),
		filepath.Join("templates", "partials", "domain_rows.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "pagination.html"),
	}
	h.templates["index"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(baseFiles...))

	// Detail View
	detailFiles := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "domain_detail.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
	}
	h.templates["detail"] = template.Must(template.New("base").Funcs(funcMap).ParseFiles(detailFiles...))

	// Partial rows
	rowFiles := []string{
		filepath.Join("templates", "partials", "domain_rows.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "pagination.html"),
	}
	h.templates["rows"] = template.Must(template.New("rows").Funcs(funcMap).ParseFiles(rowFiles...))
}

func (h *DomainHandler) List(w http.ResponseWriter, r *http.Request) {
	filters := repositories.DomainFilters{
		Search:       r.URL.Query().Get("search"),
		Category:     r.URL.Query().Get("category"),
		Confidence:   r.URL.Query().Get("confidence"),
		IsBookmarked: r.URL.Query().Get("bookmarked") == "true",
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 { page = 1 }
	limit := 20
	offset := (page - 1) * limit

	items, err := h.Repo.List(r.Context(), limit, offset, filters)
	if err != nil {
		log.Printf("Error fetching domains: %v", err)
		http.Error(w, "Failed to fetch domains", http.StatusInternalServerError)
		return
	}

	total, _ := h.Repo.Count(r.Context(), filters)
	data := struct {
		CurrentPage string
		Domains     []repositories.DomainListItem
		Filters     repositories.DomainFilters
		Page        int
		TotalPages  int
		TotalItems  int
		Limit       int
	}{
		CurrentPage: "domains",
		Domains:     items,
		Filters:     filters,
		Page:        page,
		TotalPages:  (total + limit - 1) / limit,
		TotalItems:  total,
		Limit:       limit,
	}

	if r.Header.Get("HX-Request") == "true" {
		h.templates["rows"].ExecuteTemplate(w, "domain_rows", data)
		h.templates["rows"].ExecuteTemplate(w, "pagination", data)
	} else {
		h.templates["index"].ExecuteTemplate(w, "base", data)
	}
}

func (h *DomainHandler) Detail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	detail, err := h.Repo.GetDomainDetails(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching domain details: %v", err)
		http.Error(w, "Domain not found", http.StatusNotFound)
		return
	}

	data := struct {
		CurrentPage string
		Domain      repositories.DomainDetail
	}{
		CurrentPage: "domains",
		Domain:      detail,
	}

	if err := h.templates["detail"].ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering detail: %v", err)
	}
}

func (h *DomainHandler) RedirectByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	var id int
	err := h.Repo.Pool.QueryRow(r.Context(), "SELECT id FROM domains WHERE name = $1", name).Scan(&id)
	if err != nil {
		http.Redirect(w, r, "/domains", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/domains/%d", id), http.StatusSeeOther)
}
