# SEO Audit Tool

A four-service web application for crawling web pages, checking technical SEO signals, and generating structured audit reports using **Go (Gin & Go-Rod)**, **Python (FastAPI)**, **React (Vite & Tailwind CSS)**, and **Google Gemini**.

---

## Architecture

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
                                  ▼            ▼
              ┌─────────────────────────┐   ┌───────────────────────────────┐
              │  Audit Service (Go/Rod) │◄──┤ Report Service (FastAPI)      │
              │       Port: 8081        │   │        Port: 8082             │
              └─────────────────────────┘   └───────────────┬───────────────┘
                                                            │
                                                            │ Google GenAI SDK
                                                            ▼
                                             ┌───────────────────────────────┐
                                             │   Google Gemini 2.0 Flash     │
                                             │  (+ Rule-Based Score Fallback)│
                                             └───────────────────────────────┘
```

### Services
1. **`api-gateway` (`Go / Gin`, `:8080`)**: Reverse proxy that handles CORS, request validation, and routing between the React frontend, `audit-service`, and `report-service`.
2. **`audit-service` (`Go / Go-Rod / GoQuery`, `:8081`)**: Headless Chromium crawler with a 4-tab concurrency semaphore. Renders client-side JavaScript pages, extracts meta tags, canonical/OpenGraph URLs, heading hierarchy (`H1`/`H2`), links, and images missing `alt` attributes, and validates resolved IPs before crawling to block private/loopback SSRF targets.
3. **`report-service` (`Python / FastAPI`, `:8082`)**: Calls `audit-service` to collect raw page metrics and prompts Google Gemini (`gemini-2.0-flash`) for a structured JSON SEO report. Includes a deterministic rule-based scoring fallback (`_fallback_rule_based_report`) so reports still complete if the Gemini API times out or hits quota limits.
4. **`frontend` (`React 18 / Vite / Tailwind CSS`, `:3000`)**: Web dashboard for running audits, filtering issues by severity, inspecting raw DOM metrics, viewing local audit history, and exporting reports as JSON.

---

## Quick Start with Docker Compose

### 1. Configure Environment
```bash
git clone https://github.com/paranietharan/seo-audit-tool.git
cd seo-audit-tool
cp .env.example .env
```

Set your Gemini API key in `.env` (if omitted or quota-exhausted, `report-service` falls back to its built-in rule-based scorer):
```env
GEMINI_API_KEY=your_gemini_api_key_here
```

### 2. Start All Services
```bash
docker compose up --build
```

- **Frontend Dashboard**: `http://localhost:3000`
- **API Gateway**: `http://localhost:8080`
- **Audit Service Health**: `http://localhost:8081/health`
- **Report Service OpenAPI Docs**: `http://localhost:8082/docs`

---

## Running Services Locally

### 1. API Gateway (`api-gateway`)
```bash
cd api-gateway
go run cmd/main.go
# Listens on :8080
```

### 2. Audit Service (`audit-service`)
Requires Chrome or Chromium installed locally:
```bash
cd audit-service
go run cmd/main.go
# Listens on :8081
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
# Opens on http://localhost:3000
```

---

## API Reference

### 1. Full Crawl & SEO Report
`POST /api/report` (via Gateway `:8080` or directly on `:8082/report`)

**Request:**
```json
{
  "url": "https://example.com"
}
```

**Response (`200 OK`):**
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
    "external_links": [
      {
        "href": "https://www.iana.org/domains/example",
        "text": "More information..."
      }
    ],
    "images_missing_alt": []
  },
  "report": {
    "score": 82,
    "summary": "The page loads quickly and has a single H1 tag, but is missing a meta description and structured data.",
    "issues": [
      {
        "severity": "warning",
        "category": "Meta Tags",
        "description": "Missing meta description tag.",
        "recommendation": "Add a meta description between 120 and 160 characters."
      }
    ],
    "quick_wins": [
      "Add a meta description to improve search snippet click-through rates."
    ],
    "strengths": [
      "Single H1 heading present.",
      "Fast page load time (320 ms)."
    ]
  }
}
```

### 2. Raw Technical Crawl Only
`POST /api/audit` (via Gateway `:8080` or `:8081/audit`)

```json
{
  "url": "https://example.com"
}
```

### 3. Health Check
`GET /health` on `:8080`, `:8081`, or `:8082` returns `{"status": "ok"}`.
