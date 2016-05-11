package google

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	reImage = regexp.MustCompile("imgurl=([^&]+)&")
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

type ImageResult struct {
	URL string
}

func (c *Client) Init(tld string) error {
	r, err := c.get("https://www.google." + tld)
	if err != nil {
		return err
	}
	r.Body.Close()
	return nil
}

func (c *Client) get(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36")
	return c.client.Do(req)
}

func (c *Client) Images(tld, query, lang string, safe bool, count int) ([]ImageResult, error) {
	var doc *goquery.Document
	{
		ps := url.Values{
			"hl":   []string{lang},
			"q":    []string{query},
			"btnG": []string{"Google+Search"},
			"tbm":  []string{"isch"},
		}
		if !safe {
			ps.Set("safe", "off")
		}
		url_str := "https://www.google." + tld + "/search?" + ps.Encode()
		r, err := c.get(url_str)
		if err != nil {
			return nil, err
		}
		doc, err = goquery.NewDocumentFromResponse(r)
		if err != nil {
			return nil, err
		}
		// goquery.NewDocumentFromResponse closes body
		// r.Body.Close()
	}
	elems := doc.Find(".ivg-i a")
	if elems.Length() == 0 {
		return []ImageResult{}, nil
	}
	results := make([]ImageResult, 0, count)
	elems.Each(func(i int, s *goquery.Selection) {
		h, ok := s.Attr("href")
		if !ok {
			return
		}
		m := reImage.FindStringSubmatch(h)
		if m == nil {
			return
		}
		results = append(results, ImageResult{
			URL: m[1],
		})
	})
	if len(results) > count {
		results = results[:count]
	}
	return results, nil
}

func (c *Client) Search(tld, query, lang string, safe bool, count int) ([]SearchResult, error) {
	var doc *goquery.Document
	{
		ps := url.Values{
			"hl":   []string{lang},
			"q":    []string{query},
			"btnG": []string{"Google+Search"},
		}
		if !safe {
			ps.Set("safe", "off")
		}
		url_str := "https://www.google." + tld + "/search?" + ps.Encode()
		r, err := c.get(url_str)
		if err != nil {
			return nil, err
		}
		doc, err = goquery.NewDocumentFromResponse(r)
		if err != nil {
			return nil, err
		}
		// goquery.NewDocumentFromResponse closes body
		// r.Body.Close()
	}
	elems := doc.Find(".rc")
	if elems.Length() == 0 {
		return []SearchResult{}, nil
	}
	results := make([]SearchResult, 0, count)
	elems.Each(func(i int, s *goquery.Selection) {
		a := s.Find(".r a")
		h, ok := a.Attr("href")
		if !ok {
			return
		}
		results = append(results, SearchResult{
			Title:   a.Text(),
			URL:     h,
			Content: s.Find(".st").Text(),
		})
	})
	if len(results) > count {
		results = results[:count]
	}
	return results, nil
}
