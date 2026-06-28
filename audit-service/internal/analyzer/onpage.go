package analyzer

import (
	"github.com/PuerkitoBio/goquery"
	"seo-audit-tool/internal/crawler"
	"seo-audit-tool/internal/storage"
	"strings"
)

type OnPageAnalyzer struct{}

func NewOnPageAnalyzer() *OnPageAnalyzer {
	return &OnPageAnalyzer{}
}

func (o *OnPageAnalyzer) Analyze(page *crawler.Page, result *storage.PageResult) {
	score := 0

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page.HTMLContent))
	if err != nil {
		result.OnPageScore = 0
		return
	}

	result.Title = page.Title
	if len(result.Title) > 0 && len(result.Title) <= 60 {
		score += 20
	}

	metaDesc, exists := doc.Find("meta[name='description']").Attr("content")
	if exists {
		result.MetaDescription = metaDesc
		if len(metaDesc) > 0 && len(metaDesc) <= 160 {
			score += 20
		}
	}

	result.H1Count = doc.Find("h1").Length()
	result.H2Count = doc.Find("h2").Length()
	result.H3Count = doc.Find("h3").Length()

	if result.H1Count == 1 {
		score += 15
	}
	if result.H2Count > 0 {
		score += 10
	}

	// Image analysis
	result.ImagesCount = len(page.Images)
	for _, img := range page.Images {
		if img.Alt != "" {
			result.ImagesWithAlt++
		}
	}

	if result.ImagesCount > 0 {
		altPercentage := float64(result.ImagesWithAlt) / float64(result.ImagesCount)
		if altPercentage >= 0.8 {
			score += 15
		} else if altPercentage >= 0.5 {
			score += 8
		}
	}

	for _, link := range page.Links {
		if strings.HasPrefix(link, "http") {
			result.ExternalLinks++
		} else {
			result.InternalLinks++
		}
	}

	if result.InternalLinks > 0 {
		score += 10
	}

	if result.H1Count > 0 && result.H2Count > 0 {
		score += 10
	}

	result.OnPageScore = score
}
