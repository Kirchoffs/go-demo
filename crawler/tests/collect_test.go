package tests

import (
    "crawler/collect"
    "fmt"
    "testing"
)

func TestBrowserFetch(t *testing.T) {
    url := "https://book.douban.com/subject/1007305/"
    
    var f collect.Fetcher = collect.BrowserFetch{}  
    body, err := f.Get(url)  
    if err != nil {    
        fmt.Printf("read content failed: %v", err)    
        return
    }  
    fmt.Println(string(body))
}