package google

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	reUrl = regexp.MustCompile("/url\\?q=([^&]+)&")
)

type Client struct {
	client http.Client
}

func New() (*Client, error) {
	var client http.Client
	{
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		client = http.Client{Jar: jar}
	}
	c := new(Client)
	c.client = client
	return c, nil
}

type SearchResult struct {
	Title   string
	URL     string
	Content string
}

func (c *Client) Init(tld string) error {
	r, err := c.get(fmt.Sprintf("https://www.google.%s", tld))
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) get(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0)")
	return c.client.Do(req)
}

func (c *Client) Search(tld, query, language string, count int) ([]SearchResult, error) {
	var doc *goquery.Document
	{
		url_str := fmt.Sprintf(
			"https://www.google.%s/search?hl=%s&q=%s&btnG=Google+Search&safe=off",
			tld,
			language,
			url.QueryEscape(query),
		)
		r, err := c.get(url_str)
		if err != nil {
			return nil, err
		}
		doc, err = goquery.NewDocumentFromResponse(r)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
	}
	results := make([]SearchResult, 0, count)
	doc.Find(".g .s").Slice(0, count).Each(func(i int, s *goquery.Selection) {
		el := s.Parent()
		a := el.Find(".r a")
		h, ok := a.Attr("href")
		if !ok {
			return
		}
		m := reUrl.FindStringSubmatch(h)
		u, err := url.QueryUnescape(m[1])
		if err != nil {
			return
		}
		results = append(results, SearchResult{
			Title:   a.Text(),
			URL:     u,
			Content: s.Find(".st").Text(),
		})
	})
	return results, nil
}
