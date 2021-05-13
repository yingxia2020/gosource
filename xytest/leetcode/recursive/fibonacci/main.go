package main

import (
	"fmt"
)

func fibon(input int) int {
	if input == 0 {
		return 0
	} else if input == 1 {
		return 1
	}
	a := 0
	b := 1
	var sum int
	for i:=2; i<=input; i++ {
		sum = a + b
		a = b
		b = sum
	}
	return sum
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = 10

	fmt.Println(fibon(input))
}
