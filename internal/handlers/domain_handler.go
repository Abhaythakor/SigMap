package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

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

	// Pre-parse the main view
	baseFiles := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "domains.html"),
		filepath.Join("templates", "partials", "domain_rows.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "pagination.html"),
	}
	tmpl := template.New("base").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(baseFiles...)
	if err != nil {
		log.Fatalf("Error parsing domain templates: %v", err)
	}
	h.templates["index"] = tmpl

	// Pre-parse the partial view
	partialFiles := []string{
		filepath.Join("templates", "partials", "domain_rows.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "pagination.html"),
	}
	tmplPartial := template.New("rows").Funcs(funcMap)
	tmplPartial, err = tmplPartial.ParseFiles(partialFiles...)
	if err != nil {
		log.Fatalf("Error parsing domain partial templates: %v", err)
	}
	h.templates["rows"] = tmplPartial
}

func (h *DomainHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse Filters
	filters := repositories.DomainFilters{
		Search:       r.URL.Query().Get("search"),
		Category:     r.URL.Query().Get("category"),
		Confidence:   r.URL.Query().Get("confidence"),
		IsBookmarked: r.URL.Query().Get("bookmarked") == "true",
	}

	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	items, err := h.Repo.List(r.Context(), limit, offset, filters)
	if err != nil {
		log.Printf("Error fetching domains: %v", err)
		http.Error(w, "Failed to fetch domains", http.StatusInternalServerError)
		return
	}

	total, _ := h.Repo.Count(r.Context(), filters)
	totalPages := (total + limit - 1) / limit

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
		TotalPages:  totalPages,
		TotalItems:  total,
		Limit:       limit,
	}

	isHX := r.Header.Get("HX-Request") == "true"

	if isHX {
		err = h.templates["rows"].ExecuteTemplate(w, "domain_rows", data)
		if err == nil {
			err = h.templates["rows"].ExecuteTemplate(w, "pagination", data)
		}
	} else {
		err = h.templates["index"].ExecuteTemplate(w, "base", data)
	}

	if err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}
}
