package main

import (
	"fmt"
)

func fillmovies(input1 []int, input2 int) bool {
	var movies = map[int]int{}
	
	for _, t1 := range input1 {
		t2 := input2 - t1
		if _, ok := movies[t2]; ok {
			return true
		}
		movies[t1]=1
	}
	
	return false
}

func main() {
	fmt.Println("Hello, playground")
	
	var input1 = []int{1, 1, 2, 3, 1, 4}
	var input2 = 3
	fmt.Println(fillmovies(input1, input2))
}
