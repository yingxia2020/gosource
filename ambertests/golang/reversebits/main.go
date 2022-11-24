package main

import (
	"fmt"
)

func setBit(n int, pos int) int {
	n |= (1 << pos)
	return n
}

func clearBit(n int, pos int) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func hasBit(n int, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func reverseBits(v int, n, p int) int {
	var start = n
	var end = n + p - 1
	for start < end {
		endIsSet := hasBit(v, end)
		startIsSet := hasBit(v, start)
		if endIsSet {
			v = setBit(v, start)
		} else {
			v = clearBit(v, start)
		}
		if startIsSet {
			v = setBit(v, end)
		} else {
			v = clearBit(v, end)
		}
		start += 1
		end -= 1
	}
	return v
}

func main() {
	var v = 0xF80AB0A5
	fmt.Printf("Original value: %b\n", v)
	revertv := reverseBits(v, 0, 4)
	fmt.Printf("Reversed value: %b\n", revertv)
	revertv = reverseBits(v, 5, 5)
	fmt.Printf("Reversed value: %b\n", revertv)
}
