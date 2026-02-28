# TODO: Vulnerability Intelligence Integration (v2.1)

## Phase 1: Foundation & Data Layer
- [x] Database Migration: Create `technology_vuln_profile` table
- [x] Define VulnProfile and VulnFinding models in `internal/vulnintel/models.go`
- [x] Implement NVD Source Connector (Simulated)
- [x] Implement Vulners Source Connector (Simulated)
- [x] Implement Correlation Engine (`internal/vulnintel/correlator.go`)
- [x] Implement Risk Scoring Logic (`internal/vulnintel/risk.go`)

## Phase 2: Services & Background Jobs
- [x] Create Vulnerability Service (`internal/vulnintel/service.go`)
- [x] Implement Background Refresh Job (`internal/jobs/vuln_refresh.go`)
- [x] Integrate Vulnerability mapping into the enrichment pipeline
- [x] Implement Caching mechanism (via DB profile storage)

## Phase 3: UI & API Enhancements
- [x] Expose internal API: `GET /internal/vuln/{technology}`
- [x] Technologies View: Add Risk Level, CVE Count, and Exploit Available columns
- [x] Domains View: Add Stack Risk Summary
- [x] Dashboard Metrics: Add Critical Technologies counts

## Phase 4: Advanced Sources & Alerts (Future)
- [ ] GitHub POC detection connector
- [ ] OSV & VulnCheck connectors
- [x] Real-time Alerting for new critical vulnerabilities (AlertWorker)
