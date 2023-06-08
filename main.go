package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type Newspaper struct {
	Name    string
	Address string
	Base    string
}

type Article struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

var newspapers = []Newspaper{
	{
		Name:    "nyp",
		Address: "https://nypost.com/business/",
		Base:    "",
	},
	{
		Name:    "fintimes",
		Address: "https://www.nytimes.com/section/business",
		Base:    "https://www.nytimes.com",
	},
	{
		Name:    "nytimes",
		Address: "https://www.ft.com/",
		Base:    "https://www.ft.com",
	},
	{
		Name:    "yahooFinance",
		Address: "https://finance.yahoo.com/?guccounter=1&guce_referrer=aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS8&guce_referrer_sig=AQAAAJQBSP36mCJgLU9y681i3zi9_qWXiI4KUoB8b8LQt_eSJUyABBVe4A9w3v3NVufc2-5F31sEJvCEjdzzO8x7KXHPDT5_k3x3l3ZNciY_wGz1k2gK6gPJZp3EH0cYbRfEhZHpkfhInEGCUIl1vvfK8kMoaebH1ihI0qaGtSnYPBIW",
		Base:    "https://finance.yahoo.com",
	},
	{
		Name:    "cnbc",
		Address: "https://www.cnbc.com/finance/",
		Base:    "",
	},
	// ... add other newspapers
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to my Tech Company News API")
	})

	router.GET("/news", func(c *gin.Context) {
		articles := fetchArticles()
		c.JSON(http.StatusOK, articles)
	})

	router.GET("/news/:newspaperId", func(c *gin.Context) {
		newspaperID := c.Param("newspaperId")

		newspaper := findNewspaperByID(newspaperID)
		if newspaper == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Newspaper not found"})
			return
		}

		resp, err := http.Get(newspaper.Address)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch newspaper"})
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML"})
			return
		}

		specificArticles := []Article{}
		doc.Find("a:contains('Amazon')").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			url, _ := s.Attr("href")

			specificArticles = append(specificArticles, Article{
				Title:  title,
				URL:    newspaper.Base + url,
				Source: newspaper.Name,
			})
		})

		c.JSON(http.StatusOK, specificArticles)
	})

	const port = 8000
	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("Server running on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

// fetchArticles fetches articles from the newspapers and returns a slice of Article.
func fetchArticles() []Article {
	articles := []Article{}

	for _, newspaper := range newspapers {
		resp, err := http.Get(newspaper.Address)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		doc.Find("a:contains('Amazon')").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			url, _ := s.Attr("href")

			articles = append(articles, Article{
				Title:  title,
				URL:    newspaper.Base + url,
				Source: newspaper.Name,
			})
		})
	}

	return articles
}

// findNewspaperByID finds and returns a pointer to the Newspaper with the given ID.
func findNewspaperByID(id string) *Newspaper {
	for _, newspaper := range newspapers {
		if newspaper.Name == id {
			return &newspaper
		}
	}
	return nil
}
