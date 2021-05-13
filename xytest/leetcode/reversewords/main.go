ckage main

import (
	"fmt"
)

func swap(orig []string, start, end int) {
	for start < end {
		orig[start], orig[end] = orig[end], orig[start]
		start++
		end--
	}
}

func main() {
	fmt.Println("Hello, playground")
	var orig = []string{"w", "o", "r", "l", "d", "", "h", "e", "l", "l", "o", "", "s", "a", "y"}
	var start = 0
	var end = len(orig) - 1
	swap(orig, start, end)

	for i, c := range orig {
		if c == "" {
			swap(orig, start, i-1)
			start = i+1
		} else if i == len(orig)-1 {
			swap(orig, start, i)
		}
	}	
	fmt.Println(orig)
}

