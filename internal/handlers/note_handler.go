package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/user/webtechview/internal/repositories"
)

type NoteHandler struct {
	Repo      *repositories.DomainRepository
	templates map[string]*template.Template
}

func NewNoteHandler(repo *repositories.DomainRepository) *NoteHandler {
	h := &NoteHandler{Repo: repo, templates: make(map[string]*template.Template)}
	h.parseTemplates()
	return h
}

func (h *NoteHandler) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}
}

func (h *NoteHandler) parseTemplates() {
	funcMap := h.getFuncMap()

	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "notes.html"),
		filepath.Join("templates", "partials", "toast.html"),
	}
	tmpl := template.New("base").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error parsing note templates: %v", err)
	}
	h.templates["index"] = tmpl

	// Form partial
	formFiles := []string{
		filepath.Join("templates", "partials", "note_form.html"),
		filepath.Join("templates", "partials", "toast.html"),
	}
	tmplForm := template.New("form").Funcs(funcMap)
	tmplForm, err = tmplForm.ParseFiles(formFiles...)
	if err != nil {
		log.Fatalf("Error parsing note form partial: %v", err)
	}
	h.templates["form"] = tmplForm
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Repo.ListNotes(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage string
		Notes       []repositories.NoteListItem
	}{
		CurrentPage: "notes",
		Notes:       items,
	}

	if err := h.templates["index"].ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering notes: %v", err)
	}
}

func (h *NoteHandler) New(w http.ResponseWriter, r *http.Request) {
	domainIDStr := r.URL.Query().Get("domain_id")
	domainID, _ := strconv.Atoi(domainIDStr)

	data := struct {
		DomainID int
		IsEdit   bool
	}{
		DomainID: domainID,
		IsEdit:   false,
	}

	if err := h.templates["form"].ExecuteTemplate(w, "note_form", data); err != nil {
		log.Printf("Error rendering note form: %v", err)
	}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	domainID, _ := strconv.Atoi(r.FormValue("domain_id"))
	content := r.FormValue("content")
	author := r.FormValue("author")

	if err := h.Repo.CreateNote(r.Context(), domainID, content, author); err != nil {
		http.Error(w, "Failed to save note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/notes")
	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if err := h.Repo.DeleteNote(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	// Trigger removal from UI
	w.WriteHeader(http.StatusOK)
}
