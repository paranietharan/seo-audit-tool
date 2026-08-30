package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	auditProxy  *httputil.ReverseProxy
	reportProxy *httputil.ReverseProxy
}

func NewProxy(auditURL, reportURL string) (*Proxy, error) {
	aURL, err := url.Parse(auditURL)
	if err != nil {
		return nil, err
	}
	rURL, err := url.Parse(reportURL)
	if err != nil {
		return nil, err
	}

	auditProxy := httputil.NewSingleHostReverseProxy(aURL)
	reportProxy := httputil.NewSingleHostReverseProxy(rURL)

	// Deduplicate backend CORS headers so gateway controls CORS cleanly
	modifyResponse := func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")
		return nil
	}
	auditProxy.ModifyResponse = modifyResponse
	reportProxy.ModifyResponse = modifyResponse

	// Upstream error handlers returning JSON
	auditProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("audit proxy upstream error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"audit_service_unavailable","detail":"Unable to reach audit service backend"}`))
	}

	reportProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("report proxy upstream error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"report_service_unavailable","detail":"Unable to reach report service backend"}`))
	}

	// Rewrite paths and update Host header
	origAuditDirector := auditProxy.Director
	auditProxy.Director = func(req *http.Request) {
		origAuditDirector(req)
		req.Host = aURL.Host
		req.URL.Path = "/audit"
	}

	origReportDirector := reportProxy.Director
	reportProxy.Director = func(req *http.Request) {
		origReportDirector(req)
		req.Host = rURL.Host
		req.URL.Path = "/report"
	}

	return &Proxy{
		auditProxy:  auditProxy,
		reportProxy: reportProxy,
	}, nil
}

func (p *Proxy) AuditHandler(c *gin.Context) {
	p.auditProxy.ServeHTTP(c.Writer, c.Request)
}

func (p *Proxy) ReportHandler(c *gin.Context) {
	p.reportProxy.ServeHTTP(c.Writer, c.Request)
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
