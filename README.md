# SigMap

**SigMap** is a professional-grade technology discovery and infrastructure monitoring dashboard. It enables security researchers, OSINT analysts, and engineering teams to map out technology stacks across thousands of domains, track changes over time (Delta tracking), and identify security risks.

## 🚀 Features

- **Technology Intelligence:** Integrated with 7,500+ technology signatures via ProjectDiscovery.
- **Dynamic Exploration:** HTMX-powered live filtering and search for domains, technologies, and categories.
- **Delta Tracking:** Monitor technical stack changes across your infrastructure in real-time.
- **Investigation Suite:** Built-in technical notes and one-click bookmarking system.
- **Data Export:** Export filtered domain intelligence to CSV for reporting.
- **High Performance:** Built with Go, PostgreSQL, and Materialized Views for low-latency analysis.

## 🛠️ Tech Stack

- **Backend:** Go (Chi Router, pgxpool)
- **Frontend:** HTMX, Go Templates, Tailwind CSS
- **Database:** PostgreSQL (with GIN indexes for fuzzy search)
- **Icons:** Wappalyzer CDN integration

## 🚦 Getting Started

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.21+

### 2. Setup Database
```bash
docker compose up -d
```

### 3. Sync Signatures & Ingest Data
```bash
# Download latest tech signatures
go run cmd/server/main.go -sync

# (Optional) Ingest sample domains for testing
go run cmd/server/main.go -ingest
```

### 4. Run Server
```bash
go run cmd/server/main.go
```
Access the dashboard at [http://localhost:8080](http://localhost:8080)

## 📄 License
MIT
