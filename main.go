package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Make HTTP request
	res, err := http.Get("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Extract title
	title := doc.Find("title").Text()

	// Print title
	fmt.Println(title)

	// Take screenshot
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(`https://google.com`, 100, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(contentSize.Width), int64(contentSize.Height)
			if err := chromedp.EmulateViewport(width, height).Do(ctx); err != nil {
				return err
			}
			return nil
		}),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  1920,
					Height: 1080,
					Scale:  1,
				}).
				WithFormat(page.CaptureScreenshotFormatPng).
				Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
