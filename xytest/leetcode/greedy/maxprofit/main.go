package main

import (
	"fmt"
)

func maxprofit(input []int) int {
	if len(input) <= 1 {
		return 0
	}	
	
	min := input[0]
	maxprofit := input[1] - input[0]
	
	for i:= 1; i< len(input)-1; i++ {
		if input[i] < min {
			min = input[i]
		}
		if input[i+1] - min > maxprofit {
			maxprofit = input[i+1] - min
		}
	}
	return maxprofit
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []int{10, 7, 5, 8, 11, 9, 12}
	fmt.Println(maxprofit(input))
}

