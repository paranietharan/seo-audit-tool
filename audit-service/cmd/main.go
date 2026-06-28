package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paranietharan/seo-audit-tool/audit-service/internal/handler"
	"github.com/paranietharan/seo-audit-tool/audit-service/internal/scraper"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		log.Fatalf("failed to initialize scraper: %v", err)
	}
	defer scraperInstance.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	h := handler.NewHandler(scraperInstance)

	r.POST("/audit", h.Audit)
	r.GET("/health", h.Health)

	log.Printf("audit-service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
