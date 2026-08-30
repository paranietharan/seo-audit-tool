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
	"github.com/paranietharan/seo-audit-tool/audit-service/config"
	"github.com/paranietharan/seo-audit-tool/audit-service/internal/handler"
	"github.com/paranietharan/seo-audit-tool/audit-service/internal/scraper"
)

func main() {
	cfg := config.NewConfig()

	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		log.Fatalf("failed to initialize scraper: %v", err)
	}
	defer scraperInstance.Close()

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	h := handler.NewHandler(scraperInstance)

	r.POST("/audit", h.Audit)
	r.GET("/health", h.Health)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("audit-service starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to run audit-service: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down audit-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("audit-service forced shutdown: %v", err)
	}

	log.Println("audit-service exited gracefully")
}
