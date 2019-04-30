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
	S := DefaultStops
	K := TextRank(string(src), 2, S)
	for w := range K.RankMap() {
		if S.Contains(w) || strings.TrimSpace(w) == "" {
			t.Fatalf("Stop-word found in TextRank payload: %s", w)
		}
	}
	t.Log(K.Tokens)
}
