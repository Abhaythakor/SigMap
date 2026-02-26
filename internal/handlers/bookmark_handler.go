package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type BookmarkHandler struct {
	Repo     *repositories.DomainRepository
	templates map[string]*template.Template
}

func NewBookmarkHandler(repo *repositories.DomainRepository) *BookmarkHandler {
	h := &BookmarkHandler{Repo: repo, templates: make(map[string]*template.Template)}
	h.parseTemplates()
	return h
}

func (h *BookmarkHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "bookmarks.html"),
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "toast.html"),
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing bookmark templates: %v", err)
	}
	h.templates["index"] = tmpl

	// Partial for HTMX toggle
	buttonFiles := []string{
		filepath.Join("templates", "partials", "bookmark_button.html"),
		filepath.Join("templates", "partials", "toast.html"),
	}
	tmplButton, err := template.ParseFiles(buttonFiles...)
	if err != nil {
		log.Fatalf("Error parsing bookmark button partial: %v", err)
	}
	h.templates["button"] = tmplButton
}

func (h *BookmarkHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.ListBookmarks(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch bookmarks", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage string
		Bookmarks   []repositories.BookmarkListItem
	}{
		CurrentPage: "bookmarks",
		Bookmarks:   items,
	}

	if err := h.templates["index"].ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering bookmarks: %v", err)
	}
}

func (h *BookmarkHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	isBookmarked, err := h.Repo.ToggleBookmark(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to toggle bookmark", http.StatusInternalServerError)
		return
	}

	data := struct {
		ID           int
		IsBookmarked bool
		Message      string
		Type         string
		OOB          bool
	}{
		ID:           id,
		IsBookmarked: isBookmarked,
		Message:      "Bookmark updated successfully",
		Type:         "success",
		OOB:          true,
	}

	// Execute button template (main swap)
	if err := h.templates["button"].ExecuteTemplate(w, "bookmark_button", data); err != nil {
		log.Printf("Error rendering bookmark button: %v", err)
	}

	// Execute toast template (OOB swap)
	if err := h.templates["button"].ExecuteTemplate(w, "toast", data); err != nil {
		log.Printf("Error rendering toast: %v", err)
	}
}
