package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paranietharan/seo-audit-tool/api-gateway/internal/proxy"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	auditURL := os.Getenv("AUDIT_SERVICE_URL")
	if auditURL == "" {
		auditURL = "http://localhost:8081"
	}

	reportURL := os.Getenv("REPORT_SERVICE_URL")
	if reportURL == "" {
		reportURL = "http://localhost:8082"
	}

	p, err := proxy.NewProxy(auditURL, reportURL)
	if err != nil {
		log.Fatalf("failed to initialize proxy: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(proxy.CORS())

	r.POST("/api/audit", p.AuditHandler)
	r.POST("/api/report", p.ReportHandler)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("api-gateway starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
