package textrank

import (
	"math"
	"unicode"
	"hash/adler32"
	"sort"
	"strings"

	"golang.org/x/text/transform"
	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"

	"github.com/kavorite/kwdx"
	"github.com/alixaxel/pagerank"
	"gopkg.in/jdkato/prose.v2"
)

type Stopwords map[string]struct{}

func Stops(tokens ...string) (S Stopwords) {
	S = make(Stopwords, len(tokens))
	for _, t := range tokens {
		S[t] = struct{}{}
	}
	return
}

func (S Stopwords) Contains(t string) bool {
	if S == nil {
		return false
	}
	_, ok := S[t]
	return ok
}

var normTx = transform.Chain(
	runes.Map(unicode.ToLower),
	norm.NFD,
	transform.RemoveFunc(func(r rune) bool {
		return !(unicode.Is(unicode.L, r) || unicode.Is(unicode.N, r))
	}),
	norm.NFC)

func normalize(x string) string {
	rtn, _, _ := transform.String(normTx, x)
	return strings.TrimSpace(rtn)
}

func tokenize(x string, S Stopwords) []string {
	doc, _ := prose.NewDocument(x)
	rtn := make([]string, 0, len(x) / 3)
	for _, t := range doc.Tokens() {
		p := t.Tag[0]
		if p != 'V' && p != 'F' && p != 'N' && t.Tag != "PRP" {
			continue
		}
		w := normalize(t.Text)
		if S.Contains(w) || strings.TrimSpace(w) == "" {
			continue
		}
		rtn = append(rtn, w)
	}
	return rtn
}

func hash(x string) uint32 {
	return adler32.Checksum([]byte(x))
}

func TextRank(x string, w uint, S Stopwords) (K kwdx.Keywords) {
	G := pagerank.NewGraph()
	T := tokenize(x, S)
	D := make(map[uint32]string, len(T))
	c := float64(w/2)
	for i := w; i < uint(len(T))-w; i++ {
		W := T[i-w:i+w+1]
		t := hash(T[i])
		for j, w := range W {
			k := hash(w)
			D[k] = w
			G.Link(t, k, math.Abs(c-float64(j)))
		}
	}
	K.Tokens = make([]string, 0, len(D))
	K.Rankings = make([]float64, 0, len(D))
	G.Rank(0.85, 1e-6, func(t uint32, r float64) {
		K.Tokens = append(K.Tokens, D[t])
		K.Rankings = append(K.Rankings, r)
	})
	sort.Sort(K)
	return
}
