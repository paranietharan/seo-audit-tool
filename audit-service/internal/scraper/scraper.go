package scraper

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type ScrapeResult struct {
	URL               string   `json:"url"`
	Title             string   `json:"title"`
	MetaDescription   string   `json:"meta_description"`
	MetaKeywords      string   `json:"meta_keywords"`
	CanonicalURL      string   `json:"canonical_url"`
	RobotsMeta        string   `json:"robots_meta"`
	ViewportMeta      string   `json:"viewport_meta"`
	FaviconURL        string   `json:"favicon_url"`
	OGTitle           string   `json:"og_title"`
	OGDescription     string   `json:"og_description"`
	OGImage           string   `json:"og_image"`
	TwitterCard       string   `json:"twitter_card"`
	TwitterTitle      string   `json:"twitter_title"`
	TwitterImage      string   `json:"twitter_image"`
	H1Count           int      `json:"h1_count"`
	H2Count           int      `json:"h2_count"`
	H3Count           int      `json:"h3_count"`
	H4Count           int      `json:"h4_count"`
	H5Count           int      `json:"h5_count"`
	H6Count           int      `json:"h6_count"`
	HeadingIssues     []string `json:"heading_issues"`
	InternalLinks     []Link   `json:"internal_links"`
	ExternalLinks     []Link   `json:"external_links"`
	ImagesMissingAlt  []string `json:"images_missing_alt"`
	WordCount         int      `json:"word_count"`
	PageLoadMs        int64    `json:"page_load_ms"`
	StatusCode        int      `json:"status_code"`
	HasStructuredData bool     `json:"has_structured_data"`
}

type Scraper struct {
	browser   *rod.Browser
	semaphore chan struct{}
	mu        sync.Mutex
}

func NewScraper() (*Scraper, error) {
	binPath := os.Getenv("CHROMIUM_PATH")
	if binPath == "" {
		for _, path := range []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/bin/google-chrome",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		} {
			if _, err := os.Stat(path); err == nil {
				binPath = path
				break
			}
		}
	}

	l := launcher.New()
	if binPath != "" {
		l = l.Bin(binPath)
	}
	l = l.NoSandbox(true).
		Set("disable-dev-shm-usage").
		Set("disable-gpu").
		Headless(true)

	controlURL, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch chromium: %w", err)
	}

	browser := rod.New().ControlURL(controlURL).MustConnect()

	return &Scraper{
		browser:   browser,
		semaphore: make(chan struct{}, 4), // Throttle concurrent tabs to max 4
	}, nil
}

func (s *Scraper) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.browser != nil {
		_ = s.browser.Close()
	}
}

// validatePublicURL prevents SSRF by blocking private/internal IP addresses
func validatePublicURL(targetURL string) error {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return errors.New("empty hostname in URL")
	}

	// Check direct localhost / loopback string names
	lowerHost := strings.ToLower(hostname)
	if lowerHost == "localhost" || strings.HasSuffix(lowerHost, ".local") || strings.HasSuffix(lowerHost, ".internal") {
		return errors.New("cannot audit private or local network hosts")
	}

	// Resolve IPs for SSRF defense
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve host %s: %w", hostname, err)
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return fmt.Errorf("access to internal/private IP (%s) is forbidden", ip.String())
		}
	}

	return nil
}

