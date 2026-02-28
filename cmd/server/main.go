package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/Abhaythakor/SigMap/internal/database"
	"github.com/Abhaythakor/SigMap/internal/handlers"
	"github.com/Abhaythakor/SigMap/internal/jobs"
	customMiddleware "github.com/Abhaythakor/SigMap/internal/middleware"
	"github.com/Abhaythakor/SigMap/internal/repositories"
	"github.com/Abhaythakor/SigMap/internal/services"
	"github.com/Abhaythakor/SigMap/internal/vulnintel"
	"github.com/Abhaythakor/SigMap/internal/vulnintel/sources"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Flags
	syncFlag := flag.Bool("sync", false, "Sync technology metadata from Wappalyzer")
	ingestFlag := flag.Bool("ingest", false, "Ingest mock sample data for domains")
	vulnFlag := flag.Bool("vuln", false, "Refresh vulnerability profiles")
	alertFlag := flag.Bool("alert", false, "Run alert worker once")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize Database Connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Services
	vulnConnectors := []vulnintel.SourceConnector{sources.NewNVDConnector()}
	vulnService := vulnintel.NewService(db.Pool, vulnConnectors)
	alertService := services.NewAlertService(db.Pool)

	// Handle Flags
	if *syncFlag {
		if err := services.NewSyncService(db.Pool).Sync(context.Background()); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}
		return
	}

	if *ingestFlag {
		repo := repositories.NewDomainRepository(db.Pool)
		if err := services.NewIngestionService(repo).IngestFromDirectory(context.Background(), "testDir"); err != nil {
			log.Fatalf("Ingestion failed: %v", err)
		}
		return
	}

	if *vulnFlag {
		if err := jobs.NewVulnRefreshJob(db.Pool, vulnService).Run(context.Background()); err != nil {
			log.Fatalf("Vuln refresh failed: %v", err)
		}
		return
	}

	if *alertFlag {
		if err := jobs.NewAlertWorker(db.Pool, alertService).Run(context.Background()); err != nil {
			log.Fatalf("Alert worker failed: %v", err)
		}
		return
	}

	// Background Workers
	go startBackgroundJobs(db.Pool, vulnService, alertService)

	// Repositories & Handlers
	dashboardRepo := repositories.NewDashboardRepository(db.Pool)
	domainRepo := repositories.NewDomainRepository(db.Pool)
	techRepo := repositories.NewTechRepository(db.Pool)
	categoryRepo := repositories.NewCategoryRepository(db.Pool)
	trendRepo := repositories.NewTrendRepo(db.Pool)

	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)
	domainHandler := handlers.NewDomainHandler(domainRepo)
	techHandler := handlers.NewTechHandler(techRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	bookmarkHandler := handlers.NewBookmarkHandler(domainRepo)
	noteHandler := handlers.NewNoteHandler(domainRepo)
	trendHandler := handlers.NewTrendHandler(trendRepo)
	deltaHandler := handlers.NewDeltaHandler(domainRepo)
	exportHandler := handlers.NewExportHandler(domainRepo)
	ingestionService := services.NewIngestionService(domainRepo)
	scanHandler := handlers.NewScanHandler(domainRepo, ingestionService)
	settingsHandler := handlers.NewSettingsHandler(domainRepo)

	// Initialize Router
	r := chi.NewRouter()
	r.Use(middleware.Logger, customMiddleware.Recovery, middleware.Recoverer, middleware.Compress(5), middleware.RealIP, middleware.CleanPath, customMiddleware.SanitizeInput)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes
	r.Get("/", dashboardHandler.ServeHTTP)
	r.Get("/domains", domainHandler.List)
	r.Get("/domains/{id}", domainHandler.Detail)
	r.Get("/technologies", techHandler.List)
	r.Get("/categories", categoryHandler.List)
	r.Get("/bookmarks", bookmarkHandler.List)
	r.Post("/bookmarks/toggle", bookmarkHandler.Toggle)
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", noteHandler.List)
		r.Get("/new", noteHandler.New)
		r.Post("/", noteHandler.Create)
		r.Get("/{id}/edit", noteHandler.Edit)
		r.Post("/{id}", noteHandler.Update)
		r.Delete("/{id}", noteHandler.Delete)
	})
	r.Get("/trends", trendHandler.List)
	r.Get("/delta", deltaHandler.List)
	r.Get("/export/domains", exportHandler.Domains)
	r.Post("/scan", scanHandler.Trigger)
	r.Get("/settings/alerts", settingsHandler.AlertsView)
	r.Post("/settings/alerts", settingsHandler.AddChannel)
	r.Delete("/settings/alerts/{id}", settingsHandler.DeleteChannel)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("OK")) })

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Printf("SigMap starting on port %s...", port)
	http.ListenAndServe(":"+port, r)
}

func startBackgroundJobs(pool *pgxpool.Pool, vulnSvc *vulnintel.Service, alertSvc *services.AlertService) {
	alertWorker := jobs.NewAlertWorker(pool, alertSvc)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := alertWorker.Run(context.Background()); err != nil {
			log.Printf("Background alert worker error: %v", err)
		}
	}
}
