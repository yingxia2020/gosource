package main

import (
	"fmt"
)

func isPermutationPalindrome(input string) bool {
	var sum = map[string]int{}
	for _, c := range input {
		if _, ok := sum[string(c)]; ok {
			delete(sum, string(c))
		} else {
			sum[string(c)] = 1
		}
	}

	return len(sum) <=1	
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = "civiccentcenti"
	fmt.Println(isPermutationPalindrome(input))
}

