package google

import (
	"io"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	client http.Client
}

func (c *Client) Init() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	c.client.Jar = jar
	return nil
}

func (c *Client) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36")
	return c.client.Do(req)
}
