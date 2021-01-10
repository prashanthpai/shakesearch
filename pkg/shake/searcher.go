package shake

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io"
	"io/ioutil"
	"sort"
)

const (
	textGrabLength = 250
)

func NewSearcher() *Searcher {
	return new(Searcher)
}

type Searcher struct {
	completeWorks string
	suffixArray   *suffixarray.Index
}

func (s *Searcher) Load(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %w", err)
	}

	s.completeWorks = string(b)
	// suffixarray.FindAllIndex uses regexp and it would be slower
	// than suffixArray.Lookup for doing a case insensitive search
	s.suffixArray = suffixarray.New(bytes.ToLower(b))

	return nil
}

func (s *Searcher) Search(query string, filter string) []string {
	// even though a specific work is specified, we search the whole
	// corpus, at least for now, to keep things simple

	idxs := s.suffixArray.Lookup(bytes.ToLower([]byte(query)), -1)
	sort.Ints(idxs)

	_, filterWork := worksIndex[filter]
	if filterWork {
		idxs = s.filterByWork(idxs, filter)
	}

	idxs = s.removeOverlaps(idxs)

	results := []string{}
	for _, idx := range idxs {
		// handle out-of-bounds panic bug
		start := idx - textGrabLength
		if start < 0 {
			start = 0
		}
		end := idx + textGrabLength
		if end > len(s.completeWorks) {
			end = len(s.completeWorks)
		}
		text := s.completeWorks[start:end]

		// apply filters/cleaners (can be applied conditionally in future
		// based on a configurable parameter)
		text = startAtSentence(text, query)
		text = endAtSentenceOrWord(text)
		text = removeEmptyLines(text)

		results = append(results, text)
	}

	return results
}

func (s *Searcher) removeOverlaps(idxs []int) []int {
	if len(idxs) <= 1 {
		return idxs
	}

	var result = []int{idxs[0]}
	for i := 1; i < len(idxs); i++ {
		diff := idxs[i] - result[len(result)-1]
		if diff > textGrabLength {
			result = append(result, idxs[i])
		}
	}

	return result
}

func (s *Searcher) filterByWork(idxs []int, work string) []int {
	start, end := worksIndex[work][0], worksIndex[work][1]

	var result []int
	for _, idx := range idxs {
		if start <= idx && idx <= end {
			result = append(result, idx)
		}
		if idx > end {
			break
		}
	}

	return result
}

func (s *Searcher) Filters() []string {
	var l []string
	for k := range worksIndex {
		l = append(l, k)
	}

	return l
}
