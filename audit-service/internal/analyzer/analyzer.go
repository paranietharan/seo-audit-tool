package analyzer

import (
	"seo-audit-tool/internal/crawler"
	"seo-audit-tool/internal/storage"
)

type Analyzer struct {
	technical *TechnicalAnalyzer
	onpage    *OnPageAnalyzer
	content   *ContentAnalyzer
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		technical: NewTechnicalAnalyzer(),
		onpage:    NewOnPageAnalyzer(),
		content:   NewContentAnalyzer(),
	}
}

func (a *Analyzer) AnalyzePage(page *crawler.Page) *storage.PageResult {
	result := &storage.PageResult{
		URL:        page.URL,
		StatusCode: page.StatusCode,
		LoadTime:   int(page.LoadTime.Milliseconds()),
	}

	a.technical.Analyze(page, result)

	a.onpage.Analyze(page, result)

	a.content.Analyze(page, result)

	result.OverallScore = (result.TechnicalScore + result.OnPageScore + result.ContentScore2) / 3

	return result
}
