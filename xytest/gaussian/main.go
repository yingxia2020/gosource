package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/leesper/go_rng"
)

func main() {
	fmt.Println("=====Testing for GaussianGenerator begin=====")
	grng := rng.NewPoissonGenerator(time.Now().UnixNano())
	//grng := rng.NewPoissonGenerator(int64(1582221789833178776))
	fmt.Println("Poisson(30.0): ")
	hist := map[int64]int{}
	for i := 0; i < 10000; i++ {
		hist[int64(grng.Poisson(10.0))]++
	}

	var keys []int
	for k := range hist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("%d:\t%s\n", key, strings.Repeat("*", hist[int64(key)]/50))
	}
	
	for _, key := range keys {
		fmt.Printf("%d:\t%d\n", key, hist[int64(key)])
	}
	fmt.Println("=====Testing for GaussianGenerator end=====")
	fmt.Println()
}
