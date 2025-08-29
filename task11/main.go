package main

import (
	"fmt"
	"slices"
	"strings"
)

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	result := processAnagrams(input)
	fmt.Println(result)
}

func processAnagrams(input []string) map[string][]string {
	anagrams := make(map[string][]string)
	for _, word := range input {
		runes := []rune(strings.ToLower(word))
		//pattern-defeating quicksort в среднем O(n log n)
		slices.Sort(runes)
		anagrams[string(runes)] = append(anagrams[string(runes)], word)
	}

	result := make(map[string][]string)
	for _, group := range anagrams {
		if len(group) == 1 {
			continue
		}

		result[group[0]] = group
		slices.Sort(group)
	}

	return result
}
