package analyzer

import (
	"seo-audit-tool/internal/crawler"
	"seo-audit-tool/internal/storage"
	"strings"
)

type ContentAnalyzer struct{}

func NewContentAnalyzer() *ContentAnalyzer {
	return &ContentAnalyzer{}
}

func (c *ContentAnalyzer) Analyze(page *crawler.Page, result *storage.PageResult) {
	score := 0

	words := strings.Fields(page.Content)
	result.WordCount = len(words)

	result.ReadingTime = result.WordCount / 200
	if result.ReadingTime == 0 {
		result.ReadingTime = 1
	}

	if result.WordCount >= 300 {
		score += 30
	} else if result.WordCount >= 150 {
		score += 15
	}

	content := strings.ToLower(page.Content)

	titleLower := strings.ToLower(result.Title)
	if titleLower != "" && strings.Contains(content, titleLower) {
		score += 15
	}

	sentences := strings.Split(page.Content, ".")
	if len(sentences) > 5 {
		score += 15
	}

	currentYear := "2024"
	if strings.Contains(content, currentYear) {
		score += 10
	}

	if len(sentences) > 0 {
		avgWordsPerSentence := result.WordCount / len(sentences)
		if avgWordsPerSentence >= 10 && avgWordsPerSentence <= 20 {
			score += 15
		}
	}

	if score > 100 {
		score = 100
	}

	result.ContentScore = float64(score)
	result.ContentScore2 = score
}
