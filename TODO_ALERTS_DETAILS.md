# TODO: Advanced Monitoring & Intelligence (v2.2)

## Phase 1: Domain Deep Dive
- [x] Repository: Implement `GetDomainDetails` (Aggregated view of history, metadata, and notes)
- [x] Handler: Implement `DomainDetail` view handler
- [x] Template: Create `templates/domain_detail.html` with:
    - [x] Technology Timeline (Scan history)
    - [x] Infrastructure metadata (IP, Source)
    - [x] Linked Notes component
- [x] UI: Link domain names in tables to the detail page

## Phase 2: Notification & Alerting
- [x] Define `AlertChannel` model (Webhook, Slack, Discord)
- [x] Implement `NotificationService` for dispatching alerts
- [x] Create `internal/jobs/alert_worker.go`:
    - [x] Scan for new Critical/High risks in the last hour
    - [x] Dispatch to active channels
- [ ] UI: Basic Alert Configuration page (`/settings/alerts`)

## Phase 3: Infrastructure Enrichment
- [ ] Integration: ASN & Cloud Provider lookup (via IP)
- [x] UI: Show Cloud Provider icons (AWS, GCP, Azure) next to domains (Infrastructure Grid implemented)
