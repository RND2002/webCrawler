package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type URLStatus struct {
	Url      string
	Status   int
	Title    string
	MetaData string
	Link     string
}

// function to fetch data from url
// func fetchURLData(url string, status chan URLStatus, wg *sync.WaitGroup, client *http.Client) {
// 	defer wg.Done()

// 	//client := &http.Client
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Printf("Error getting data %s", err)
// 	}

// 	//adding custom headers
// 	req.Header.Set("User-Agent", "Mozilla/5.0")
// 	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

// 	//adding cookies
// 	req.AddCookie(&http.Cookie{Name: "session", Value: "session-id"})
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error executing request:", err)
// 		return
// 	}

// 	defer res.Body.Close()
// 	//response, err := http.Get(url)

// 	// if err != nil {
// 	// 	fmt.Printf("Error fetching %s: %v\n", url, err)
// 	// 	status <- URLStatus{Url: url, Status: -1, Title: "N/A"}
// 	// 	return
// 	// }
// 	//defer response.Body.Close()

// 	// doc, err := goquery.NewDocumentFromReader(response.Body)
// 	// if err != nil {
// 	// 	fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
// 	// 	status <- URLStatus{Url: url, Status: response.StatusCode, Title: "N/A"}
// 	// 	return
// 	// }

// 	//selecting element with complex css
// 	// doc.Find("div.content > ul > li.item a").Each(func(index int, element *goquery.Selection) {
// 	// 	linkText := element.Text()
// 	// 	linkHref, _ := element.Attr("href")
// 	// 	fmt.Printf("Link: %s, URL: %s\n", linkText, linkHref)
// 	// })

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
// 		status <- URLStatus{Url: url, Status: res.StatusCode, Title: "N/A"}
// 		return
// 	}

// 	metaData := ""
// 	link := ""

// 	doc.Find("div.content > ul > li.item a").Each(func(index int, element *goquery.Selection) {
// 		metaData = element.Text()
// 		link, _ = element.Attr("href")
// 		//fmt.Printf("Link: %s, URL: %s\n", linkText, linkHref)
// 	})

// 	// Extract title from <h1> tag (fallback to <title> if not found)
// 	title := ""
// 	doc.Find("h1").Each(func(index int, element *goquery.Selection) {
// 		title = element.Text()
// 	})

// 	// Fallback to <title> tag if no <h1> found
// 	if title == "" {
// 		title = doc.Find("title").Text()
// 	}

// 	status <- URLStatus{Url: url, Status: res.StatusCode, Title: title, MetaData: metaData, Link: link}

// }
func fetchURLData(url string, status chan URLStatus, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request %s: %v\n", url, err)
		status <- URLStatus{Url: url, Status: -1, Title: "N/A"}
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.AddCookie(&http.Cookie{Name: "session", Value: "session-id"})

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request for %s: %v\n", url, err)
		status <- URLStatus{Url: url, Status: -1, Title: "N/A"}
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
		status <- URLStatus{Url: url, Status: res.StatusCode, Title: "N/A"}
		return
	}

	// Extract title from <h1> tag (fallback to <title> if not found)
	title := ""
	doc.Find("h1").Each(func(index int, element *goquery.Selection) {
		title = element.Text()
	})
	if title == "" {
		title = doc.Find("title").Text()
	}

	// Extract metadata from <meta> tags
	metaData := ""
	doc.Find("meta").Each(func(index int, element *goquery.Selection) {
		// Looking for specific meta tags, e.g., description, keywords, author
		if name, exists := element.Attr("name"); exists && (name == "description" || name == "keywords" || name == "author") {
			content, _ := element.Attr("content")
			metaData += fmt.Sprintf("%s: %s\n", name, content)
		}
		if property, exists := element.Attr("property"); exists && property == "og:title" {
			content, _ := element.Attr("content")
			metaData += fmt.Sprintf("OpenGraph Title: %s\n", content)
		}
	})

	// Extract the first link from the page, if needed
	link := ""
	doc.Find("a").First().Each(func(index int, element *goquery.Selection) {
		link, _ = element.Attr("href")
	})

	status <- URLStatus{
		Url:      url,
		Status:   res.StatusCode,
		Title:    title,
		MetaData: metaData,
		Link:     link,
	}
}

func main() {
	urls := []string{
		"https://aryandwivedi.me/",
		"https://www.google.com",
		"https://www.example.com",
		"https://www.medium.com",
		"https://www.github.com",
		"https://aryandwi.netlify.app",
	}

	//Handling headers and cookies
	client := &http.Client{}

	var wg sync.WaitGroup
	ch := make(chan URLStatus)
	start := time.Now()

	for _, url := range urls {
		fmt.Printf("Starting go routine for %s\n", url)
		// Start go routine to fetch URL data
		wg.Add(1)
		go fetchURLData(url, ch, &wg, client)
	}

	// Close the channel after all goroutines finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	dataSlice := []ScrapedData{}

	// Read results from the channel
	for status := range ch {
		dataSlice = append(dataSlice, ScrapedData{
			Title:    status.Title,
			URL:      status.Url,
			MetaData: status.MetaData,
			Link:     status.Link,
		})
		fmt.Printf("URL: %s, Status: %d, Title: %s\n", status.Url, status.Status, status.Title)
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)

	if err := saveToJSON(dataSlice, "scraped_data.json"); err != nil {
		fmt.Printf("Error saving data to JSON: %v\n", err)
	}

}

type ScrapedData struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	MetaData string `json:"metadata"`
	Link     string
}

func saveToJSON(data []ScrapedData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	return encoder.Encode(data)
}
