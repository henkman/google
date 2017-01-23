package google

import (
	"fmt"
	"testing"
)

var (
	c Session
)

func TestTranslateInit(t *testing.T) {
	if err := s.Init(); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestTranslate(t *testing.T) {
	tr, err := s.Translate(`Dies ist ein etwas längerer Text.
Hoffentlich frisst der Übersetzer ihn. Der Hund springt manchmal
über Hürden.`, "auto", "en")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(tr.Text)
}
