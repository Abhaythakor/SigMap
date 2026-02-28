package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

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
)

func main() {
	// Flags
	syncFlag := flag.Bool("sync", false, "Sync technology metadata from Wappalyzer")
	ingestFlag := flag.Bool("ingest", false, "Ingest mock sample data for domains")
	vulnFlag := flag.Bool("vuln", false, "Refresh vulnerability profiles for detected technologies")
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

	// Vulnerability Service Initialization
	vulnConnectors := []vulnintel.SourceConnector{
		sources.NewNVDConnector(),
	}
	vulnService := vulnintel.NewService(db.Pool, vulnConnectors)

	// Handle Sync
	if *syncFlag {
		syncService := services.NewSyncService(db.Pool)
		if err := syncService.Sync(context.Background()); err != nil {
			log.Fatalf("Sync failed: %v", err)
		}
		log.Println("Sync completed. Exiting.")
		return
	}

	// Handle Ingest
	if *ingestFlag {
		domainRepo := repositories.NewDomainRepository(db.Pool)
		ingestionService := services.NewIngestionService(domainRepo)
		if err := ingestionService.IngestSampleData(context.Background()); err != nil {
			log.Fatalf("Ingestion failed: %v", err)
		}
		log.Println("Mock ingestion completed. Exiting.")
		return
	}

	// Handle Vuln Refresh
	if *vulnFlag {
		vulnJob := jobs.NewVulnRefreshJob(db.Pool, vulnService)
		if err := vulnJob.Run(context.Background()); err != nil {
			log.Fatalf("Vulnerability refresh failed: %v", err)
		}
		log.Println("Vulnerability refresh completed. Exiting.")
		return
	}

	// Initialize Repositories
	dashboardRepo := repositories.NewDashboardRepository(db.Pool)
	domainRepo := repositories.NewDomainRepository(db.Pool)
	techRepo := repositories.NewTechRepository(db.Pool)
	categoryRepo := repositories.NewCategoryRepository(db.Pool)
	trendRepo := repositories.NewTrendRepo(db.Pool)

	// Initialize Handlers
	dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)
	domainHandler := handlers.NewDomainHandler(domainRepo)
	techHandler := handlers.NewTechHandler(techRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	bookmarkHandler := handlers.NewBookmarkHandler(domainRepo)
	noteHandler := handlers.NewNoteHandler(domainRepo)
	trendHandler := handlers.NewTrendHandler(trendRepo)
	deltaHandler := handlers.NewDeltaHandler(domainRepo)
	exportHandler := handlers.NewExportHandler(domainRepo)
	scanHandler := handlers.NewScanHandler(domainRepo, nil)

	// Initialize Router
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.Logger)
	r.Use(customMiddleware.Recovery)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.RealIP)
	r.Use(middleware.CleanPath)
	r.Use(customMiddleware.SanitizeInput)

	// Static Files
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
	r.Get("/notes", noteHandler.List)
	r.Get("/notes/new", noteHandler.New)
	r.Post("/notes", noteHandler.Create)
	r.Get("/notes/{id}/edit", noteHandler.Edit)
	r.Post("/notes/{id}", noteHandler.Update)
	r.Delete("/notes/{id}", noteHandler.Delete)
	r.Get("/trends", trendHandler.List)
	r.Get("/delta", deltaHandler.List)
	r.Get("/export/domains", exportHandler.Domains)
	r.Post("/scan", scanHandler.Trigger)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Server Execution
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
