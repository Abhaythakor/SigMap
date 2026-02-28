# TODO: External Integrations & Expansion (v2.4)

## Phase 1: Subdomain Discovery (Chaos)
- [ ] Create `internal/integrations/chaos/client.go`
- [ ] Implement `FetchSubdomains(domain string)` using Chaos API (or CLI wrapper)
- [ ] Add `ChaosService` to ingest found subdomains into `domains` table
- [ ] Update `ScanHandler` to trigger subdomain discovery

## Phase 2: Real Infrastructure Data (IPInfo)
- [ ] Create `internal/integrations/ipinfo/client.go`
- [ ] Implement `GetIPDetails(ip string)` (ASN, Geo, Cloud)
- [ ] Replace mock `LookupInfrastructure` in `IngestionService` with real IPInfo calls
- [ ] Add caching for IP lookups to avoid rate limits

## Phase 3: Live Port Scanning (Naabu/HTTPX)
- [ ] Create `internal/integrations/runner/runner.go` to execute CLI tools
- [ ] Integrate `httpx` for liveness and tech detection (json output)
- [ ] Parse `httpx` results and feed into `DomainRepository.AddDetection`

## Phase 4: UI Enhancements for Discovery
- [ ] Add "Discovery" tab to Domain Detail view (Subdomains list)
- [ ] Add "Live Status" badge (Up/Down) based on HTTPX result
