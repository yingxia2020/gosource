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
	"bytes"
	"fmt"
	"math/rand"
)

func ecommerceSim(buffer *bytes.Buffer, action string, result []string) {
	switch action {
	case "2":
		productLines(buffer)
	case "3":
		productModels(buffer)
	case "4":
		productDetail(buffer)
	case "11":
		getRegions(buffer)
	default:
		errorCase(buffer)
	}
}

func productLines(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")
	for i := minCurr; i <= SCALEDLOAD; i++ {
		buffer.WriteString(fmt.Sprintf("%s%02d\n", Nouns[rand.Intn(currSize)], i))
	}
	buffer.WriteString("</pre>\n")
}

func productModels(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")
	for i := minCurr; i <= SCALEDLOAD; i++ {
		buffer.WriteString(fmt.Sprintf("%s%02d&%s&%s&%s&",
			Nouns[rand.Intn(currSize)], i, getHighlights(), getFeature(), getHighlights()))
		buffer.WriteString(fmt.Sprintf("%s&%s&", getFeature(), getHighlights()))
		buffer.WriteString(fmt.Sprintf("%s&%s\n", getFeature(), getHighlights()))
	}
	buffer.WriteString("</pre>\n")
}

func productDetail(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")
	for i := minCurr; i <= SCALEDLOAD*2; i++ {
		buffer.WriteString(fmt.Sprintf("(%1d) %s\n", i, getOverview()))
	}
	buffer.WriteString("</pre>\n")
}

func getRegions(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("0\n")
	for i := minCurr; i <= maxCurr; i++ {
		buffer.WriteString(fmt.Sprintf("%s%02d\n", Currencies[rand.Intn(currSize)], i))
	}
	buffer.WriteString("</pre>\n")
}

func errorCase(buffer *bytes.Buffer) {
	buffer.WriteString("<pre>\n")
	buffer.WriteString("1\n")
	buffer.WriteString("Unknown action type!\n")
	buffer.WriteString("</pre>\n")
}
