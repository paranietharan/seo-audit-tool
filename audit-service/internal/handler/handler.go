package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paranietharan/seo-audit-tool/audit-service/internal/scraper"
)

type Handler struct {
	scraper *scraper.Scraper
}

func NewHandler(s *scraper.Scraper) *Handler {
	return &Handler{scraper: s}
}

type AuditRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *Handler) Audit(c *gin.Context) {
	var req AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	targetURL := strings.TrimSpace(req.URL)
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	result, err := h.scraper.Scrape(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
