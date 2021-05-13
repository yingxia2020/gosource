package main

import (
	"fmt"
)

func sortscores(input []int, highest int) []int {
	var result []int
	var store = make([]int, highest+1)
	for _, i := range input {
		store[i]++
	}
	
	for i := len(store)-1; i>=0; i-- {
		for j:=0; j<store[i]; j++ {
			result = append(result, i)
		}
	}	
	
	return result
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []int{53, 78, 64, 72, 41, 37, 91, 89, 84, 60, 65, 65, 53, 78, 64, 72}
	var highest = 100
	fmt.Println(sortscores(input, highest))
}

