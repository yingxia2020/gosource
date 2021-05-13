package main

import (
	"fmt"
)

func Max(nums ...int) int {
	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func Min(nums ...int) int {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}


func maxproductof3(input []int) int {
	if len(input) <= 2 {
		return 0
	}	
	
	min := Min(input[0], input[1])
	max := Max(input[0], input[1])
	min2 := input[0]*input[1]
	max2 := input[0]*input[1]
	max3 := input[0]*input[1]*input[2]
	
	for i:= 2; i< len(input); i++ {
		max3 = Max(max3, input[i]*max2, input[i]*min2)
		max2 = Max(max2, input[i]*max, input[i]*min)
		min2 = Min(min2, input[i]*max, input[i]*min)
		if input[i] < min {
			min = input[i]
		}
		if input[i] > max {
			max = input[i]
		}
	}
	return max3
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []int{1, 10, -5, 2, 15, 20, 50, -100}
	fmt.Println(maxproductof3(input))
}
