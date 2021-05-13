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
	var x = 0.0001

	for i := 0; i <= 5000000; i++ {
		x += math.Sqrt(x)
	}
	fmt.Println("Result=", x)
}
