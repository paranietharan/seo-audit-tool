package main

import (
	"encoding/json"
	"log"
	"net/http"
	"seo-audit-tool/config"
	"seo-audit-tool/internal/analyzer"
	"seo-audit-tool/internal/crawler"
	"seo-audit-tool/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	config   *configs.Config
	db       *storage.Database
	crawler  *crawler.Crawler
	analyzer *analyzer.Analyzer
}

type AuditRequest struct {
	URL string `json:"url"`
}

type AuditResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func NewServer() *Server {
	config := configs.NewConfig()
	db := storage.NewDatabase(config.DatabaseURL)

	c := crawler.NewCrawler()
	a := analyzer.NewAnalyzer()

	return &Server{
		config:   config,
		db:       db,
		crawler:  c,
		analyzer: a,
	}
}

func (s *Server) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/audit", s.handleAudit).Methods("POST")
	r.HandleFunc("/report/{id}", s.handleReport).Methods("GET")
	return r
}

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	var req AuditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	audit := &storage.Audit{
		URL:    req.URL,
		Status: "processing",
	}

	if err := s.db.CreateAudit(audit); err != nil {
		http.Error(w, "Failed to create audit", http.StatusInternalServerError)
		return
	}

	go s.runAudit(audit)

	response := AuditResponse{
		ID:     audit.ID,
		URL:    audit.URL,
		Status: audit.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	audit, err := s.db.GetAudit(id)
	if err != nil {
		http.Error(w, "Audit not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audit)
}

func (s *Server) runAudit(audit *storage.Audit) {
	defer func() {
		if r := recover(); r != nil {
			audit.Status = "failed"
			audit.Error = "Internal error occurred"
			s.db.UpdateAudit(audit)
		}
	}()

	// Crawl the website
	pages, err := s.crawler.CrawlSite(audit.URL, 10) // Limit to 10 pages for demo
	if err != nil {
		audit.Status = "failed"
		audit.Error = err.Error()
		s.db.UpdateAudit(audit)
		return
	}

	// Analyze pages
	results := make([]*storage.PageResult, 0, len(pages))
	for _, page := range pages {
		result := s.analyzer.AnalyzePage(page)
		result.AuditID = audit.ID
		results = append(results, result)
		s.db.CreatePageResult(result)
	}

	// Update audit status
	audit.Status = "completed"
	audit.PagesAnalyzed = len(results)
	s.db.UpdateAudit(audit)
}

func main() {
	server := NewServer()

	if err := server.db.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := server.setupRoutes()

	log.Printf("Server starting on port %s", server.config.Port)
	log.Fatal(http.ListenAndServe(":"+server.config.Port, router))
}
