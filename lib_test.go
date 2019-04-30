package textrank

import (
	"testing"
	"io/ioutil"
	"strings"
)

func TestInspec(t *testing.T) {
	src, err := ioutil.ReadFile("inspec.txt")
	if err != nil {
		t.Fatal(err)
	}
	K := TextRank(string(src), 2, DefaultStops)
	for w := range K.RankMap() {
		if DefaultStops.Contains(w) || strings.TrimSpace(w) == "" {
			t.Fatalf("Stop-word found in TextRank payload: %s", w)
		}
	}
	t.Log(K.Tokens)
}
