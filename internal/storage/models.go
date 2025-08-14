package storage

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Audit struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	URL           string    `json:"url"`
	Status        string    `json:"status"` // processing, completed, failed
	Error         string    `json:"error,omitempty"`
	PagesAnalyzed int       `json:"pages_analyzed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	PageResults []PageResult `gorm:"foreignKey:AuditID" json:"page_results,omitempty"`
}

type PageResult struct {
	ID      string `gorm:"primaryKey" json:"id"`
	AuditID string `json:"audit_id"`
	URL     string `json:"url"`

	// Technical SEO
	LoadTime         int  `json:"load_time"`
	StatusCode       int  `json:"status_code"`
	IsHTTPS          bool `json:"is_https"`
	HasSitemap       bool `json:"has_sitemap"`
	HasRobotsTxt     bool `json:"has_robots_txt"`
	IsMobileFriendly bool `json:"is_mobile_friendly"`

	// On-Page SEO
	Title           string `json:"title"`
	MetaDescription string `json:"meta_description"`
	H1Count         int    `json:"h1_count"`
	H2Count         int    `json:"h2_count"`
	H3Count         int    `json:"h3_count"`
	ImagesCount     int    `json:"images_count"`
	ImagesWithAlt   int    `json:"images_with_alt"`
	InternalLinks   int    `json:"internal_links"`
	ExternalLinks   int    `json:"external_links"`

	// Content Analysis
	WordCount    int     `json:"word_count"`
	ReadingTime  int     `json:"reading_time"`
	ContentScore float64 `json:"content_score"`

	// SEO Score
	TechnicalScore int `json:"technical_score"`
	OnPageScore    int `json:"onpage_score"`
	ContentScore2  int `json:"content_score2"`
	OverallScore   int `json:"overall_score"`

	CreatedAt time.Time `json:"created_at"`
}

func (a *Audit) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New().String()
	return nil
}

func (pr *PageResult) BeforeCreate(tx *gorm.DB) error {
	pr.ID = uuid.New().String()
	return nil
}
