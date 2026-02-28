# TODO: External Integrations & Expansion (v2.4)

## Phase 1: Subdomain Discovery (Chaos)
- [x] Create `internal/integrations/chaos/client.go`
- [x] Implement `FetchSubdomains(domain string)` (Simulated)
- [x] Add `ChaosService` to ingest found subdomains
- [x] Update `ScanHandler` to trigger subdomain discovery

## Phase 2: Real Infrastructure Data (IPInfo)
- [x] Create `internal/integrations/ipinfo/client.go`
- [x] Implement `GetIPDetails(ip string)` (Simulated)
- [x] Replace mock `LookupInfrastructure` in `IngestionService` with real logic
- [x] Add caching for IP lookups (via DB updates)

## Phase 3: Live Port Scanning (Naabu/HTTPX)
- [x] Create `internal/integrations/runner/runner.go` to execute CLI tools
- [x] Integrate `httpx` for liveness and tech detection
- [x] Parse `httpx` results and feed into `DomainRepository.AddDetection`

## Phase 4: UI Enhancements for Discovery
- [ ] Add "Discovery" tab to Domain Detail view (Subdomains list)
- [ ] Add "Live Status" badge (Up/Down) based on HTTPX result
