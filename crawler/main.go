package main

import (
    "crawler/collect"
    "fmt"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func main() {
    url := "https://ssr1.scrape.center/"
    var f collect.Fetcher = collect.BaseFetch{}
    body, err := f.Get(url)

    if err != nil {
        fmt.Printf("read content failed: %v", err)
        return
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
    if err != nil {
        fmt.Printf("read content failed: %v", err)
    }

    doc.Find("div.small_toplink__GmZhY a h2").Each(func(i int, s *goquery.Selection) {
        title := s.Text()
        fmt.Printf("title %d: %s\n", i, title)
    })
}
