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
)

func main() {

	/* Use 5000000 times original SQRT calculations to scale up on powerful machines */
	var x = 0.0001

	for i := 0; i <= 5000000; i++ {
		x += math.Sqrt(x)
	}

	fmt.Printf("result=%f OK\n\n", x)
}
