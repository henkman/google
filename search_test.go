package google

import (
	"fmt"
	"testing"
)

const (
	TLD = "de"
)

func TestSearch(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Log(err)
		t.Fail()
	}
	rs, err := s.Search(TLD, "kittens", "de", true, 0, 5)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, r := range rs {
		fmt.Println(r.Title, r.URL)
	}
}

func TestImages(t *testing.T) {
	var s Session
	if err := s.Init(); err != nil {
		t.Log(err)
		t.Fail()
	}
	rs, err := s.Images(TLD, "kittens", "de", true, ImageType_Any, 0, 5)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, r := range rs {
		fmt.Println(r.URL)
	}
}
