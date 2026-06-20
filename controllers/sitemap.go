package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SitemapController struct{}

func NewSitemapController() *SitemapController {
	return &SitemapController{}
}

type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type SitemapURLSet struct {
	XMLName string       `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

func (sc *SitemapController) GenerateSitemap(c *gin.Context) {
	baseURL := "https://www.glynsgardeningservices.com"
	currentDate := time.Now().Format("2006-01-02")

	urls := []SitemapURL{
		{
			Loc:        baseURL + "/",
			LastMod:    currentDate,
			ChangeFreq: "weekly",
			Priority:   "1.0",
		},
	}

	sitemap := SitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	c.Header("Content-Type", "application/xml")
	c.XML(http.StatusOK, sitemap)
}
