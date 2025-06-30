package util

import (
	"errors"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math/rand"
	"strings"
	"unicode"
)

// ErrInvalidSize is returned when an invalid size is passed
var ErrInvalidSize = errors.New("wordMap: no words with the specified size")

// ErrNotEnoughWords is returned when there aren't enough words to choose
var ErrNotEnoughWords = errors.New("wordMap: not enough words with the specified size")

// WordMap is a utility type for efficiently storing a word list
type WordMap struct {
	// min/max size of all words in the list
	minSize, maxSize uint32

	// Map from clean to original word
	cleanToOrigMap map[string]string

	// Map organizing words by size; words here are cleaned
	sizeMap map[uint32][]string
}

func (w WordMap) MinWordSize() uint32 {
	return w.minSize
}

func (w WordMap) MaxWordSize() uint32 {
	return w.maxSize
}

func (w WordMap) CleanWord(word string) string {
	return RemoveDiacritics(RemoveWhitespace(strings.ToLower(strings.TrimSpace(word))))
}

// ChooseRandom chooses words randomly. Returned words are cleaned (without diacritics)
//
//   - If there are no words with the specified size, returns ErrInvalidSize
//   - If there aren't enough words with the specified size and count, returns ErrNotEnoughWords
func (w WordMap) ChooseRandom(wordLength, count uint32) ([]string, error) {
	// Ensure wordLength is between min/max sizes
	if wordLength < w.minSize || wordLength > w.maxSize {
		return nil, ErrInvalidSize
	}

	// Get the list of words with the specified wordLength
	var words []string
	var ok bool
	if words, ok = w.sizeMap[wordLength]; !ok {
		return nil, ErrInvalidSize
	}

	// Ensure there are at least the provided count of words in the list
	if count > uint32(len(words)) {
		return nil, ErrNotEnoughWords
	}

	copied := make([]string, len(words))
	copy(copied, words)

	// Shuffle in-place
	rand.Shuffle(len(copied), func(i, j int) {
		copied[i], copied[j] = copied[j], copied[i]
	})

	return copied[:count], nil
}

// GetOriginalWord returns the original word given a cleaned word as input, along with whether it is valid
func (w WordMap) GetOriginalWord(cleanedWord string) (string, bool) {
	origWord, ok := w.cleanToOrigMap[cleanedWord]
	return origWord, ok
}

func WordMapFromList(words []string) WordMap {
	cleanToOrigMap := make(map[string]string)
	sizeMap := make(map[uint32][]string)

	minSize, maxSize := uint32(1000000), uint32(0)
	for _, word := range words {
		// Clean word and store in the clean-to-orig map
		word = strings.ToLower(strings.TrimSpace(word))
		cleaned := RemoveDiacritics(RemoveWhitespace(word))
		if cleaned == "" {
			continue
		}

		cleanToOrigMap[cleaned] = word

		// Store in the size map
		wordLen := uint32(len(cleaned))
		if _, ok := sizeMap[wordLen]; !ok {
			sizeMap[wordLen] = make([]string, 0)
		}
		sizeMap[wordLen] = append(sizeMap[wordLen], word)

		// Update min/max sizes
		if wordLen < minSize {
			minSize = wordLen
		}
		if wordLen > maxSize {
			maxSize = wordLen
		}
	}

	return WordMap{
		minSize:        minSize,
		maxSize:        maxSize,
		cleanToOrigMap: cleanToOrigMap,
		sizeMap:        sizeMap,
	}
}

// RemoveDiacritics removes all diacritics from a text
func RemoveDiacritics(text string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)
	return result
}

// RemoveWhitespace removes all diacritics from a text
func RemoveWhitespace(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)
}
