package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("api-gateway starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to run gateway: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down api-gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("api-gateway forced shutdown: %v", err)
	}

	log.Println("api-gateway exited gracefully")
}