func (s *Scraper) Scrape(targetURL string) (*ScrapeResult, error) {
	// 1. SSRF Protection Validation
	if err := validatePublicURL(targetURL); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	// 2. Acquire Concurrency Semaphore
	s.semaphore <- struct{}{}
	defer func() { <-s.semaphore }()

	// 3. Measure initial HTTP status and redirect chain
	statusCode := 200
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}

	httpReq, err := http.NewRequest("GET", targetURL, nil)
	if err == nil {
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		resp, err := client.Do(httpReq)
		if err == nil {
			statusCode = resp.StatusCode
			_ = resp.Body.Close()
		} else {
			statusCode = 500
		}
	}

	// 4. Headless Scraping via Rod
	start := time.Now()
	page, err := s.browser.Page(proto.TargetCreateTarget{URL: targetURL})
	if err != nil {
		return nil, fmt.Errorf("failed to open page in browser: %w", err)
	}
	defer func() {
		_ = page.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	_ = page.Context(ctx).WaitLoad()
	pageLoadMs := time.Since(start).Milliseconds()

	html, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("failed to extract page HTML: %w", err)
	}

	// 5. Parse DOM with GoQuery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML DOM: %w", err)
	}

	result := &ScrapeResult{
		URL:           targetURL,
		StatusCode:    statusCode,
		PageLoadMs:    pageLoadMs,
		HeadingIssues: make([]string, 0),
	}

	parsedBase, _ := url.Parse(targetURL)
	var baseHost string
	if parsedBase != nil {
		baseHost = parsedBase.Host
	}

	// Extract Title
	result.Title = strings.TrimSpace(doc.Find("title").First().Text())

	// Extract Metas
	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		name := strings.ToLower(strings.TrimSpace(sel.AttrOr("name", "")))
		property := strings.ToLower(strings.TrimSpace(sel.AttrOr("property", "")))
		content := strings.TrimSpace(sel.AttrOr("content", ""))

		switch name {
		case "description":
			result.MetaDescription = content
		case "keywords":
			result.MetaKeywords = content
		case "robots":
			result.RobotsMeta = content
		case "viewport":
			result.ViewportMeta = content
		case "twitter:card":
			result.TwitterCard = content
		case "twitter:title":
			result.TwitterTitle = content
		case "twitter:image":
			result.TwitterImage = resolveURL(content, parsedBase)
		}

		switch property {
		case "og:title":
			result.OGTitle = content
		case "og:description":
			result.OGDescription = content
		case "og:image":
			result.OGImage = resolveURL(content, parsedBase)
		}
	})

	// Extract Favicon
	doc.Find("link[rel*='icon']").Each(func(i int, sel *goquery.Selection) {
		if result.FaviconURL == "" {
			if href, exists := sel.Attr("href"); exists {
				result.FaviconURL = resolveURL(href, parsedBase)
			}
		}
	})

	// Extract Canonical URL (Resolved)
	if canonical, exists := doc.Find("link[rel='canonical']").Attr("href"); exists {
		result.CanonicalURL = resolveURL(canonical, parsedBase)
	}

	// Extract Headings
	result.H1Count = doc.Find("h1").Length()
	result.H2Count = doc.Find("h2").Length()
	result.H3Count = doc.Find("h3").Length()
	result.H4Count = doc.Find("h4").Length()
	result.H5Count = doc.Find("h5").Length()
	result.H6Count = doc.Find("h6").Length()

	if result.H1Count == 0 {
		result.HeadingIssues = append(result.HeadingIssues, "Missing primary <h1> heading tag")
	} else if result.H1Count > 1 {
		result.HeadingIssues = append(result.HeadingIssues, fmt.Sprintf("Multiple (%d) <h1> tags detected", result.H1Count))
	}

	// Extract Structured Data
	result.HasStructuredData = doc.Find("script[type='application/ld+json']").Length() > 0

	// Extract Links
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "#") {
			return
		}

		text := strings.TrimSpace(sel.Text())
		resolved := resolveURL(href, parsedBase)
		linkItem := Link{Href: resolved, Text: text}

		if isInternalLink(resolved, baseHost) {
			result.InternalLinks = append(result.InternalLinks, linkItem)
		} else {
			result.ExternalLinks = append(result.ExternalLinks, linkItem)
		}
	})

	// Extract Images Missing Alt
	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		alt, altExists := sel.Attr("alt")
		src, srcExists := sel.Attr("src")
		if srcExists {
			cleanSrc := resolveURL(src, parsedBase)
			if !altExists || strings.TrimSpace(alt) == "" {
				result.ImagesMissingAlt = append(result.ImagesMissingAlt, cleanSrc)
			}
		}
	})

	// Accurate Word Count (Stripping non-content nodes first)
	cleanDoc := goquery.CloneDocument(doc)
	cleanDoc.Find("script, style, noscript, svg, iframe, nav, footer, header").Remove()
	bodyText := cleanDoc.Find("body").Text()
	words := strings.Fields(bodyText)
	result.WordCount = len(words)

	return result, nil
}

func resolveURL(href string, base *url.URL) string {
	href = strings.TrimSpace(href)
	if href == "" || base == nil {
		return href
	}
	parsedHref, err := url.Parse(href)
	if err != nil {
		return href
	}
	if parsedHref.IsAbs() {
		return href
	}
	return base.ResolveReference(parsedHref).String()
}

func isInternalLink(linkStr, baseHost string) bool {
	u, err := url.Parse(linkStr)
	if err != nil {
		return false
	}
	return u.Host == "" || strings.EqualFold(u.Host, baseHost)
}
