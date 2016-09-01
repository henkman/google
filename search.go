package google

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	reImage = regexp.MustCompile("imgurl=([^&]+)&")
)

type SearchResult struct {
	Title   string
	URL     string
	Content string
}

type ImageResult struct {
	URL string
}

type ImageType string

const (
	ImageType_Any      ImageType = ""
	ImageType_Animated           = "animated"
	ImageType_Face               = "face"
	ImageType_Clipart            = "clipart"
	ImageType_Photo              = "photo"
	ImageType_Lineart            = "lineart"
)

func (c *Client) Images(tld, query, lang string, safe bool, t ImageType, start, num uint) ([]ImageResult, error) {
	var doc *goquery.Document
	{
		ps := url.Values{
			"hl":    []string{lang},
			"q":     []string{query},
			"btnG":  []string{"Google+Search"},
			"start": []string{fmt.Sprint(start)},
			"num":   []string{fmt.Sprint(num)},
			"tbm":   []string{"isch"},
		}
		if t != ImageType_Any {
			ps.Set("tbs", "itp:"+string(t))
		}
		if safe {
			ps.Set("safe", "on")
		} else {
			ps.Set("safe", "off")
		}
		url_str := "https://www.google." + tld + "/search?" + ps.Encode()
		r, err := c.request("GET", url_str, nil)
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
	results := make([]ImageResult, 0, num)
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
	if uint(len(results)) > num {
		results = results[:num]
	}
	return results, nil
}

func (c *Client) Search(tld, query, lang string, safe bool, start, num uint) ([]SearchResult, error) {
	var doc *goquery.Document
	{
		ps := url.Values{
			"hl":    []string{lang},
			"q":     []string{query},
			"btnG":  []string{"Google+Search"},
			"start": []string{fmt.Sprint(start)},
			"num":   []string{fmt.Sprint(num)},
		}
		if safe {
			ps.Set("safe", "on")
		} else {
			ps.Set("safe", "off")
		}
		url_str := "https://www.google." + tld + "/search?" + ps.Encode()
		r, err := c.request("GET", url_str, nil)
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
	results := make([]SearchResult, 0, num)
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
	if uint(len(results)) > num {
		results = results[:num]
	}
	return results, nil
}
