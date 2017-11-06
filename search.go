package google

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	reImage = regexp.MustCompile("imgurl=(.*?)&amp;")
	reWeb   = regexp.MustCompile("/url\\?q=(.*?)&")
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
	ImageType_Animated ImageType = "animated"
	ImageType_Face     ImageType = "face"
	ImageType_Clipart  ImageType = "clipart"
	ImageType_Photo    ImageType = "photo"
	ImageType_Lineart  ImageType = "lineart"
)

func (s *Session) Images(tld, query, lang string,
	safe bool, t ImageType, start, num uint) ([]ImageResult, error) {
	ps := url.Values{
		"hl":    []string{lang},
		"q":     []string{query},
		"start": []string{fmt.Sprint(start)},
		"ijn":   []string{fmt.Sprint(num)},
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
	r, err := s.request("GET", url_str, nil)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	m := reImage.FindAllStringSubmatch(string(raw), -1)
	if m == nil {
		return []ImageResult{}, nil
	}
	results := make([]ImageResult, 0, num)
	for _, img := range m {
		results = append(results, ImageResult{
			URL: img[1],
		})
	}
	if uint(len(results)) > num {
		results = results[:num]
	}
	return results, nil
}

func (s *Session) Search(tld, query, lang string,
	safe bool, start, num uint) ([]SearchResult, error) {
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
	r, err := s.request("GET", url_str, nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		return nil, err
	}
	elems := doc.Find("#search .g")
	if elems.Length() == 0 {
		return []SearchResult{}, nil
	}
	results := make([]SearchResult, 0, num)
	elems.Each(func(i int, s *goquery.Selection) {
		a := s.Find(".r a")
		if h, ok := a.Attr("href"); ok {
			m := reWeb.FindStringSubmatch(h)
			if m != nil {
				results = append(results, SearchResult{
					Title:   a.Text(),
					URL:     m[1],
					Content: s.Find(".st").Text(),
				})
			}
		}
	})
	if uint(len(results)) > num {
		results = results[:num]
	}
	return results, nil
}
