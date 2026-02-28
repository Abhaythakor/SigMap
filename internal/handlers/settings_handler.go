package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Abhaythakor/SigMap/internal/models"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type SettingsHandler struct {
	Repo      *repositories.DomainRepository
	templates map[string]*template.Template
}

func NewSettingsHandler(repo *repositories.DomainRepository) *SettingsHandler {
	h := &SettingsHandler{
		Repo:      repo,
		templates: make(map[string]*template.Template),
	}
	h.parseTemplates()
	return h
}

func (h *SettingsHandler) parseTemplates() {
	files := []string{
		filepath.Join("templates", "layouts", "base.html"),
		filepath.Join("templates", "partials", "sidebar.html"),
		filepath.Join("templates", "partials", "header.html"),
		filepath.Join("templates", "settings_alerts.html"),
		filepath.Join("templates", "partials", "toast.html"),
	}
	h.templates["alerts"] = template.Must(template.New("base").ParseFiles(files...))
}

func (h *SettingsHandler) AlertsView(w http.ResponseWriter, r *http.Request) {
	channels, err := h.Repo.ListAlertChannels(r.Context())
	if err != nil {
		http.Error(w, "Failed to load channels", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentPage string
		Channels    []models.AlertChannel
	}{
		CurrentPage: "settings",
		Channels:    channels,
	}

	if err := h.templates["alerts"].ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error rendering settings: %v", err)
	}
}

func (h *SettingsHandler) AddChannel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	cType := r.FormValue("type")
	url := r.FormValue("url")

	if err := h.Repo.AddAlertChannel(r.Context(), name, cType, url); err != nil {
		http.Error(w, "Failed to add channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/settings/alerts")
	w.WriteHeader(http.StatusOK)
}

func (h *SettingsHandler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if err := h.Repo.DeleteAlertChannel(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete channel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
