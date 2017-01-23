package google

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
)

var (
	reTranslations = regexp.MustCompile("^\\[{2}(.*?)\\]{2},")
	reTranslation  = regexp.MustCompile("\\[\"([^\"]+)\"")
)

type Translation struct {
	Text string
}

func (s *Session) Translate(text, sourcelang, targetlang string) (Translation, error) {
	var t Translation
	r, err := s.request("POST",
		"https://translate.googleapis.com/translate_a/single",
		bytes.NewBufferString(url.Values{
			"client": []string{"gtx"},
			"ie":     []string{"UTF-8"},
			"oe":     []string{"UTF-8"},
			"sl":     []string{sourcelang},
			"tl":     []string{targetlang},
			"dt":     []string{"t"},
			"q":      []string{text},
		}.Encode()))
	if err != nil {
		return t, err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		r.Body.Close()
		return t, err
	}
	r.Body.Close()
	trs := reTranslations.FindSubmatch(data)
	if trs == nil {
		return t, errors.New("server returned unexpected result")
	}
	strs := strings.Replace(string(trs[1]), "\\n", "\n", -1)
	tr := reTranslation.FindAllStringSubmatch(strs, -1)
	if tr == nil {
		return t, errors.New("server returned unexpected result")
	}
	b := bytes.NewBufferString("")
	for _, t := range tr {
		b.WriteString(t[1])
	}
	t.Text = b.String()
	return t, nil
}
