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
	"encoding/csv"
	"os"
)

func main() {
	file, err := os.OpenFile("test.csv", os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	x := []string{"Country", "City", "Population"}
	y := []string{"Japan", "Tokyo", "923456"}
	z := []string{"Australia", "Sydney", "789650"}
	csvWriter := csv.NewWriter(file)
	strWrite := [][]string{x, y, z}
	csvWriter.WriteAll(strWrite)
	csvWriter.Write([]string{"hello", "world"})
	csvWriter.Flush()
}
