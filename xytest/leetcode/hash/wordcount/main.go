package main

import (
	"fmt"
	"regexp"
	"strings"
)

func wordcount(input string) map[string]int {
	var sum = map[string]int{}
	// one way
	sm := regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(input, " ")
	// split it by space.
	words := strings.Split(sm, " ")
	fmt.Println(words)
	
	// or another
	s := regexp.MustCompile(`[^a-zA-Z]+`)
	tokens := s.Split(input, -1)
	fmt.Println(tokens)
	for i, token := range tokens {
		// this step is needed, since last element is ""
		if token == "" {
			fmt.Println(i, "really")
			continue
		}

		if _, ok := sum[token]; ok {
			sum[token]++
		} else {
			sum[token] = 1
		}
	}
	return sum
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = "Cliff finished his 'cake' and paid the bill. Bill  finished his cake at the edge of the cliff."
	fmt.Println(wordcount(input))
}
