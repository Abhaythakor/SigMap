# SigMap Technical Testing Guide (CLI & Logic)

Follow these steps to verify that the backend engine, database, and integrations are working correctly.

## 🛠️ Prerequisites
- [ ] PostgreSQL is running: `docker compose up -d`
- [ ] Database is accessible on port `5433` (as per `.env`)
- [ ] Binary is built: `go build -o server cmd/server/main.go`

---

## 🏁 Step 1: Technology Signature Sync
**Goal:** Verify connection to ProjectDiscovery and metadata ingestion.
- [ ] Run command: `DB_PORT=5433 ./server -sync`
- [ ] **Expected Result:** Logs show "Starting Wappalyzer metadata sync" and complete successfully.
- [ ] **Verification:** Run `docker exec -i webtechview_db psql -U admin -d webtechview -c "SELECT count(*) FROM technologies;"`. Count should be > 7,000.

## 📂 Step 2: Real Data Ingestion
**Goal:** Ingest JSON files from `testDir` and verify parsing.
- [ ] Run command: `DB_PORT=5433 ./server -ingest`
- [ ] **Expected Result:** Logs show "Ingesting file: testDir/test1.json" and "Ingestion from testDir completed."
- [ ] **Verification:** Run `docker exec -i webtechview_db psql -U admin -d webtechview -c "SELECT name, version FROM detections d JOIN domains dom ON d.domain_id = dom.id LIMIT 5;"`. You should see real domains and versions (or empty versions).

## 🛡️ Step 3: Vulnerability Intelligence Refresh
**Goal:** Test correlation engine and multi-source (NVD/Vulners) simulation.
- [ ] Run command: `DB_PORT=5433 ./server -vuln`
- [ ] **Expected Result:** Logs show "Refreshing real vulnerability details" for various technologies.
- [ ] **Verification:** Run `SELECT count(*) FROM vulnerability_details;`. It should be > 0.

## 🔔 Step 4: Real-time Alert Worker
**Goal:** Test the Sentry mode logic and webhook dispatching.
- [ ] First, add a dummy webhook in the Web UI (Settings > Alerts).
- [ ] Run command: `DB_PORT=5433 ./server -alert`
- [ ] **Expected Result:** Logs show "Starting alert worker pass." If new high-risk tech was found, it should show "ALERT TRIGGERED."
- [ ] **Verification:** Check `alert_history` table in DB to see if the event was logged.

## 🔍 Step 5: Live Scanning Integration
**Goal:** Verify HTTPX, Chaos, and IPInfo discovery chain.
- [ ] Start server: `DB_PORT=5433 ./server`
- [ ] From another terminal, trigger a scan: `curl -X POST -d "domain=google.com" http://localhost:8080/scan`
- [ ] **Expected Result:** Check server logs. You should see:
    1. "Triggering full intelligence scan"
    2. "IPInfo: Fetching details"
    3. "Chaos: Fetching subdomains"
    4. "HTTPX: Scanning google.com"
