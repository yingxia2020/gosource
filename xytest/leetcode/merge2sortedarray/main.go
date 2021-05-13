package main

import (
	"fmt"
)


func main() {
	fmt.Println("Hello, playground")
	
	var input1 = []int{1, 3, 5, 7, 8, 9, 10}
	var input2 = []int{2, 3, 4, 6, 7}
	var output []int
	
	var start1 = 0
	var start2 = 0
	
	for start1 < len(input1) && start2 < len(input2) {

		if input1[start1] < input2[start2] {
			output = append(output, input1[start1])
			start1++
		} else if input1[start1] > input2[start2] {
			output = append(output, input2[start2])
			start2++
		} else {
			output = append(output, input1[start1])
			start1++
			start2++
		}	
		fmt.Println(start1, start2)
	}
	
		if start1 == len(input1) {
			output = append(output, input2[start2:]...)
			/*
			for start2 < len(input2) {
				output = append(output, input2[start2])
				start2++
			}
			*/
		} else 	if start2 == len(input2) {
			output = append(output, input1[start1:]...)
		}
	fmt.Println(output)
}
