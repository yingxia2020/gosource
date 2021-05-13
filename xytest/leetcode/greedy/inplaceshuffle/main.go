ckage main

import (
	"fmt"
	"math/rand"
)

func inplaceshuffle(input []int) {
	
	for i:=0; i<len(input)-1; i++ {
		// get a random number range from i to len(input)-1, note: rand.Intn(0) cause error!
		randIndex := rand.Intn(len(input)-1-i) + i
		if i != randIndex {
			input[i], input[randIndex] = input[randIndex], input[i]
		}
	}	
}

func main() {
	fmt.Println("Hello, playground")
	
	var input = []int{1, 2, 3, 4, 5, 6, 7}
	inplaceshuffle(input)
	fmt.Println(input)
}

