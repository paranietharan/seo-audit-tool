package analyzer

import (
	"seo-audit-tool/internal/crawler"
	"seo-audit-tool/internal/storage"
	"strings"
)

type TechnicalAnalyzer struct{}

func NewTechnicalAnalyzer() *TechnicalAnalyzer {
	return &TechnicalAnalyzer{}
}

func (t *TechnicalAnalyzer) Analyze(page *crawler.Page, result *storage.PageResult) {
	score := 0

	if strings.HasPrefix(page.URL, "https://") {
		result.IsHTTPS = true
		score += 20
	}

	if page.LoadTime.Milliseconds() < 3000 {
		score += 30
	} else if page.LoadTime.Milliseconds() < 5000 {
		score += 15
	}

	if result.StatusCode == 200 {
		score += 25
	}

	if strings.Contains(page.HTMLContent, "viewport") {
		result.IsMobileFriendly = true
		score += 25
	}

	result.TechnicalScore = score
}
