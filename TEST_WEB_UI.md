# SigMap Web UI Testing Guide (Step-by-Step)

Start the server (`DB_PORT=5433 ./server`) and open [http://localhost:8080](http://localhost:8080).

---

## 📊 1. Dashboard Overview
- [ ] **Metrics:** Do the 5 metric cards (Total Detections, Avg Confidence, Risky Tech, Critical Threats, Bookmarks) show real numbers?
- [ ] **Trends Chart:** Is the line chart rendering detection counts over the last 7 days?
- [ ] **Distribution:** Is the doughnut chart showing the Top 5 Technologies?

## 🌐 2. Domain Explorer
- [ ] **Live Search:** Type "stripe" in the filter bar. Does the table update instantly without refreshing the page?
- [ ] **Multi-filtering:** Select "Confidence: High" and check "Bookmarked." Do results narrow down correctly?
- [ ] **Pagination:** Scroll to the bottom. Use "Next" and "Previous" buttons. Does the URL update (e.g., `?page=2`)?
- [ ] **Icons & Versions:** Do you see the Wappalyzer icons and the `vX.Y.Z` text next to technologies?

## 🔎 3. Domain Deep-Dive (Detail Page)
- [ ] **Navigation:** Click on a domain name (e.g., `stripe.com`).
- [ ] **Live Badge:** Does it show a pulsing green "Live" badge if tech was detected?
- [ ] **Vulnerability Cards:** Expand/View the "Active Vulnerabilities" section. Do you see Nuclei findings with severity colors?
- [ ] **CVE Intel:** Do you see descriptions and bug types under the "CVE Intelligence" sections for specific techs?
- [ ] **Subdomains:** Are the subdomains discovered by Chaos listed? Click one—does it redirect to that subdomain's detail page?

## 📑 4. Investigation Tools
- [ ] **Bookmarking:** Click a bookmark icon in the table. Does it turn blue immediately? Does a success toast appear in the bottom-right?
- [ ] **Notes CRUD:**
    - [ ] Click "Add Note" on a domain. Does the modal appear?
    - [ ] Save a note. Does it redirect you to the Notes page?
    - [ ] On the Notes page, click "Edit" on a note. Change text and save. Does it update?
    - [ ] Click "Delete" on a note. Does the card disappear from the UI?

## ⚙️ 5. Alert Settings
- [ ] **Configuration:** Go to Sidebar > Settings.
- [ ] **Add Channel:** Add a test webhook (e.g., `https://httpbin.org/post`).
- [ ] **Verification:** Does the new channel appear in the "Active Integrations" list?
- [ ] **Deletion:** Delete the test channel. Does it remove instantly via HTMX?

## 📥 6. Export Intelligence
- [ ] **CSV Export:** Go to Domains and click "Export CSV."
- [ ] **Verification:** Open the downloaded file. Does it contain all filtered domains with their technology stacks?
