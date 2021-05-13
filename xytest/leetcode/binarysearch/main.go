package main

import (
	"fmt"
)

func findrotationpoint(input []string) int {
	start := 0
	end := len(input)-1
	
	for input[start] > input[end] && (end - start) > 1 {
		mid := (start+end)/2
		if input[mid] < input[end] {
			end = mid
		} else {
			start = mid
		}
	}
	return end
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []string{
	       "undulate",
	       "xenoepist",
	       "zena",
	       "asymptote", // <-- rotates here!
	       "babka",
	       "banoffee",
	       "coffee",
	       "engender",
	       "karpatka",
	       "othellolagkage",
	       "ptolemaic",
	       "retrograde",
	       "supplant",
        }

	fmt.Println(findrotationpoint(input))
}
