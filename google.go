package google

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Session struct {
	cli http.Client
}

func (s *Session) Init() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.cli.Jar = jar
	s.cli.Timeout = time.Second * 10
	return nil
}

func (s *Session) IsInitialized() bool {
	return s.cli.Jar != nil
}

func (s *Session) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "google api")
	return s.cli.Do(req)
}
