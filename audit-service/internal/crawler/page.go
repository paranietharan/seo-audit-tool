package crawler

import "time"

type Page struct {
	URL         string
	Content     string
	Title       string
	StatusCode  int
	LoadTime    time.Duration
	Links       []string
	Images      []Image
	Headers     map[string]string
	HTMLContent string
}

type Image struct {
	Src string
	Alt string
}
