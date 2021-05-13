package main

import (
	"fmt"
)

func productofothers(input []int) []int {
	if len(input) <= 1 {
		return []int{}
	}	
	
	var left int
	var right int
	var total = make([]int, len(input))
	for i, _ := range input {
		if i==0 {
			left = 1
		} else {
			left = left*input[i-1]
		}
		total[i] = left
	}

	for i:= len(input)-1; i>=0; i-- {
		if i == len(input)-1 {
			right = 1			
		} else {
			right= right*input[i+1]
		}
		total[i] = total[i]*right
	}

	return total	
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []int{1, 5, 7, 3, 4}
	fmt.Println(productofothers(input))
}

