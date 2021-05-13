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
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

var OCRSLA = []int{1500, 2000, 2500}
var OCRTHROUGHPUT = []int{10, 20, 30}

func main() {
	// create result file
	rf, err := os.Create("a.result")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rf.Close()
	rTable := tablewriter.NewWriter(rf)
	rTable.SetCenterSeparator("+")
	rTable.SetColumnSeparator("|")
	rTable.SetRowSeparator("-")
	rHeader := []string{"SLA (MSEC)", "RATE (REQ/SEC)"}
	rTable.SetHeader(rHeader)
	rTable.SetAutoFormatHeaders(false)

	// Write to result table
	for idx, item := range OCRSLA {
		contents := []string{fmt.Sprintf("%d", item),
			fmt.Sprintf("%d", OCRTHROUGHPUT[idx])}
		rTable.Append(contents)
	}
	rTable.Render()
}
