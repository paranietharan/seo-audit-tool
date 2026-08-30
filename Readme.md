# 🚀 Technical SEO Auditor & AI Analyzer

A production-grade, distributed microservices platform for automated technical website audits, DOM scraping, Core Web Vitals metric extraction, and AI-driven SEO recommendation reports powered by **Google Gemini** and headless **Chromium (Go-Rod)**.

---

## 🏗️ Architecture Overview

```
                          ┌──────────────────────────┐
                          │   Frontend (React/Vite)  │
                          │        Port: 3000        │
                          └─────────────┬────────────┘
                                        │
                                        │ HTTP / REST
                                        ▼
                          ┌──────────────────────────┐
                          │   API Gateway (Go/Gin)   │
                          │        Port: 8080        │
                          └──────┬────────────┬──────┘
                                 │            │
             POST /api/audit     │            │ POST /api/report
             (Internal / Proxy)  │            │
                                 ▼            ▼
             ┌─────────────────────────┐   ┌───────────────────────────────┐
             │  Audit Service (Go/Rod) │◄──┤ Report Service (FastAPI)      │
             │       Port: 8081        │   │        Port: 8082             │
             └─────────────────────────┘   └───────────────┬───────────────┘
                                                           │
                                                           │ Google GenAI API
                                                           ▼
                                            ┌───────────────────────────────┐
                                            │ Google Gemini 2.0 Flash / Pro │
                                            └───────────────────────────────┘
```

### Services Breakdown
1. **`api-gateway`** *(Go / Gin)*: Central entrypoint handling path routing, CORS header harmonization, upstream failover error formatting, and reverse proxying to backend services.
2. **`audit-service`** *(Go / Go-Rod / GoQuery)*: Headless Chromium browser automation engine that executes JavaScript, extracts DOM metadata, calculates accurate word counts, resolves canonical/OG URLs, verifies heading hierarchies, and protects against SSRF.
3. **`report-service`** *(Python / FastAPI / Google GenAI)*: Orchestration layer that triggers raw audits and synthesizes technical SEO signals into executive summaries, categorized issue tables, and prioritized quick wins using Google Gemini (with deterministic rule-based fallback).
4. **`frontend`** *(React 18 / Vite 6 / Tailwind CSS / Nginx)*: Interactive dashboard featuring live SEO score gauges, issue severity filters, raw DOM inspection accordions, localStorage history, and JSON export.

---

## ⚡ Quick Start with Docker Compose

### 1. Clone the repository & set up environment
```bash
git clone https://github.com/paranietharan/seo-audit-tool.git
cd seo-audit-tool
cp .env.example .env
```

Edit `.env` and add your Google Gemini API key:
```env
GEMINI_API_KEY=your_gemini_api_key_here
```

### 2. Launch the entire microservices stack
```bash
docker compose up --build
```

### 3. Access the services
- **Frontend Dashboard**: [http://localhost:3000](http://localhost:3000)
- **API Gateway**: [http://localhost:8080](http://localhost:8080)
- **Audit Service**: [http://localhost:8081/health](http://localhost:8081/health)
- **Report Service Docs**: [http://localhost:8082/docs](http://localhost:8082/docs)

---

## 🛠️ Local Development Setup

### 1. API Gateway (`api-gateway`)
```bash
cd api-gateway
go run cmd/main.go
# Starts on port 8080
```

### 2. Audit Service (`audit-service`)
```bash
cd audit-service
go run cmd/main.go
# Starts on port 8081 (requires local Chrome or Chromium installed)
```

### 3. Report Service (`report-service`)
```bash
cd report-service
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
export GEMINI_API_KEY="your-api-key"
uvicorn main:app --host 0.0.0.0 --port 8082 --reload
```

### 4. Frontend (`frontend`)
```bash
cd frontend
npm install
npm run dev
# Starts on http://localhost:3000
```

---

## 📡 API Reference

### 1. Generate Full Audit & AI Report
`POST /api/report` (or `POST :8082/report`)

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response (200 OK):**
```json
{
  "audit": {
    "url": "https://example.com",
    "title": "Example Domain",
    "meta_description": "",
    "canonical_url": "https://example.com",
    "robots_meta": "",
    "viewport_meta": "width=device-width, initial-scale=1",
    "h1_count": 1,
    "h2_count": 0,
    "heading_issues": [],
    "word_count": 125,
    "page_load_ms": 320,
    "status_code": 200,
    "has_structured_data": false,
    "internal_links": [],
    "external_links": [{"href": "https://www.iana.org/domains/example", "text": "More information..."}],
    "images_missing_alt": []
  },
  "report": {
    "score": 82,
    "summary": "The website has a clean DOM structure with fast response time, but lacks meta description and structured data.",
    "issues": [
      {
        "severity": "warning",
        "category": "Meta Tags",
        "description": "Missing meta description tag.",
        "recommendation": "Add a compelling meta description between 120-160 characters."
      }
    ],
    "quick_wins": [
      "Add a meta description to improve SERP click-through rates."
    ],
    "strengths": [
      "Single H1 heading correctly implemented.",
      "Fast page load time (320 ms)."
    ]
  }
}
```

### 2. Generate Raw Technical Audit Only
`POST /api/audit` (or `POST :8081/audit`)

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

### 3. Health Checks
`GET /health` on any service returns `{"status": "ok"}`.

---