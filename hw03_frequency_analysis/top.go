package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type keyValue struct {
	key   string
	value int
}

func Top10(s string) []string {
	var trimmed string

	result := make(map[string]int)
	words := strings.Fields(s)

	for _, word := range words {
		trimmed = strings.TrimFunc(strings.ToLower(word), unicode.IsPunct)

		if len(word) > 1 && trimmed == "" {
			result[word]++
		}

		if trimmed != "" {
			result[trimmed]++
		}
	}

	sliceResult := make([]keyValue, 0, len(result))
	for key, value := range result {
		sliceResult = append(sliceResult, keyValue{key, value})
	}

	sort.Slice(sliceResult, func(i, j int) bool {
		if sliceResult[i].value == sliceResult[j].value {
			return sliceResult[i].key < sliceResult[j].key
		}

		return sliceResult[i].value > sliceResult[j].value
	})

	answer := make([]string, 0, 10)
	for i, kv := range sliceResult {
		if i >= 10 {
			break
		}
		answer = append(answer, kv.key)
	}

	return answer
}
