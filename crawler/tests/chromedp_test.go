package tests

import (
    "context"
    "log"
    "testing"
    "time"

    "github.com/chromedp/chromedp"
)

func TestChromeDP(t *testing.T) {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 15 * time.Second)
    defer cancel()
    
    var example string
    err := chromedp.Run(
        ctx,
        chromedp.Navigate(`https://pkg.go.dev/time`),   
        chromedp.WaitVisible(`body > footer`),    
        chromedp.Click(`#example-After`, chromedp.NodeVisible),    
        chromedp.Value(`#example-After textarea`, &example),
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Go's time.After example:\\n%s", example)
}