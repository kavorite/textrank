// TODO: n-gram keyphrase extraction?
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
	"github.com/kljensen/snowball"
)

func hash(x string) uint32 {
	return adler32.Checksum([]byte(x))
}

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

// Tokens represent a series of unigrams ascribed to a document.
type Tokens []string

// Tokenize splits the given document into its constituent unigrams.
func Tokenize(x string, S Stopwords) Tokens {
	doc, _ := prose.NewDocument(x)
	rtn := make([]string, 0, len(x) / 3)
	for _, t := range doc.Tokens() {
		p := t.Tag[0]
		if p != 'V' && p != 'F' && p != 'N' && t.Tag != "PRP" {
			continue
		}
		w := normalize(t.Text)
		if S.Contains(w) || w == "" {
			continue
		}
		rtn = append(rtn, w)
	}
	return rtn
}

// Stem performs stemming on the tokens in-place in the given language. Errors
// out if the language given is unsupported.
func (T Tokens) Stem(lang string) error {
	for i := 0; i < len(T); i++ {
		stemmed, err := snowball.Stem(T[i], lang, true)
		if err != nil {
			return err
		}
		T[i] = stemmed
	}
	return nil
}

// StemTable returns a vocabulary dictionary mapping tokens to their stems.
// Errors out if the language given is unsupported.
func (T Tokens) StemTable(lang string) (map[string]string, error) {
	D := make(map[string]string, len(T))
	for _, t := range T {
		s, err := snowball.Stem(t, lang, true)
		if err != nil {
			return nil, err
		}
		D[t] = s
	}
	return D, nil
}

// TextRank tokenizes and performs keyword extraction on the given tokens with
// window size = `w` (i.e., a context of 2w+1 words is examined each
// iteration).
func TextRank(T Tokens, w uint) (K kwdx.Keywords) {
	G := pagerank.NewGraph()
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
	sort.Sort(sort.Reverse(K))
	return
}
