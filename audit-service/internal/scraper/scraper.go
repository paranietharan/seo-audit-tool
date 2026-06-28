package scraper

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	URL                string   `json:"url"`
	Title              string   `json:"title"`
	MetaDescription    string   `json:"meta_description"`
	MetaKeywords       string   `json:"meta_keywords"`
	CanonicalURL       string   `json:"canonical_url"`
	RobotsMeta         string   `json:"robots_meta"`
	OGTitle            string   `json:"og_title"`
	OGDescription      string   `json:"og_description"`
	OGImage            string   `json:"og_image"`
	H1Count            int      `json:"h1_count"`
	H2Count            int      `json:"h2_count"`
	H3Count            int      `json:"h3_count"`
	H4Count            int      `json:"h4_count"`
	H5Count            int      `json:"h5_count"`
	H6Count            int      `json:"h6_count"`
	InternalLinks      []Link   `json:"internal_links"`
	ExternalLinks      []Link   `json:"external_links"`
	ImagesMissingAlt   []string `json:"images_missing_alt"`
	WordCount          int      `json:"word_count"`
	PageLoadMs         int64    `json:"page_load_ms"`
	StatusCode         int      `json:"status_code"`
	HasStructuredData  bool     `json:"has_structured_data"`
}

type Scraper struct {
	browser *rod.Browser
}

func NewScraper() (*Scraper, error) {
	binPath := os.Getenv("CHROMIUM_PATH")
	if binPath == "" {
		for _, path := range []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/bin/google-chrome",
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
		Headless(true)

	controlURL, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(controlURL).MustConnect()
	return &Scraper{browser: browser}, nil
}

func (s *Scraper) Close() {
	if s.browser != nil {
		_ = s.browser.Close()
	}
}

func (s *Scraper) Scrape(targetURL string) (*ScrapeResult, error) {
	// 1. Get HTTP Status Code
	statusCode := 200
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(targetURL)
	if err == nil {
		statusCode = resp.StatusCode
		resp.Body.Close()
	} else {
		statusCode = 500
	}

	// 2. Perform Headless Scraping with Rod
	start := time.Now()
	page, err := s.browser.Page(proto.TargetCreateTarget{URL: targetURL})
	if err != nil {
		return nil, err
	}
	defer page.MustClose()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = page.Context(ctx).WaitLoad()
	if err != nil {
		return nil, err
	}

	pageLoadMs := time.Since(start).Milliseconds()

	html, err := page.HTML()
	if err != nil {
		return nil, err
	}

	// 3. Parse Rendered DOM with Goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	result := &ScrapeResult{
		URL:        targetURL,
		StatusCode: statusCode,
		PageLoadMs: pageLoadMs,
	}

	// Extract Title
	result.Title = strings.TrimSpace(doc.Find("title").First().Text())

	// Extract Metas
	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		name, _ := sel.Attr("name")
		property, _ := sel.Attr("property")
		content, _ := sel.Attr("content")

		name = strings.ToLower(name)
		property = strings.ToLower(property)

		switch name {
		case "description":
			result.MetaDescription = content
		case "keywords":
			result.MetaKeywords = content
		case "robots":
			result.RobotsMeta = content
		}

		switch property {
		case "og:title":
			result.OGTitle = content
		case "og:description":
			result.OGDescription = content
		case "og:image":
			result.OGImage = content
		}
	})

	// Extract Canonical Link
	canonical, exists := doc.Find("link[rel='canonical']").Attr("href")
	if exists {
		result.CanonicalURL = canonical
	}

	// Extract Headings
	result.H1Count = doc.Find("h1").Length()
	result.H2Count = doc.Find("h2").Length()
	result.H3Count = doc.Find("h3").Length()
	result.H4Count = doc.Find("h4").Length()
	result.H5Count = doc.Find("h5").Length()
	result.H6Count = doc.Find("h6").Length()

	// Extract Structured Data Presence
	result.HasStructuredData = doc.Find("script[type='application/ld+json']").Length() > 0

	// Parse Base Host for Link Classification
	parsedBase, err := url.Parse(targetURL)
	var baseHost string
	if err == nil {
		baseHost = parsedBase.Host
	}

	// Extract Links
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
			return
		}

		text := strings.TrimSpace(sel.Text())

		// Resolve relative paths
		resolved := href
		if parsedHref, err := url.Parse(href); err == nil && !parsedHref.IsAbs() {
			resolved = parsedBase.ResolveReference(parsedHref).String()
		}

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
			if !altExists || strings.TrimSpace(alt) == "" {
				result.ImagesMissingAlt = append(result.ImagesMissingAlt, src)
			}
		}
	})

	// Extract Word Count (body text)
	bodyText := doc.Find("body").Text()
	words := strings.Fields(bodyText)
	result.WordCount = len(words)

	return result, nil
}

func isInternalLink(linkStr, baseHost string) bool {
	u, err := url.Parse(linkStr)
	if err != nil {
		return false
	}
	return u.Host == "" || u.Host == baseHost
}
