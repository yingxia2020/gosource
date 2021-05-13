/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

import (
	"fmt"
	"math"
	"math/rand"
	//"os"
	"runtime"
	//"strconv"
	"sync"
	"time"
)

func monte_carlo_pi(radius float64, reps int, result *int, wait *sync.WaitGroup) {
	var x, y float64
	count := 0
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	for i := 0; i < reps; i++ {
		x = random.Float64() * radius
		y = random.Float64() * radius

		if num := math.Sqrt(x*x + y*y); num < radius {
			count++
		}
	}

	*result = count
	wait.Done()
}

func main() {
	cores := runtime.NumCPU()
	fmt.Println(cores)
	cores = 4
	runtime.GOMAXPROCS(cores)

	var wait sync.WaitGroup

	counts := make([]int, cores)

	// samples, _ := strconv.Atoi(os.Args[1])
	samples := 1000000000

	start := time.Now()
	wait.Add(cores)

	for i := 0; i < cores; i++ {
		go monte_carlo_pi(100.0, samples/cores, &counts[i], &wait)
	}

	wait.Wait()

	total := 0
	for i := 0; i < cores; i++ {
		total += counts[i]
	}

	pi := (float64(total) / float64(samples)) * 4

	fmt.Println("Time: ", time.Since(start))
	fmt.Println("pi: ", pi)
	fmt.Println("")
}
