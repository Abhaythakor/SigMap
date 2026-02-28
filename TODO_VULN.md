# TODO: Vulnerability Intelligence Integration (v2.1)

## Phase 1: Foundation & Data Layer
- [ ] Database Migration: Create `technology_vuln_profile` table
- [ ] Define VulnProfile and VulnFinding models in `internal/vulnintel/models.go`
- [ ] Implement NVD Source Connector (`internal/vulnintel/sources/nvd.go`)
- [ ] Implement Vulners Source Connector (`internal/vulnintel/sources/vulners.go`)
- [ ] Implement Correlation Engine (`internal/vulnintel/correlator.go`)
- [ ] Implement Risk Scoring Logic (`internal/vulnintel/risk.go`)

## Phase 2: Services & Background Jobs
- [ ] Create Vulnerability Service (`internal/vulnintel/service.go`)
- [ ] Implement Background Refresh Job (`internal/jobs/vuln_refresh.go`)
- [ ] Integrate Vulnerability mapping into the enrichment pipeline
- [ ] Implement Caching mechanism (24h TTL)

## Phase 3: UI & API Enhancements
- [ ] Expose internal API: `GET /internal/vuln/{technology}`
- [ ] Technologies View: Add Risk Level, CVE Count, and Exploit Available columns
- [ ] Domains View: Add Stack Risk Summary (e.g., High Risk Tech: 2)
- [ ] Dashboard Metrics: Add Risky and Critical Technologies counts

## Phase 4: Advanced Sources & Alerts (Future)
- [ ] GitHub POC detection connector
- [ ] OSV & VulnCheck connectors
- [ ] Real-time Alerting for new critical vulnerabilities
