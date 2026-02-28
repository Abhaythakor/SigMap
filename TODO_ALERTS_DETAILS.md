# TODO: Advanced Monitoring & Intelligence (v2.2)

## Phase 1: Domain Deep Dive
- [ ] Repository: Implement `GetDomainDetails` (Aggregated view of history, metadata, and notes)
- [ ] Handler: Implement `DomainDetail` view handler
- [ ] Template: Create `templates/domain_detail.html` with:
    - [ ] Technology Timeline (Scan history)
    - [ ] Infrastructure metadata (IP, Source)
    - [ ] Linked Notes component
- [ ] UI: Link domain names in tables to the detail page

## Phase 2: Notification & Alerting
- [ ] Define `AlertChannel` model (Webhook, Slack, Discord)
- [ ] Implement `NotificationService` for dispatching alerts
- [ ] Create `internal/jobs/alert_worker.go`:
    - [ ] Scan for new Critical/High risks in the last hour
    - [ ] Dispatch to active channels
- [ ] UI: Basic Alert Configuration page (`/settings/alerts`)

## Phase 3: Infrastructure Enrichment
- [ ] Integration: ASN & Cloud Provider lookup (via IP)
- [ ] UI: Show Cloud Provider icons (AWS, GCP, Azure) next to domains
