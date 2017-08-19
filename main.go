package main

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func Crawl(url string, depth int, m *message) {
	defer func() { m.quit <- 0 }()

	// WebページからURLを取得
	urls, err := Fetch(url)

	// 結果送信
	m.res <- &respons{
		url: url,
		err: err,
	}

	if err == nil {
		for _, url := range urls {
			// 新しいリクエスト送信
			m.req <- &request{
				url:   url,
				depth: depth - 1,
			}
		}
	}
}

func Fetch(u string) (urls []string, err error) {
	baseUrl, err := url.Parse(u)
	if err != nil {
		return
	}

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	urls = make([]string, 0)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			reqUrl, err := baseUrl.Parse(href)
			if err == nil {
				urls = append(urls, reqUrl.String())
			}
		}
	})

	return
}

func main() {

}
