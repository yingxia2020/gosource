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
	"time"
)

func main() {
	fmt.Println(time.Now().Format(time.RFC850))
	fmt.Println(time.Now().Format("20060102150405"))
}
