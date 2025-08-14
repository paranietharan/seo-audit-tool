package crawler

import (
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Crawler struct {
	browser *rod.Browser
}

func NewCrawler() *Crawler {
	l := launcher.New().Headless(true)
	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()

	return &Crawler{browser: browser}
}

func (c *Crawler) CrawlSite(startURL string, maxPages int) ([]*Page, error) {
	visited := make(map[string]bool)
	queue := []string{startURL}
	pages := make([]*Page, 0)

	baseURL, err := url.Parse(startURL)
	if err != nil {
		return nil, err
	}

	for len(queue) > 0 && len(pages) < maxPages {
		currentURL := queue[0]
		queue = queue[1:]

		if visited[currentURL] {
			continue
		}

		page, err := c.CrawlPage(currentURL)
		if err != nil {
			continue
		}

		visited[currentURL] = true
		pages = append(pages, page)

		// Add internal links to queue
		for _, link := range page.Links {
			if c.isInternalLink(link, baseURL) && !visited[link] {
				queue = append(queue, link)
			}
		}

		// Be respectful - add delay
		time.Sleep(1 * time.Second)
	}

	return pages, nil
}

func (c *Crawler) CrawlPage(pageURL string) (*Page, error) {
	start := time.Now()

	rodPage := c.browser.MustPage(pageURL)
	defer rodPage.MustClose()

	// Wait for page to load
	rodPage.MustWaitLoad()

	loadTime := time.Since(start)

	page := &Page{
		URL:      pageURL,
		LoadTime: loadTime,
	}

	// Get basic info
	page.Title = rodPage.MustInfo().Title
	page.HTMLContent = rodPage.MustHTML()

	// Get status code from response
	page.StatusCode = 200 // Rod doesn't easily expose this, assume 200 if page loads

	// Extract links
	links := rodPage.MustElements("a[href]")
	for _, link := range links {
		href := link.MustAttribute("href")
		if href != nil {
			page.Links = append(page.Links, c.resolveURL(*href, pageURL))
		}
	}

	// Extract images
	images := rodPage.MustElements("img")
	for _, img := range images {
		src := img.MustAttribute("src")
		alt := img.MustAttribute("alt")

		if src != nil {
			image := Image{Src: *src}
			if alt != nil {
				image.Alt = *alt
			}
			page.Images = append(page.Images, image)
		}
	}

	// Get text content for analysis
	page.Content = rodPage.MustElement("body").MustText()

	return page, nil
}

func (c *Crawler) isInternalLink(link string, baseURL *url.URL) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	return linkURL.Host == "" || linkURL.Host == baseURL.Host
}

func (c *Crawler) resolveURL(href, baseURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	link, err := url.Parse(href)
	if err != nil {
		return href
	}

	resolved := base.ResolveReference(link)
	return resolved.String()
}

func (c *Crawler) Close() {
	if c.browser != nil {
		c.browser.MustClose()
	}
}
