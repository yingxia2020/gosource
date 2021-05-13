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

//go:generate go run main.go

import (
	"os"

	"github.com/wcharczuk/go-chart"
)

func main() {
	/*
	   In this example we set a custom range for the y-axis, overriding the automatic range generation.
	   Note: the chart will still generate the ticks automatically based on the custom range, so the intervals may be a bit weird.
	*/

	graph := chart.Chart{
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 10.0,
			},
		},
		XAxis: chart.XAxis{
			Range: &chart.Range{
				Min: 0,
				Max: 4,
			},
			Name: "This is a test",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
