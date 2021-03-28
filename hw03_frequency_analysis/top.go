package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const TopWordAmount = 10

type WordDescr struct {
	Counter int
	Word    string
}

func Top10(userText string) []string {
	result := make([]string, 0)
	wordList := splitText(userText)

	if len(wordList) == 0 {
		return result
	}
	scoreList := countWords(wordList)
	sortScore(scoreList)

	topCounter := TopWordAmount
	if len(scoreList) < topCounter {
		topCounter = len(scoreList)
	}
	for i := 0; i < topCounter; i++ {
		result = append(result, scoreList[i].Word)
	}
	return result
}

func splitText(s string) []string {
	return regexp.MustCompile(`\s`).Split(s, -1)
}

func countWords(wordList []string) []*WordDescr {
	scoreMap := make(map[string]*WordDescr)
	scoreList := make([]*WordDescr, 0)

	for _, w := range wordList {
		w = strings.Trim(w, `,.!-;`)
		w = strings.ToLower(w)
		if len(w) == 0 {
			continue
		}
		descr, ok := scoreMap[w]
		if !ok {
			descr = &WordDescr{Counter: 0, Word: w}
			scoreMap[w] = descr
			scoreList = append(scoreList, descr)
		}
		descr.Counter++
	}
	return scoreList
}

func sortScore(s []*WordDescr) {
	sort.Slice(s, func(i, j int) bool {
		if s[i].Counter == s[j].Counter {
			return s[i].Word < s[j].Word
		}
		return s[i].Counter > s[j].Counter
	})
}
