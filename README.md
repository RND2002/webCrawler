# Go Web Crawler

This is a simple Go program that fetches and parses data from a list of URLs concurrently. It demonstrates Go's concurrency model using goroutines and channels. The program also utilizes custom HTTP headers and cookies when making requests. HTML parsing is done using the `goquery` package to extract elements like the page title and links from the response.

## Features

- **Concurrency**: The program uses Go's goroutines and the `sync.WaitGroup` to fetch data from multiple URLs concurrently, making the program highly efficient.
- **Channels**: Go's channels are used to pass the fetched results (URL status and title) between goroutines and the main thread.
- **Custom Headers and Cookies**: HTTP requests include custom headers like `User-Agent` and `Accept-Language`, as well as cookies.
- **HTML Parsing**: The program uses the `goquery` package to extract elements from the HTML response, such as `<h1>` tags and links.
- **Graceful Handling of Errors**: Proper error handling is implemented to ensure the program continues even if some requests fail.

## Packages Used

- [Goquery](https://github.com/PuerkitoBio/goquery): A Go package for easily parsing and manipulating HTML documents, similar to jQuery.
- `net/http`: Standard Go package for making HTTP requests and handling responses.
- `sync`: Standard Go package used for synchronizing concurrent goroutines using `WaitGroup`.
- `time`: Standard Go package to measure the execution time of the program.

## How It Works

1. The list of URLs is defined in the `urls` slice.
2. A new goroutine is started for each URL using the `fetchURLData` function.
3. HTTP requests are made using the `http.Client`, with custom headers and cookies added to the request.
4. The `goquery` package parses the HTML response and extracts the page title (from `<h1>` or `<title>` tags).
5. The fetched results (URL, status code, and title) are sent to a channel.
6. Once all goroutines complete their execution, the results are read from the channel and printed to the console.
7. The total elapsed time is displayed at the end.

## Setup Instructions

### Prerequisites

- Go (version 1.16 or above)
- Internet connection (to fetch the URLs)

### Steps to Run

1. Clone this repository.
   
   ```bash
   git clone https://github.com/RND2002/webCrawler.git
   cd go-url-fetcher

   go get github.com/PuerkitoBio/goquery

   go run main.go

