package textrank

import (
	"testing"
	"io/ioutil"
	"strings"
    "sort"
    "github.com/aaaton/golem"
    "github.com/aaaton/golem/dicts/en"
)

var (
    inspec string
    L *golem.Lemmatizer
)

func init() {
	src, err := ioutil.ReadFile("inspec.txt")
	if err != nil {
		panic(err)
	}
	inspec = string(src)
    L, err = golem.New(en.NewPackage())
    if err != nil {
        panic(err)
    }
}

func TestInspec(t *testing.T) {
	S := DefaultStops
	T := Tokenize(inspec, DefaultStops)
    T.Lemmatize(L)
	R := TextRank(T, 2)
	for w := range R {
		if S.Contains(w) || strings.TrimSpace(w) == "" {
			t.Fatalf("Stop-word found in TextRank payload: %s", w)
		}
	}
    T = make([]string, 0, len(R))
    for t := range R {
        T = append(T, t)
    }
    sort.Slice(T, func(i, j int) bool {
        return R[T[i]] > R[T[j]]
    })
	t.Log(T)
}

