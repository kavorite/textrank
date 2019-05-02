package textrank

import (
	"testing"
	"io/ioutil"
	"strings"
)

var inspec string

func init() {
	src, err := ioutil.ReadFile("inspec.txt")
	if err != nil {
		panic(err)
	}
	inspec = string(src)
}

func TestInspec(t *testing.T) {
	S := DefaultStops
	T := Tokenize(inspec, DefaultStops)
	K := TextRank(T, 2)
	R := K.RankMap()
	for w := range R {
		if S.Contains(w) || strings.TrimSpace(w) == "" {
			t.Fatalf("Stop-word found in TextRank payload: %s", w)
		}
	}
	t.Log(K.Tokens)
}

func TestStem(t *testing.T) {
	T := Tokenize(inspec, DefaultStops)
	err := T.Stem("english")
	if err != nil {
		t.Fatal(err)
	}
}
